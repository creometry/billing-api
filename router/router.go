package router

import (
	controllers "billing-api/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
)

type Namespace struct {
	ClientSet *kubernetes.Clientset
}

func SetupRouter(db *gorm.DB) {

	router := gin.Default()
	router.POST("/v1/CreateBillingAccount", controllers.CreateBillingAccount)
	router.GET("/v1/getBillingAccount/:uuid", controllers.GetBillingAccount)
	router.GET("/v1/GetBillingAccountsByAdminUUID/:uuid", controllers.GetBillingAccountsByAdminUUID)
	router.POST("/v1/addProject", controllers.AddProject)
	router.GET("/v1/listBillingAccountNamespaces", controllers.ListBillingAccountNamespaces)
	router.Run("localhost:8080")
}
