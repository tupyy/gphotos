package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

const SlowQueryThreshold = 250 * time.Millisecond

type driverLogger struct{}

// Don't overwrite logrus log level.
func (d driverLogger) LogMode(logger.LogLevel) logger.Interface {
	return d
}

func (d driverLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	logrus.Infof(msg, args...)
}

func (d driverLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	logrus.Warnf(msg, args...)
}

func (d driverLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	logrus.Errorf(msg, args...)
}

func (d driverLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	l := logrus.
		WithField("elapsed", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)).
		WithField("rows", rows)

	if err != nil {
		l.WithError(err).Errorf("error with SQL query: %v", sql)
		return
	}

	if elapsed > SlowQueryThreshold {
		l.Warnf("slow SQL query: %v", sql)
		return
	}

	l.Debugf("SQL query: %v", sql)
}
