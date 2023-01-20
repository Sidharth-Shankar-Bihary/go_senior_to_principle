package models

import (
	"log"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbTest *gorm.DB
var err error

func init() {
	connStr := "host=localhost port=5432 user=root dbname=pro password= sslmode=disable"
	dbTest, err = gorm.Open(postgres.Open(connStr))
	if err = dbTest.AutoMigrate(
		&User{},
	); err != nil {
		log.Fatal(err.Error())
	}
}

func TestUser_CreateUser(t *testing.T) {
	mockUser := &User{}
	mockUser.Username = "test123"
	mockUser.Email = gofakeit.Email()
	mockUser.Address = "test"
	mockUser.Password = "12345678"

	err = mockUser.CreateUser(mockUser, dbTest)
	assert.Equal(t, err, nil)
}

func TestUser_GetUserByID(t *testing.T) {
	mockUser := &User{}
	mockUser.ID = 6

	user, testErr := mockUser.GetUserByID(mockUser.ID, dbTest)
	assert.Equal(t, testErr, nil)
	assert.Equal(t, user.ID, mockUser.ID)
}

func TestUser_GetUserByName(t *testing.T) {
	mockUser := &User{}

	// test a username does not exist
	mockUser.Username = gofakeit.Username()
	user1, testErr := mockUser.GetUserByName(mockUser.Username, dbTest)
	assert.Equal(t, testErr, nil)
	assert.Equal(t, user1.Username, "")

	// test a username exists
	mockUser.Username = "test123"
	user2, testErr := mockUser.GetUserByName(mockUser.Username, dbTest)
	assert.Equal(t, user2.Username, mockUser.Username)
	assert.Equal(t, testErr, nil)
}

func Test_CleanUp(t *testing.T) {
	dbTest.Exec("delete from users where address=?", "test")
}
