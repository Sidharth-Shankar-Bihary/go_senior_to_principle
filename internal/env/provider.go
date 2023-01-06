package env

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/ramseyjiang/go_senior_to_principle/internal/db"
	"github.com/ramseyjiang/go_senior_to_principle/pkg/baseenv"
	"github.com/ramseyjiang/go_senior_to_principle/pkg/utils"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

func provideLogger() (logger *zap.Logger) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:       "T",
		LevelKey:      "L",
		NameKey:       "N",
		CallerKey:     "C",
		MessageKey:    "M",
		StacktraceKey: "S",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(utils.FmtDateTimeStr))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewJSONEncoder(cfg.EncoderConfig)

	var core zapcore.Core
	if os.Getenv("env") == "test" {
		core = zapcore.NewCore(encoder, os.Stdout, zap.InfoLevel)
	} else {
		core = zapcore.NewCore(encoder, os.Stdout, zap.DebugLevel)
	}

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
}

func providerTracer() (io.Closer, error) {
	var cfg = jaegercfg.Configuration{
		ServiceName: os.Getenv("JAEGER_TRACE_NAME"),
		Sampler: &jaegercfg.SamplerConfig{
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: os.Getenv("JAEGER_LOCAL_AGENT_HOST_PORT"),
		},
		Tags: []opentracing.Tag{{Key: "env", Value: os.Getenv("ENV")}, {Key: "projectEnv", Value: os.Getenv("JAEGER_PROJECT_ENV")}},
	}

	jMetricsFactory := metrics.NullFactory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Metrics(jMetricsFactory),
	)
	opentracing.SetGlobalTracer(tracer)
	return closer, err
}

func provideBEnv(ctx context.Context, logger *zap.Logger) (env *baseenv.Environment) {
	return baseenv.NewBaseEnv(ctx, logger)
}

func providerRedis(ctx context.Context) (client *redis.Client, err error) {
	// NewClient returns a client to the Redis Server specified by Options. Password and DB can be removed, it is up to you.
	client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"), // use default Addr
		Password: "",                                                      // no password set
	})

	// Check the connection
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return
}

func provideConnDB() (database *gorm.DB, err error) {
	err = db.InitDB()
	if err != nil {
		return nil, err
	}
	database = db.GetDB()
	return database, nil
}
