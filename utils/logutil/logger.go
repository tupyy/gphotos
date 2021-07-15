package logutil

import (
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetLogger(c *gin.Context) *logrus.Entry {
	logger := newLogger(c)

	methodName := getMethodName()

	if methodName != "" {
		return logger.WithField("method", methodName)
	}

	return logger
}

func GetDefaultLogger() *logrus.Entry {
	return logrus.WithField("method", getMethodName())
}

func newLogger(c *gin.Context) *logrus.Entry {
	fields := logrus.Fields{}

	if c.GetString("username") != "" {
		fields["user"] = c.GetString("username")
	}

	if c.GetString("sessionID") != "" {
		fields["sessionID"] = c.GetString("sessionID")
	}

	return logrus.WithFields(fields)
}

func getMethodName() string {
	// 4 as stack depth should be enough to get the real caller. (2 should be enough)
	stack := make([]uintptr, 4)
	depth := runtime.Callers(3, stack) // Can skip the first 3 as it's Callers < getMethodName < Get(*)Logger

	if depth < 1 {
		return ""
	}

	frames := runtime.CallersFrames(stack)

	for f, hasNext := frames.Next(); hasNext; {

		tmp := strings.Split(f.Function, "/")
		if len(tmp) == 0 {
			continue
		}

		shortName := tmp[len(tmp)-1]

		if !strings.HasPrefix(shortName, "logutil.") {
			return shortName
		}
	}

	return ""
}
