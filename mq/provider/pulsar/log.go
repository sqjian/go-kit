package pulsar

import (
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func logWrapper(debug bool) log.Logger {

	var logInst = logrus.New()

	switch debug {
	case true:
		logInst.SetOutput(os.Stdout)
	default:
		logInst.SetOutput(ioutil.Discard)
	}

	logInst.SetLevel(logrus.DebugLevel)
	logInst.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05", //时间格式化
	})
	return log.NewLoggerWithLogrus(logInst)
}
