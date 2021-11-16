package s3

import (
	"github.com/aws/smithy-go/logging"
	"log"
)

type DefLogger struct {
	dummy bool
}

func (l *DefLogger) Logf(classification logging.Classification, format string, v ...interface{}) {
	if l.dummy {
		return
	}
	log.Println(classification)
	log.Printf(format, v...)
}
