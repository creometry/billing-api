package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

type Plan string

const (
	PayPerUse string = "PayPerUse"
	Starter          = "Starter"
	Pro              = "Pro"
	Elite            = "Elite"
)

type Project struct {
	gorm.Model
	ProjectId           uuid.UUID  `json:"projectId" gorm:"primaryKey"`
	ClusterId           uuid.UUID  `json:"clusterId"`
	CreationTimeStamp   time.Time  `json:"creationTimeStamp"`
	State               string     `json:"State"`
	Plan                Plan       `json:"accountType"`
	History             []BillFile `json:"history" gorm:"foreignKey:BillingAccountRefer"`
	BillingAccountRefer string     `json:"BillingAccountUUID"`
}

type AdminDetails struct {
	gorm.Model
	UUID                uuid.UUID `json:"uuid" gorm:"primaryKey"`
	Email               string    `json:"email"`
	Phone_number        string    `json:"phone_number"`
	Name                string    `json:"name"`
	BillingAccountRefer string
	// BillingAccountUUID string
}

type Company struct {
	gorm.Model
	IsCompany bool   `json:"isCompany"`
	TaxId     string `json:"TaxId" gorm:"primaryKey"`
	Name      string `json:"name"`
}

type BillFile struct {
	BillingDate         time.Time `json:"BillingDate" gorm:"primaryKey"`
	PdfLink             string    `json:"pdfLink"`
	Amount              float64   `json:"amount"`
	BillingAccountRefer uint
}

type BillingAccount struct {
	gorm.Model
	UUID             uuid.UUID      `json:"uuid" gorm:"primaryKey;unique"`
	BillingAdmins    []AdminDetails `json:"billingAdmins" gorm:"-"`
	BillingStartDate time.Time      `json:"billingStartDate"`
	Balance          float64        `json:"balance"`
	IsActive         bool           `json:"isActive"`
	Company          Company        `json:"company" gorm:"embedded"`
	Projects         []Project      `json:"projects" gorm:"foreignKey:BillingAccountRefer"`
}

type CreateBillingAccount struct {
	BillingAdmins []AdminDetails `json:"billingAdmins"`
	Company       Company        `json:"company"`
	Projects      []Project      `json:"projects"`
}
