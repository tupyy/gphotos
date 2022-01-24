package conf

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Below all the different keys used to configure this service.
const (
	prefix       = "GOPHOTO"
	logLevel     = "LOG_LEVEL"
	logFormatter = "LOG_FORMATTER"

	gracefulShutdown        = "GRACEFUL_SHUTDOWN"
	defaultGracefulShutdown = 5 * time.Second

	defaultHttpTimeout = 5 * time.Second
)

func ParseConfiguration(confFile string) {
	setDefaults()

	viper.SetEnvPrefix(prefix)
	viper.AutomaticEnv() // read in environment variables that match

	if len(confFile) == 0 {
		logrus.Info("no config file specified")
		return
	}

	viper.SetConfigFile(confFile)

	err := viper.ReadInConfig()
	if err != nil {
		logrus.WithError(err).Errorf("failed to read config file %v", confFile)
		return
	}

	logrus.Infof("using config file: %v", viper.ConfigFileUsed())
}

func setDefaults() {
	viper.SetDefault(gracefulShutdown, defaultGracefulShutdown)

	// repo cache
	viper.SetDefault(repoCacheTTL, "1h")
	viper.SetDefault(repoCacheCleanInterval, "6h")
}

func (m MinioConfig) String() string {
	return fmt.Sprintf("Url: %s, User: %s, Password: %v", m.Url, m.User, len(m.Password) > 0)
}

func GetLogLevel() logrus.Level {
	if !viper.IsSet(logLevel) {
		return logrus.WarnLevel
	}

	level := viper.GetString(logLevel)
	switch level {
	case "TRACE":
		return logrus.TraceLevel
	case "DEBUG":
		return logrus.DebugLevel
	case "INFO":
		return logrus.InfoLevel
	case "WARN":
		return logrus.WarnLevel
	case "ERROR":
		return logrus.ErrorLevel
	}

	return logrus.WarnLevel
}

func GetLogFormatter() logrus.Formatter {
	switch strings.ToLower(viper.GetString(logFormatter)) {
	case "json":
		return &logrus.JSONFormatter{}
	default:
		return &logrus.TextFormatter{}
	}
}

func GetGracefulShutdownDuration() time.Duration {
	if !viper.IsSet(gracefulShutdown) {
		return defaultGracefulShutdown
	}

	return viper.GetDuration(gracefulShutdown)
}

func GetHttpRequestTimeout() time.Duration {
	return defaultHttpTimeout
}
