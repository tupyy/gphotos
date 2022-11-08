package conf

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tupyy/gophoto/internal/utils/logutil"
	postgres "github.com/tupyy/gophoto/internal/utils/pgclient"
)

const (
	// server params
	authCallbackURL = "AUTH_CALLBACK_URL"
	secretKey       = "SECRET_KEY"
	encryptionKey   = "ENCRYPTION_KEY"
	noAuth          = "AUTH_DISABLED"

	// cache config for repo.
	repoCacheTTL           = "REPOCACHE_TTL"
	repoCacheCleanInterval = "REPOCACHE_CLEAN_INTERVAL"

	// params for postgresql.
	pgsqlHost   = "POSTGRESQL_HOST"
	pgsqlPort   = "POSTGRESQL_PORT"
	pgsqlUser   = "POSTGRESQL_USER"
	pgsqlPwd    = "POSTGRESQL_PASSWORD"
	pgsqlDBName = "POSTGRESQL_DBNAME"

	// folders
	templateDir = "TEMPLATE_FOLDER"
	staticDir   = "STATICS_FOLDER"
)

func GetPostgresConf() postgres.ClientParams {
	ret := postgres.ClientParams{
		Host:     viper.GetString(pgsqlHost),
		Port:     viper.GetUint(pgsqlPort),
		User:     viper.GetString(pgsqlUser),
		Password: viper.GetString(pgsqlPwd),
		DBName:   viper.GetString(pgsqlDBName),
	}

	logrus.Infof("postgres conf: %+v", ret)

	return ret
}

func GetServerSecretKey() string {
	return viper.GetString(secretKey)
}

func GetEncryptionKey() string {
	return viper.GetString(encryptionKey)
}

// GetServerAuthCallback returns the url of the authentication callback.
func GetServerAuthCallback() string {
	return viper.GetString(authCallbackURL)
}

func GetTemplateFolder() string {
	return viper.GetString(templateDir)
}

func GetStaticsFolder() string {
	return viper.GetString(staticDir)
}

func HasAuth() bool {
	if viper.IsSet(noAuth) {
		return false
	}
	return true
}

func GetRepoCacheConfig() (ttl time.Duration, interval time.Duration) {
	logger := logutil.GetDefaultLogger()

	ttl = viper.GetDuration(repoCacheTTL)
	interval = viper.GetDuration(repoCacheCleanInterval)

	logger.Infof("repo cache config is: ttl=%s cleanInterval=%s", ttl, interval)

	return ttl, interval
}
