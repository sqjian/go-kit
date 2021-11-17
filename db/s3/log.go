package s3

import (
	"fmt"
	"github.com/aws/smithy-go/logging"
	"time"
)

type DefLogger struct {
	dummy bool
}

func (l *DefLogger) Logf(classification logging.Classification, format string, v ...interface{}) {
	if l.dummy {
		return
	}
	fmt.Printf("%v => level:%v,msg:", time.Now().UnixNano(), classification)
	fmt.Printf(format, v...)
}
