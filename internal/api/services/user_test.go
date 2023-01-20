package services

import (
	"log"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/models"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/repos"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbTest *gorm.DB
var err error
var mockUserService UserService

func init() {
	connStr := "host=localhost port=5432 user=root dbname=pro password= sslmode=disable"
	dbTest, err = gorm.Open(postgres.Open(connStr))
	if err = dbTest.AutoMigrate([]interface{}{
		&models.User{},
	}); err != nil {
		log.Fatal(err.Error())
	}

	logger := &zap.Logger{}
	mockUserRepo, _ := repos.New(dbTest, logger)
	mockUserService = NewUserService(mockUserRepo, logger)
}

func TestUserService_CreateUser(t *testing.T) {
	mockReq := &RegisterRequest{
		Username: gofakeit.Username(),
		Password: "12345678",
		Email:    "test@qq.com",
	}

	resp, testErr := mockUserService.CreateUser(*mockReq)
	assert.Equal(t, testErr, nil)
	assert.Equal(t, resp.Err, nil)
	assert.Equal(t, resp.Status, http.StatusOK)
}

func TestUserService_GetUserByID(t *testing.T) {
	mockReq := &GetUserRequest{ID: 6}

	resp, testErr := mockUserService.GetUserByID(*mockReq)
	assert.Equal(t, testErr, nil)
	assert.Equal(t, resp.Err, "")
	assert.Equal(t, resp.User.ID, mockReq.ID)
}

func Test_CleanUp(t *testing.T) {
	dbTest.Exec("delete from users where email=?", "test@qq.com")
}
