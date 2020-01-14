package l

import "go.uber.org/zap"

var Logger *zap.Logger

// Init initialized the logging tool.
func Init(env string) {
	if env == "production" {
		Logger, _ = zap.NewProduction()
	} else {
		Logger, _ = zap.NewDevelopment()
	}
}
