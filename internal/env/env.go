package env

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/ramseyjiang/go_senior_to_principle/pkg/baseenv"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Environment struct {
	*baseenv.Environment
	rootCancel context.CancelFunc
	Log        *zap.Logger
	Redis      *redis.Client
	DB         *gorm.DB
}

// NewEnv can be used to add more env init, such as redis, kafka, grpc, and so on.
func NewEnv(
	environment *baseenv.Environment,
	cancel context.CancelFunc,
	logger *zap.Logger,
	rds *redis.Client,
	db *gorm.DB,
) *Environment {
	return &Environment{
		rootCancel:  cancel,
		Environment: environment,
		Log:         logger,
		Redis:       rds,
		DB:          db,
	}
}

func InitEnv(ctx context.Context, cancel context.CancelFunc) (gEnv *Environment, err error) {
	err = godotenv.Load() // It reads .env file only and directly.
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	closer, err := providerTracer()
	if err != nil {
		return nil, err
	}

	logger := provideLogger()
	environment := provideBEnv(ctx, logger)

	rds, err := providerRedis(ctx)
	if err != nil {
		return nil, err
	}

	db, err := provideConnDB()
	if err != nil {
		return nil, err
	}

	gEnv = NewEnv(environment, cancel, logger, rds, db)
	gEnv.AddCloser(closer)

	return gEnv, nil
}

func (env *Environment) L() *zap.Logger {
	return env.Log
}
