package repos

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// repo is a struct name, and here does not use user as a struct name.
// The reason is this way can be used in many repos. If you use a user as a struct name, other repos should change the struct name all the time.
// repo is used to connect to all data source, such as models, redis, kafka.
type repo struct {
	db     *gorm.DB
	logger *zap.Logger
}

func New(db *gorm.DB, logger *zap.Logger) (*repo, error) {
	return &repo{
		db:     db,
		logger: logger,
	}, nil
}
