package main

import (
	config "billing-api/config"
	data "billing-api/data"
	log "billing-api/logging"
	// Router "billing-api/router"
)

func main() {
	// err := config.LoadDotEnvVariables()
	// if err != nil {
	// 	log.Fatal(log.ConfigError, err)
	// }

	data.ParseFiles()

	// db, err := data.InitializeDB()
	// if err != nil {
	// 	log.Fatal(log.ConfigError, err)
	// }

	error := config.LoadDotEnvVariables()
	if error != nil {
		log.Fatal(log.ConfigError, error)
	}

	// Router.SetupRouter(db)
	// log.Error(log.ConfigError, "Testing the error")

}
