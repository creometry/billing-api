package controllers

import (
	"net/http"
	"time"

	data "billing-api/data"
	models "billing-api/models"

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
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	pp.Println("accountDetails", accountDetails)
	c.JSON(http.StatusCreated, accountDetails)

}

func GetBillingAccount(c *gin.Context) {
	uuid := c.Param("uuid")

	var BillingAccount models.BillingAccount

	if result := DB.Find(&BillingAccount, uuid); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusOK, &BillingAccount)
}
