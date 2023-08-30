package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(productionENV string) (*zap.Logger, error) {
	if productionENV == "production" {
		return NewProduction()
	}

	return NewConsole()
}

func NewProduction() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.DisableCaller = true
	config.DisableStacktrace = true

	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	return config.Build()
}

func NewConsole() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.DisableCaller = true
	config.DisableStacktrace = true

	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	config.Encoding = "console"

	return config.Build()
}
