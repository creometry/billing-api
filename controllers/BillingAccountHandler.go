package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	data "billing-api/data"
	"billing-api/logging"
	models "billing-api/models"
	"billing-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k0kubun/pp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var DB, err = data.InitializeDB()

func CreateBillingAccount(c *gin.Context) {
	// log.Println("creating billing account by ", billingAccount.billingAdmins[0])

	// Validate input

	//TODO: generate UUID for each newly created billing account

	var input models.CreateBillingAccount
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// pp.Println("input", input)
	accountDetails := models.BillingAccount{
		UUID:             uuid.New(),
		BillingAdmins:    input.BillingAdmins,
		BillingStartDate: time.Now(),
		Balance:          0.0,
		Company:          input.Company,
		Projects:         input.Projects,
	}

	if result := DB.Table("billing_accounts").Create(&accountDetails); result.Error != nil {

		ApiError := utils.NewAPIError(409, "Conflict", "Database error", "UUID already exists", "UUID already exists")
		c.JSON(http.StatusConflict, ApiError)
		return
	}

	pp.Println("accountDetails", accountDetails)
	c.JSON(http.StatusCreated, accountDetails)

}

func fetchBillingAccountbyuuid(uuid string) (*models.BillingAccount, error) {
	var billingAccount models.BillingAccount

	if result := DB.First(&billingAccount, "uuid = ?", uuid); result.Error != nil {
		return nil, result.Error
	}

	return &billingAccount, nil
}

func GetBillingAccount(c *gin.Context) {
	uuid := c.Param("uuid")

	var BillingAccount models.BillingAccount

	if result := DB.Find(&BillingAccount, "uuid = ?", uuid); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusOK, &BillingAccount)
}

func GetBillingAccountsByAdminUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	var adminDetails models.AdminDetails
	var BillingAccounts []models.BillingAccount

	_ = BillingAccounts

	if result := DB.First(&adminDetails, "uuid = ?", uuid); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	// Retrieve all billing accounts the admin's uuid is associated with

	err := DB.Preload("Projects").
		Preload("BillingAdmins").
		Joins("INNER JOIN public.billing_account_admins ON public.billing_accounts.uuid = public.billing_account_admins.billing_account_uuid").
		Where("public.billing_account_admins.admin_details_uuid = ?", uuid).
		Find(&BillingAccounts).Error

	pp.Println("err", err)

	pp.Println("BillingAccounts", BillingAccounts)

	c.JSON(http.StatusOK, &BillingAccounts)

}

func AddProject(c *gin.Context) {
	var input models.AddProjectModel
	var billingAccount models.BillingAccount
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pp.Println("input", input)

	if result := DB.First(&billingAccount, "uuid = ?", input.BillingAccountUUID); result.Error != nil {

		ApiError := utils.NewAPIError(404, "Not Found", "input error", "Billing account UUID does not exist", "")
		c.JSON(http.StatusNotFound, ApiError)
		return
	}

	newproject := models.Project{
		ProjectId:           input.Project.ProjectId,
		ClusterId:           input.Project.ClusterId,
		ProjectCreationTS:   time.Time{},
		State:               input.Project.State,
		Plan:                input.Project.Plan,
		History:             []models.BillFile{},
		BillingAccountRefer: input.BillingAccountUUID.String(),
	}

	if result2 := DB.Create(&newproject); result2.Error != nil {

		ApiError := utils.NewAPIError(500, "Internal Server Error", "Internal Server Error", "Could not create project", "")
		c.JSON(http.StatusInternalServerError, ApiError)
		return
	}

	c.JSON(http.StatusCreated, newproject)

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

	pp.Println("namespaces", namespaces)

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

//TODO: test this function after filling the database with data
func ListBillingAccountNamespaces(c *gin.Context) {
	var BillingAccountNamespaces []string
	billingAccount_uuid := c.Request.URL.Query().Get("billingAccount_uuid")

	//get billing account by uuid
	billingAccount, err := fetchBillingAccountbyuuid(billingAccount_uuid)
	if err != nil {
		logging.Error(logging.HTTPError, err)
		c.JSON(http.StatusBadRequest, utils.NewAPIError(400, "Bad Request", "input error", "billingAccount_uuid param is required", ""))
	}

	// loop over billing account projects and add their namespace to BillingAccountNamespaces array
	for _, project := range billingAccount.Projects {
		projectNamespace, err := ListProjectNamespaces(project.ProjectId)
		if err != nil {
			logging.Error(logging.HTTPError, err)
			return
		}
		BillingAccountNamespaces = append(BillingAccountNamespaces, projectNamespace...)
	}

	// return projectsNamespaces array
	c.JSON(http.StatusOK, &BillingAccountNamespaces)
}

func getBillingAccountNamespaces(billingAccount_uuid uuid.UUID) ([]string, error) {
	var BillingAccountNamespaces []string
	//get billing account by uuid
	billingAccount, err := fetchBillingAccountbyuuid(billingAccount_uuid.String())
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

func getNamespaceMetrics(namespaceId string) (models.Metrics, error) {
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
		return models.AllocationResponse{}, err
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
		CPUCoreHours:         kubecostResponse.Data[0]["__idle__"].CPUCoreHours,
		CpuAverageUsage:      kubecostResponse.Data[0]["__idle__"].CPUCoreUsageAverage,
		RamByteMinutes:       kubecostResponse.Data[0]["__idle__"].RAMByteHours,
		RamAverageUsage:      kubecostResponse.Data[0]["__idle__"].RAMBytesUsageAverage,
		NetworkTransferBytes: kubecostResponse.Data[0]["__idle__"].NetworkTransferBytes,
		NetworkReceiveBytes:  kubecostResponse.Data[0]["__idle__"].NetworkReceiveBytes,
		//TODO: calculate PV total
		// pvByteHours: metrics.Data[0]["__idle__"].PVs,
	}
	return namespaceMetrics, nil
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
			projectNamespace, err := ListProjectNamespaces(project.ProjectId)
			if err != nil {
				logging.Error(logging.HTTPError, err)
				return 0.0, err
			}

			// get metrics for each namespace
			for _, namespace := range projectNamespace {
				namespaceMetrics, err := getNamespaceMetrics(namespace)
				if err != nil {
					logging.Error(logging.HTTPError, err)
					return 0.0, err
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
	_ = consumedCredit
	//TODO: update billing account balance
	billingAccount.Balance = billingAccount.Balance - consumedCredit
}

func fetchAllBillingAccounts() ([]models.BillingAccount, error) {
	var billingAccounts []models.BillingAccount
	if result := DB.Find(&billingAccounts); result.Error != nil {
		logging.Error(logging.HTTPError, result.Error)
		return nil, result.Error
	}
	return billingAccounts, nil
}

//TODO: implement generateAndSavePDF() function
func generateAndSavePDF() string {
	return "pdf saved!"
}

//TODO: write function that generates bills for each billing account when 5 days before next billing cycle
func generatebills() {
	// get billing accounts within 5 days of the next billing cycle
	billingAccounts, err := fetchAllBillingAccounts()
	if err != nil {
		logging.Error(logging.HTTPError, err)
	}

	for _, billingAccount := range billingAccounts {
		// TODO: verify this date comparaison works
		if billingAccount.BillingStartDate.Add(-5*24*time.Hour).Day() == time.Now().Day() {
			// get billing account consumed credit
			updateBillingAccountBalence(billingAccount)
			if err != nil {
				logging.Error(logging.HTTPError, err)
			}

		}

		// generate pdf
	}
}
