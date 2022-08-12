package billingoperations

import (
	// "fmt"
	// "testing"

	"billing-api/models"
	"time"

	_ "billing-api/utils"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func insertBillingAccountData(db *gorm.DB) error {
	admin1 := models.AdminDetails{
		// UUID:         uuid.FromString("2ad3cc71-bcbe-4fcc-a568-9e22c8b23fb2"),
		UUID:         uuid.New(),
		Email:        "creometryTestAdmin1@creometry.com",
		Phone_number: "11111111",
		Name:         "creometryAdmin1",
	}

	// the naming convention of project is projectName_<< number of admin who created it >>
	// testing needs to cover all the subscription plans.

	// this case tests for project with prepaid plans

	project1_1 := models.Project{
		//TODO: take projectId and ClusterId of system project from Rancher API
		ProjectId:         "placeholder",
		ClusterId:         "placeholder",
		CreationTimeStamp: time.Now().AddDate(0, -1, 0).Add(time.Hour * 8),
		State:             "Active",
		Plan:              "Starter",
	}

	// // this case tests for projects with prepaid plans that were created in the middle of the billing cycle

	project2_1 := models.Project{
		//TODO: take projectId and ClusterId of system project from Rancher API
		ProjectId:         "placeholder2",
		ClusterId:         "placeholder2",
		CreationTimeStamp: time.Now().AddDate(0, 0, -15).Add(time.Hour * 14),
		State:             "Active",
		Plan:              "Closed",
	}

	// // // this case tests for projects with PayPerUse plan

	project3_1 := models.Project{
		//TODO: take projectId and ClusterId of system project from Rancher API
		ProjectId:         "placeholder3",
		ClusterId:         "placeholder3",
		CreationTimeStamp: time.Now().AddDate(0, 0, -28).Add(time.Hour * 9),
		State:             "Active",
		Plan:              "PayPerUse",
	}

	// // // this case tests for when the owner of the project has closed the project or when the project
	// // // has surpassed the grace period without having any sufficient funding in any of the billing accounts linked to it

	project4_1 := models.Project{
		//TODO: take projectId and ClusterId of system project from Rancher API
		ProjectId:         "placeholder4",
		ClusterId:         "placeholder4",
		CreationTimeStamp: time.Now().AddDate(0, 0, -29).Add(time.Hour * 7),
		State:             "Active",
		Plan:              "Closed",
	}

	company1 := models.Company{
		IsCompany: true,
		TaxId:     "123-456-789",
		Name:      "technologyCompany1",
	}

	projects := []models.Project{project1_1, project2_1, project3_1, project4_1}
	admins := []models.AdminDetails{admin1}

	billingAccount1 := models.BillingAccount{
		UUID:             uuid.MustParse("2ad3cc71-bcbe-4fcc-a568-9e22c8b23fb2"),
		BillingAdmins:    admins,
		BillingStartDate: time.Now().AddDate(0, -1, 0),
		Balance:          1200,
		Company:          company1,
		Projects:         projects,
	}

	// uuid2, err2 := uuid.Parse("2ad3cc71-bcbe-4fcc-a568-9e22c8b23fb2")
	// if err2 != nil {
	// }

	// var billingAdmins2 []models.AdminDetails
	// var projects2 []models.Project

	// billingAccount2 := models.BillingAccount{
	// 	UUID:             uuid2,
	// 	BillingAdmins:    billingAdmins2,
	// 	BillingStartDate: time.Now().AddDate(0, -1, 0),
	// 	Balance:          1200,
	// 	Company:          company1,
	// 	Projects:         projects2,
	// }
	// // insert billing account data
	// // err := db.Create(billingAccount1).Error
	// _ = billingAccount1
	// err := db.Create(billingAccount2).Error

	// err := db.Create(&admin1).Error
	err := db.Create(&billingAccount1).Error
	return err

}

var _ = Describe("Repository", func() {
	BeforeEach(func() {
		// repo = &Db.Repository{Db: Db}
		// err := repo.Migrate() // auto create tables
		err := Db.AutoMigrate(&models.BillingAccount{}, &models.AdminDetails{}, &models.BillFile{}, &models.Project{}, &models.Company{})
		insertBillingAccountData(Db)
		Ω(err).To(Succeed())
	})
	It("can Insert test data into the database", func() {
		err := insertBillingAccountData(Db)
		Ω(err).To(Succeed())
	})
	Context("Load", func() {
		It("finds Billing account in database", func() {
			var billingAccount models.BillingAccount
			insertBillingAccountData(Db)
			err := Db.First(&billingAccount).Error

			Ω(err).To(Succeed())
			// Ω(blog.Content).To(Equal("hello"))
			// Ω(blog.Tags).To(Equal(pq.StringArray{"a", "b"}))
		})
		// It("Not Found", func() {
		// 	_, err := repo.Load(999)
		// 	Ω(err).To(HaveOccurred())
		// })
	})
	It("deducts from account under plan model with sufficient balence", func() {
		insertBillingAccountData(Db)

	})
})
