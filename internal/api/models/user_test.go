package models

import (
	"strconv"
	"testing"
	"time"

	"github.com/ramseyjiang/go_senior_to_principle/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestUserModelGet(t *testing.T) {
	userModel := NewUser()
	uid := uint64(0) // When connect to db, should use uid gather than 0, and test more about first name, and so on.
	userInfo, err := userModel.Get(uid)
	expected := map[string]string{
		"ID":        "",
		"CreatedAt": "0001-01-01 00:00:00 +0000 UTC",
		"UpdatedAt": "0001-01-01 00:00:00 +0000 UTC",
	}

	id, _ := strconv.ParseUint(expected["ID"], 10, 32)
	createAt, _ := time.Parse(utils.FmtDate, expected["CreatedAt"])
	updateAt, _ := time.Parse(utils.FmtDate, expected["UpdatedAt"])

	assert.Nil(t, err)
	assert.Equal(t, id, userInfo.ID)
	assert.Equal(t, createAt, userInfo.CreatedAt)
	assert.Equal(t, updateAt, userInfo.UpdatedAt)
}
