package billingoperations

import (
	"billing-api/data"
	"billing-api/logging"
	"billing-api/models"
	"billing-api/utils"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid"
	"github.com/k0kubun/pp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var DB, err = data.InitializeDB()

func FetchBillingAccountbyuuid(uuid string) (*models.BillingAccount, error) {
	var billingAccount models.BillingAccount

	if result := DB.First(&billingAccount, "uuid = ?", uuid); result.Error != nil {
		return nil, result.Error
	}

	return &billingAccount, nil
}

func calculateConsumedCredit(namespaceMetrics models.Metrics) float64 {
	var consumedCredit float64

	resourcepricing := data.GetResourcesPrices()

	RamByteHours := namespaceMetrics.RamByteMinutes / 60
	consumedCredit += RamByteHours * resourcepricing.MemoryByteHourPrice
	consumedCredit += namespaceMetrics.CpuAverageUsage * resourcepricing.CPUCoreHourPrice
	consumedCredit += namespaceMetrics.PvByteHours * resourcepricing.PVByteHourPrice
	consumedCredit += namespaceMetrics.NetworkReceiveBytes * resourcepricing.NetworkReceiveBytesPrice
	consumedCredit += namespaceMetrics.NetworkTransferBytes * resourcepricing.NetworkTransferBytesPrice

	return consumedCredit
}

func getBillingAccountConsumedCredit(billingAccount models.BillingAccount) (float64, error) {
	var billingAccountConsumedCredit float64
	// get billing account namespaces
	for _, project := range billingAccount.Projects {
		if project.Plan == "PayPerUse" {
			pp.Println("paymode: PayPerUse")
			projectNamespaces, err := ListProjectNamespaces(project.ProjectId)
			if err != nil {
				logging.Error(logging.HTTPError, err)
				continue
			}

			// get metrics for each namespace
			for _, namespace := range projectNamespaces {
				namespaceMetrics, err := GetNamespaceMetrics(namespace)
				if err != nil {
					logging.Error(logging.HTTPError, err)
					continue
				}

				// calculate consumed credit for each namespace

				consumedCredit := calculateConsumedCredit(namespaceMetrics)

				// update billing account consumed credit
				billingAccountConsumedCredit += consumedCredit
			}
		} else if utils.Contains(models.Plans, project.Plan) {
			planPrice, err := data.GetPlanPrice(project.Plan)
			if err != nil {
				logging.Error(logging.ConfigError, err)
				return 0.0, err
			}
			billingAccountConsumedCredit += planPrice
		}
	}
	// get billing account consumed credit by making the sum of all credit consumed by namespaces
	// return consumed credit
	return billingAccountConsumedCredit, nil
}

func updateBillingAccountBalence(billingAccount models.BillingAccount) {
	consumedCredit, err := getBillingAccountConsumedCredit(billingAccount)
	if err != nil {
		logging.Error(logging.HTTPError, err)
	}
	billingAccount.Balance = billingAccount.Balance - consumedCredit
	pp.Println("billingAccount UUID", billingAccount.UUID.String())
	pp.Println("consumed credit", consumedCredit)
	pp.Println("billingAccount Balance", billingAccount.Balance)
	DB.Save(&billingAccount)
}

func fetchAllBillingAccounts() ([]models.BillingAccount, error) {
	var billingAccounts []models.BillingAccount
	if result := DB.Preload("Projects").Find(&billingAccounts); result.Error != nil {
		logging.Error(logging.HTTPError, result.Error)
		return nil, result.Error
	}
	return billingAccounts, nil
}

// func ListProjectNamespaces(c *gin.Context) {
func ListProjectNamespaces(project_id string) ([]string, error) {
	var projectsNamespaces []string
	clientset := utils.GetKubernetesClient()

	listOptions := metav1.ListOptions{
		LabelSelector: "field.cattle.io/projectId=" + project_id,
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), listOptions)
	if err != nil {
		logging.Error(logging.HTTPError, err)
		return nil, errors.New("error getting namespaces from Kubernetes API")
	}

	if len(namespaces.Items) == 0 {
		// c.JSON(http.StatusNotFound, utils.NewAPIError(404, "Not Found", "input error", "No namespaces found for project or the project does not exist", ""))
		logging.Error(logging.HTTPError, "No namespaces found for project or the project does not exist")
		return nil, errors.New("NO NAMESPACES FOUND FOR PROJECT OR THE PROJECT DOES NOT EXIST")
	}

	for _, namespaceData := range namespaces.Items {
		projectsNamespaces = append(projectsNamespaces, namespaceData.ObjectMeta.Name)
	}

	return projectsNamespaces, nil
}

func getBillingAccountNamespaces(billingAccount_uuid uuid.UUID) ([]string, error) {
	var BillingAccountNamespaces []string
	//get billing account by uuid
	billingAccount, err := FetchBillingAccountbyuuid(billingAccount_uuid.String())
	if err != nil {
		logging.Error(logging.HTTPError, err)
	}

	// loop over billing account projects and add their namespace to BillingAccountNamespaces array
	for _, project := range billingAccount.Projects {
		projectNamespace, err := ListProjectNamespaces(project.ProjectId)
		if err != nil {
			logging.Error(logging.HTTPError, err)
			return nil, err
		}
		BillingAccountNamespaces = append(BillingAccountNamespaces, projectNamespace...)
	}

	// return projectsNamespaces array
	return BillingAccountNamespaces, nil
}

func GetNamespaceMetrics(namespaceId string) (models.Metrics, error) {
	var kubecostUrl string
	// ** kubecost allocation params: **
	// ref : https://github.com/kubecost/docs/blob/main/allocation.md

	// window of billed time
	window := "month"
	// accumulate results (kubecost returns dates seperated by day if this is false (default value) )
	accumulte := "true"
	//field by wich to aggrgate results
	aggregate := "namespace"
	// select namespace by which to filter results
	filterNamespaces := namespaceId

	// kubecostUrl := "kubecost-cost-analyzer"
	if os.Getenv("APP_ENV") == "development" {
		kubecostUrl = "localhost"
	} else {
		kubecostUrl = "kubecost-cost-analyzer"
	}

	// kubecost metrics api
	url := "http://" + kubecostUrl + ":9090/model/allocation?window=" + window + "&accumulate=" + accumulte + "&aggregate=" + aggregate + "&filterNamespaces=" + filterNamespaces

	resp, err := http.Get(url)
	if err != nil {
		logging.Error(logging.HTTPError, "No response from Kubecost!"+err.Error())
		return models.Metrics{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	kubecostResponse := models.AllocationResponse{}
	jsonErr := json.Unmarshal(body, &kubecostResponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	// pp.Println(kubecostResponse)

	namespaceMetrics := models.Metrics{
		CPUCoreHours:         kubecostResponse.Data[0][namespaceId].CPUCoreHours,
		CpuAverageUsage:      kubecostResponse.Data[0][namespaceId].CPUCoreUsageAverage,
		RamByteMinutes:       kubecostResponse.Data[0][namespaceId].RAMByteHours,
		RamAverageUsage:      kubecostResponse.Data[0][namespaceId].RAMBytesUsageAverage,
		NetworkTransferBytes: kubecostResponse.Data[0][namespaceId].NetworkTransferBytes,
		NetworkReceiveBytes:  kubecostResponse.Data[0][namespaceId].NetworkReceiveBytes,
		//TODO: calculate PV total
		// pvByteHours: metrics.Data[0]["__idle__"].PVs,
	}
	return namespaceMetrics, nil
}

func updatePayPerUseHourlyUse() {
	// get projects with PayPerUse plan
	var projects []models.Project
	if result := DB.Find(&projects); result.Error != nil {
		logging.Error(logging.HTTPError, result.Error)
		return nil, result.Error
	}

	// get usage for each project within the last hour
	for _, project := range projects {
	}
	// calculate consumed credit for each project
	// update corresponding billing account for each project
}

//TODO: write function that generates bills for each billing account when 5 days before next billing cycle
func Generatebills() {
	pp.Println("Generate bills called")
	// get billing accounts within 5 days of the next billing cycle
	billingAccounts, err := fetchAllBillingAccounts()
	if err != nil {
		logging.Error(logging.HTTPError, err)
	}

	for _, billingAccount := range billingAccounts {
		// TODO: verify this date comparaison works
		if billingAccount.BillingStartDate.Add(-5*24*time.Hour).Day() == time.Now().Day() || os.Getenv("APP_ENV") == "development" {
			// get billing account consumed credit
			updateBillingAccountBalence(billingAccount)
			if err != nil {
				logging.Error(logging.HTTPError, err)
			}

		}

		// generate pdf
		utils.GeneratePDF("./bills/"+billingAccount.UUID.String()+"_"+time.Now().String()+".pdf", billingAccount)
	}
}
