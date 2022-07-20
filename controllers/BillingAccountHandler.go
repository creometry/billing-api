package controllers

import (
	"net/http"
	"time"

	data "billing-api/data"
	models "billing-api/models"
	"billing-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k0kubun/pp"
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
