package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/kubecost/opencost/pkg/kubecost"
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

var Plans = []Plan{"Starter", "Pro", "Elite"}

type State string

const (
	Grace  State = "Grace"
	Closed       = "Closed"
	Active       = "Active"
)

// project_id and Cluster_id are taken from rancher from the Id field (not the UUID)
type Project struct {
	gorm.Model
	ProjectId           string     `json:"projectId" gorm:"primaryKey;unique"`
	ClusterId           string     `json:"clusterId"`
	CreationTimeStamp   time.Time  `json:"creationTimeStamp"`
	State               State      `json:"State"`
	Plan                Plan       `json:"plan"`
	History             []BillFile `json:"history" gorm:"foreignKey:ProjectRefer;references:ProjectId"`
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
	ProjectRefer        string    `json:"ProjectRefer"`
	BillingAccountRefer uint
}

type BillingAccount struct {
	gorm.Model
	UUID             uuid.UUID      `json:"uuid" gorm:"primaryKey;unique"`
	BillingAdmins    []AdminDetails `json:"billingAdmins" gorm:"many2many:BillingAccount_Admin;"`
	BillingStartDate time.Time      `json:"billingStartDate"`
	Balance          float64        `json:"balance"`
	Company          Company        `json:"company" gorm:"embedded"`
	Projects         []Project      `json:"projects" gorm:"foreignKey:BillingAccountRefer;references:UUID"`
}

type CreateBillingAccount struct {
	BillingAdmins []AdminDetails `json:"billingAdmins"`
	Company       Company        `json:"company"`
	Projects      []Project      `json:"projects"`
}

type AddProjectModel struct {
	BillingAccountUUID uuid.UUID `json:"billing_account_uuid"`
	ProjectId          string    `json:"project_id"`
	ClusterId          string    `json:"clusterId"`
	CreationTimeStamp  time.Time `json:"creationTimeStamp"`
	Plan               Plan      `json:"plan"`
	State              State     `json:"state"`
}
type AllocationResponse struct {
	Code int                              `json:"code"`
	Data []map[string]kubecost.Allocation `json:"data"`
}

type Metrics struct {
	CPUCoreHours         float64
	CpuAverageUsage      float64
	RamByteMinutes       float64
	RamAverageUsage      float64
	NetworkTransferBytes float64
	NetworkReceiveBytes  float64
	PvByteHours          float64
}

type ResourcePricing struct {
	CPUCoreHourPrice          float64
	MemoryByteHourPrice       float64
	NetworkReceiveBytesPrice  float64
	NetworkTransferBytesPrice float64
	PVByteHourPrice           float64
}
