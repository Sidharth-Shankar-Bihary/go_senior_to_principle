package repos

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/models"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbTest       *gorm.DB
	err          error
	mockUserRepo *repo
)

func TestUserRepo(t *testing.T) {
	// RegisterFailHandler is from the package gomega, it is used to connect Ginkgo with Gomega.
	// When a matcher fails the fail handler passed into RegisterFailHandler is called.
	RegisterFailHandler(Fail)
	// RunSpecs is from the package ginkgo, RunSpecs is the entry point for the Ginkgo test runner
	RunSpecs(t, "User Suite")
}

// BeforeSuite is from the package ginkgo
var _ = BeforeSuite(func() {
	// Init container, open connection, run migrations seed database, init repository
	rand.Seed(time.Now().UnixNano())
	connStr := "host=localhost port=5432 user=root dbname=pro password= sslmode=disable"
	dbTest, err = gorm.Open(postgres.Open(connStr))
	if err = dbTest.AutoMigrate([]interface{}{
		&models.User{},
	}); err != nil {
		log.Fatal(err.Error())
	}

	logger := &zap.Logger{}
	mockUserRepo, err = New(dbTest, logger)
})

var _ = Describe("Test the User repository", func() {
	Context("Test Create User", func() {
		When("A user info has been verified and filled, the request has been sent.", func() {
			It("should return nil", func() {
				mockUser := &models.User{}
				mockUser.Username = "test"
				mockUser.Email = gofakeit.Email()
				mockUser.Address = "test"
				mockUser.Password = "12345678"
				result := mockUserRepo.CreateUser(mockUser)

				Expect(result).Should(BeNil())
			})
		})
	})

	Context("Test GetUserByID", func() {
		When("when userID is 6", func() {
			It("should return user info", func() {
				mockUser := &models.User{}
				mockUser.ID = 6

				user, testErr := mockUserRepo.GetUserByID(uint(mockUser.ID))

				Expect(testErr).Should(BeNil())
				Expect(user.ID).Should(BeEquivalentTo(mockUser.ID))
			})
		})
	})

	Context("Test GetUserByName", func() {
		When("when a username does not exist", func() {
			It("should return nil", func() {
				mockUser := &models.User{}
				// test a username does not exist
				mockUser.Username = gofakeit.Username()
				user, testErr := mockUserRepo.GetUserByName(mockUser.Username)

				Expect(testErr).Should(BeNil())
				Expect(user.Username).Should(BeEquivalentTo(""))
			})
		})
	})

	Context("Test GetUserByName", func() {
		When("when a username exists", func() {
			It("should return nil", func() {
				mockUser := &models.User{}
				// test a username does not exist
				mockUser.Username = "test"
				user, testErr := mockUserRepo.GetUserByName(mockUser.Username)

				Expect(testErr).Should(BeNil())
				Expect(user.Username).Should(BeEquivalentTo(mockUser.Username))
			})
		})
	})
})

// Purge function destroys container
var _ = AfterSuite(func() {
	dbTest.Exec("delete from users where address=?", "test")
})
