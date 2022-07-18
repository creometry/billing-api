package config

import (
	log "billing-api/logging"
	"os"

	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
)

func LoadDotEnvVariables() error {

	// // load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(log.ConfigError, err)
		return err
	}

	pp.Println("Postgresqlpassword", os.Getenv("Postgresqlpassword"))
	return nil
}
