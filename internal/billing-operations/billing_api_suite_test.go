package billingoperations

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "billing-api/utils"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBillingApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BillingApi Suite")
}

// func TestDocker(t *testing.T) {
// 	RegisterFailHandler(Fail)
// 	RunSpecs(t, "Docker Suite")
// }

var Db *gorm.DB
var cleanupDocker func()

var _ = BeforeSuite(func() {
	// setup *gorm.Db with docker
	Db, cleanupDocker = setupGormWithDocker()
})

var _ = AfterSuite(func() {
	// cleanup resource
	cleanupDocker()
})

var _ = BeforeEach(func() {
	// clear db tables before each test
	err := Db.Exec(`DROP SCHEMA public CASCADE;CREATE SCHEMA public;`).Error
	Ω(err).To(Succeed())
})

const (
	dbName = "test"
	// passwd = "test"
	passwd = "mypasswordbob"
)

func setupGormWithDocker() (*gorm.DB, func()) {
	pool, err := dockertest.NewPool("")
	chk(err)

	runDockerOpt := &dockertest.RunOptions{
		Repository: "postgres", // image
		Tag:        "14",       // version
		Env:        []string{"POSTGRES_PASSWORD=" + passwd, "POSTGRES_DB=" + dbName},
	}

	fnConfig := func(config *docker.HostConfig) {
		config.AutoRemove = true                     // set AutoRemove to true so that stopped container goes away by itself
		config.RestartPolicy = docker.NeverRestart() // don't restart container
	}

	resource, err := pool.RunWithOptions(runDockerOpt, fnConfig)
	chk(err)
	// call clean up function to release resource
	fnCleanup := func() {
		err := resource.Close()
		chk(err)
	}

	// conStr := fmt.Sprintf("host=localhost port=%s user=postgres dbname=%s password=%s sslmode=disable",
	// conStr := fmt.Sprintf("host=172.17.0.3 port=%s user=bobsql dbname=%s password=%s sslmode=disable",
	// 	resource.GetPort("5432/tcp"), // get port of localhost
	// 	dbName,
	// 	passwd,
	// )

	// conStr := "host=" + os.Getenv("Postgresqlhost") + " user=" + os.Getenv("Postgresqluser") + " password=" + os.Getenv("Postgresqlpassword") + " dbname=" + os.Getenv("Postgresqldbname") + " port=" + os.Getenv("Postgresqlport")
	conStr := "host=172.17.0.3 port=5432 user=bobsql password=mypasswordbob dbname=test"
	var gdb *gorm.DB
	// retry until db server is ready
	err = pool.Retry(func() error {
		// gdb, err = gorm.Open(postgres.Open(conStr))
		gdb, err = gorm.Open(postgres.Open(conStr), &gorm.Config{})

		if err != nil {
			return err
		}
		db, err := gdb.DB()
		if err != nil {
			return err
		}
		return db.Ping()
	})
	chk(err)

	// container is ready, return *gorm.Db for testing
	return gdb, fnCleanup
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

// var _ = Describe("Repository", func() {
// 	It("can Insert test data into the database", func() {
// 		err := insertBillingAccountData(Db)
// 		logging.Fatal(logging.HTTPError, err)
// 		Ω(err).To(Succeed())
// 	})
// })
