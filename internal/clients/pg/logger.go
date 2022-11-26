package pg

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

const SlowQueryThreshold = 250 * time.Millisecond

type driverLogger struct{}

// Don't overwrite logrus log level.
func (d driverLogger) LogMode(logger.LogLevel) logger.Interface {
	return d
}

func (d driverLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	zap.S().Infow(msg, args...)
}

func (d driverLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	zap.S().Warnw(msg, args...)
}

func (d driverLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	zap.S().Errorw(msg, args...)
}

func (d driverLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		zap.S().Errorw("error with sql query", "error", err, "rows", rows)
		return
	}

	if elapsed > SlowQueryThreshold {
		zap.S().Warnw("slow sql query", "query", sql, "rows", rows)
		return
	}

	zap.S().Debugw("SQL query", "query", sql, "rows", rows)
}
