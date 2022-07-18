package controllers

import (
	"net/http"
	"time"

	data "billing-api/data"
	models "billing-api/models"
	utils "billing-api/utils"

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

	pp.Println("input", input)
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
		// c.JSON(http.StatusConflict, ApiError)
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

	// if result := DB.Find(&BillingAccounts, "uuid = ?", uuid); result.Error != nil {
	// 	c.AbortWithError(http.StatusNotFound, result.Error)
	// 	return
	// }

	// DB.Model(&adminDetails).Select("adminDetails.uuid").Joins("left join billing_accounts on billing_accounts.uuid = users.id").Scan(&result{})

	// c.JSON(http.StatusOK, &BillingAccounts)
	c.JSON(http.StatusOK, &adminDetails)

}
