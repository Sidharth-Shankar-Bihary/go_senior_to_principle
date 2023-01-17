package repos

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbTest *gorm.DB
var err error
var mockUserRepo *repo

func init() {
	connStr := "host=localhost port=5432 user=root dbname=pro password= sslmode=disable"
	dbTest, err = gorm.Open(postgres.Open(connStr))

	logger := &zap.Logger{}
	mockUserRepo, err = New(dbTest, logger)
}

func TestRepo_CreateUser(t *testing.T) {
	mockUser := &models.User{}
	mockUser.Username = gofakeit.Username()
	mockUser.Email = gofakeit.Email()
	mockUser.Address = "test"
	mockUser.Password = "12345678"

	result := mockUserRepo.CreateUser(mockUser)
	assert.Equal(t, result, nil)
}

func TestRepo_GetUserByID(t *testing.T) {
	mockUser := &models.User{}
	mockUser.ID = 6

	user, testErr := mockUserRepo.GetUserByID(uint(mockUser.ID))
	assert.Equal(t, testErr, nil)
	assert.Equal(t, user.ID, mockUser.ID)
}

func TestRepo_GetUserByName(t *testing.T) {
	mockUser := &models.User{}

	// test a username does not exist
	mockUser.Username = gofakeit.Username()
	user1, testErr := mockUserRepo.GetUserByName(mockUser.Username)
	assert.Equal(t, testErr, nil)
	assert.Equal(t, user1.Username, "")

	// test a username exists
	mockUser.Username = "wq112"
	user2, testErr := mockUserRepo.GetUserByName(mockUser.Username)
	assert.Equal(t, testErr, nil)
	assert.Equal(t, user2.Username, mockUser.Username)
}

func Test_CleanUp(t *testing.T) {
	dbTest.Exec("delete from users where address=?", "test")
}
