package main

import (
	data "billing-api/data"
	log "billing-api/logging"
	Router "billing-api/router"
	"billing-api/utils"
)

func main() {
	// err := config.LoadDotEnvVariables()
	// if err != nil {
	// 	log.Fatal(log.ConfigError, err)
	// }

	data.ParseFiles()

	utils.GetKubernetesClient()

	db, err := data.InitializeDB()
	if err != nil {
		log.Fatal(log.ConfigError, err)
	}

	data.InitializeMigrations()

	Router.SetupRouter(db)
	// log.Error(log.ConfigError, "Testing the error")

}
