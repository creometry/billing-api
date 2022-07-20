package router

import (
	controllers "billing-api/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) {

	// db.Migrator().CreateTable(&BillingAccount{})
	// db.Migrator().CreateTable(&Company{})
	// db.Migrator().CreateTable(&Project{})

	router := gin.Default()
	router.POST("/v1/CreateBillingAccount", controllers.CreateBillingAccount)
	router.GET("/v1/getBillingAccount/:uuid", controllers.GetBillingAccount)
	router.GET("/v1/GetBillingAccountsByAdminUUID/:uuid", controllers.GetBillingAccountsByAdminUUID)
	router.POST("/v1/addProject", controllers.AddProject)
	router.Run("localhost:8080")
}
