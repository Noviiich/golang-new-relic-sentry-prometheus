package config

import (
	"os"

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrzap"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ZapConfig(app *newrelic.Application) *zap.Logger {
	core := zapcore.NewCore(zapcore.NewJSONEncoder(
		zap.NewProductionEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zap.InfoLevel, // Info Level для продакшена
	)

	backgroundCore, err := nrzap.WrapBackgroundCore(core, app)
	if err != nil && err != nrzap.ErrNilApp {
		panic(err)
	}

	return zap.New(
		backgroundCore,
		zap.AddCaller(),                       // добавляет информацию о месте вызова
		zap.AddStacktrace(zapcore.ErrorLevel), // добавляет stacktrace для ошибок
	)
}

func ZapTestConfig() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err.Error())
	}
	return logger
}
