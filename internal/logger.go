package internal

import (
	"go.uber.org/zap"
	"log"
)

var (
	// logger
	Sugar *zap.SugaredLogger
)
// use zap from uber as logger
// https://github.com/uber-go/zap/blob/master/example_test.go
// https://blog.sandipb.net/2018/05/02/using-zap-simple-use-cases/
// https://github.com/uber-go/zap/issues/261
func InitLogger() {

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	defer logger.Sync()
	Sugar = logger.Sugar()
}
