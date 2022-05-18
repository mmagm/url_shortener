package logger

import (
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

// https://github.com/uber-go/zap/blob/master/config.go
func newZapCustomConfig() zap.Config {
	return zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// NewZapCustom builds a custom Logger that writes InfoLevel and above
// logs to standard error in a human-friendly format.
//
// It's a shortcut for newZapCustomConfig().Build(...Option).
func NewZapCustom(options ...zap.Option) (*zap.Logger, error) {
	return newZapCustomConfig().Build(options...)
}

func init() {
	// Zap logger (uber-go): https://github.com/uber-go/zap
	// https://godoc.org/go.uber.org/zap
	// logger, err := zap.NewDevelopment()
	logger, err := NewZapCustom()

	if err != nil {
		println("error in getting the ZAP for production, fallback to EXAMPLE")
		logger = zap.NewExample()
	}

	// In contexts where performance is nice, but not critical, use
	// the SugaredLogger. It's 4-10x faster than other structured
	// logging packages and includes both structured and printf-style
	// APIs.
	log = logger.Sugar()
	log.Infof("LOGGINMIDDLEWARE (logger) initialization")
}

func Logger() (logger *zap.SugaredLogger) {
	return log
}
