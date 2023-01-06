package repos

import (
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/models"
	"go.uber.org/zap"
)

var u = models.User{}

func (r *repo) CreateUser(user *models.User) (err error) {
	if err = u.CreateUser(user, r.db); err != nil {
		r.logger.Debug("CreateUser ", zap.Any("error: ", err))
	}

	return nil
}

func (r *repo) GetUserByID(id uint) (*models.User, error) {
	user, err := u.GetUserByID(uint64(id), r.db)
	if err != nil {
		r.logger.Debug("GetUserByID ", zap.Any("error: ", err))
		return nil, err
	}

	return user, nil
}

func (r *repo) GetUserByName(username string) (*models.User, error) {
	user, err := u.GetUserByName(username, r.db)
	if err != nil {
		r.logger.Debug("GetUserByName ", zap.Any("error: ", err))
		return nil, err
	}

	return user, nil
}
