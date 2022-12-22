package env

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/ramseyjiang/go_senior_to_principle/pkg/apierrors"
	"github.com/ramseyjiang/go_senior_to_principle/pkg/baseenv"
	"github.com/ramseyjiang/go_senior_to_principle/pkg/utils"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func provideLogger(c *Conf) (logger *zap.Logger) {
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
	core := zapcore.NewCore(encoder, os.Stdout, c.Logger.Level)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
}

func providerTracer(c *Conf) (io.Closer, error) {
	var cfg = jaegercfg.Configuration{
		ServiceName: c.Jaeger.TracerName,
		Sampler: &jaegercfg.SamplerConfig{
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: c.Jaeger.LocalAgentHostPort,
		},
		Tags: []opentracing.Tag{{Key: "env", Value: c.Env}, {Key: "projectEnv", Value: c.Jaeger.ProjectEnv}},
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

func provideErrorHandler() (strErr *apierrors.StrErrFace, err error) {
	strErr = apierrors.NewStrErrFace()
	// customer error code
	err = apierrors.SetCustomizeErr(strErr)
	if err != nil {
		return nil, err
	}
	return strErr, nil
}
