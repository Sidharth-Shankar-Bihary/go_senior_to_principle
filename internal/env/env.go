package env

import (
	"context"

	"github.com/ramseyjiang/go_senior_to_principle/pkg/apierrors"
	"github.com/ramseyjiang/go_senior_to_principle/pkg/baseenv"
	"go.uber.org/zap"
)

type Environment struct {
	*baseenv.Environment
	rootCancel context.CancelFunc
	cfg        *Conf
	StrErr     *apierrors.StrErrFace
	Log        *zap.Logger
}

// NewEnv can be used to add more env init, such as kafka, grpc, and so on.
func NewEnv(
	environment *baseenv.Environment,
	cancel context.CancelFunc,
	cfg *Conf,
	strErr *apierrors.StrErrFace,
	logger *zap.Logger,
) *Environment {
	return &Environment{
		rootCancel:  cancel,
		Environment: environment,
		cfg:         cfg,
		StrErr:      strErr,
		Log:         logger,
	}
}

func InitEnv(ctx context.Context, cancel context.CancelFunc) (gEnv *Environment, err error) {
	conf, err := provideConf("config.yaml")
	if err != nil {
		return nil, err
	}

	closer, err := providerTracer(conf)
	if err != nil {
		return nil, err
	}

	logger := provideLogger(conf)
	environment := provideBEnv(ctx, logger)
	strErr, err := provideErrorHandler()
	if err != nil {
		return nil, err
	}

	gEnv = NewEnv(environment, cancel, conf, strErr, logger)
	gEnv.AddCloser(closer)

	return gEnv, nil
}

func (env *Environment) C() *Conf {
	return env.cfg
}

func (env *Environment) L() *zap.Logger {
	return env.Log
}

func (env *Environment) Err() *apierrors.StrErrFace {
	return env.StrErr
}
