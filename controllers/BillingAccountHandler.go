package controllers

import (
	"context"
	"errors"
	"net/http"
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
		IsActive:         true,
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
		CreationTimeStamp:   time.Time{},
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
