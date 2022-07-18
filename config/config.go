package config

import (
	log "billing-api/logging"

	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
)

func LoadDotEnvVariables() error {

	pp.Println("loading environnement variables from .env file")
	// // load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(log.ConfigError, err)
		return err
	}

	return nil
}
