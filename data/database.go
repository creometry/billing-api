package data

import (
	"log"
	"os"

	config "billing-api/config"
	models "billing-api/models"

	"github.com/k0kubun/pp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

var DB *gorm.DB = nil

func InitializeDB() (*gorm.DB, error) {
	config.LoadDotEnvVariables()
	pp.Println("Postgresqlhost", os.Getenv("Postgresqlhost"))
	dsn := "host=" + os.Getenv("Postgresqlhost") + " user=" + os.Getenv("Postgresqluser") + " password=" + os.Getenv("Postgresqlpassword") + " dbname=" + os.Getenv("Postgresqldbname") + " port=" + os.Getenv("Postgresqlport")
	//+ " TimeZone=" + os.Getenv("PostgresqlTimezone")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return db, nil
}

func InitializeMigrations() {
	db, err := InitializeDB()
	if err != nil {
		log.Fatal(err)
	}
	DB = db
	db.AutoMigrate(&models.BillingAccount{}, &models.AdminDetails{}, &models.BillFile{}, &models.Project{}, &models.Company{})
}
