package conf

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tupyy/gophoto/utils/logutil"
	postgres "github.com/tupyy/gophoto/utils/pgclient"
)

// Below all the different keys used to configure this service.
const (
	prefix   = "GOPHOTO"
	logLevel = "LOG_LEVEL"

	gracefulShutdown        = "GRACEFUL_SHUTDOWN"
	defaultGracefulShutdown = 5 * time.Second

	defaultHttpTimeout = 5 * time.Second

	// folders
	templateDir = "TEMPLATE_FOLDER"
	staticDir   = "STATICS_FOLDER"

	// params for keycloak
	keycloakClientID     = "KEYCLOAK_CLIENT_ID"
	keycloakClientSecret = "KEYCLOAK_CLIENT_SECRET"
	keycloakBaseURL      = "KEYCLOAK_URL"
	keycloakRealm        = "KEYCLOAK_REALM"
	keycloakAdmin        = "KEYCLOAK_ADMIN_USERNAME"
	keycloakAdminPwd     = "KEYCLOAK_ADMIN_PWD"

	// server params
	authCallbackURL = "AUTH_CALLBACK_URL"
	secretKey       = "SECRET_KEY"
	encryptionKey   = "ENCRYPTION_KEY"

	// cache config for repo.
	repoCacheTTL           = "REPOCACHE_TTL"
	repoCacheCleanInterval = "REPOCACHE_CLEAN_INTERVAL"

	// params for postgresql.
	pgsqlHost   = "POSTGRESQL_HOST"
	pgsqlPort   = "POSTGRESQL_PORT"
	pgsqlUser   = "POSTGRESQL_USER"
	pgsqlPwd    = "POSTGRESQL_PASSWORD"
	pgsqlDBName = "POSTGRESQL_DBNAME"

	// params for minio
	minioUrl  = "MINIO_SERVER_URL"
	minioUser = "MINIO_ACCESS_ID"
	minioPwd  = "MINIO_ACCESS_KEY"
)

type KeycloakConfig struct {
	ClientID      string
	ClientSecret  string
	BaseURL       string
	Realm         string
	AdminUsername string
	AdminPwd      string
}

type MinioConfig struct {
	Url      string
	User     string
	Password string
}

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

func GetKeycloakConfig() KeycloakConfig {
	return KeycloakConfig{
		ClientID:      viper.GetString(keycloakClientID),
		ClientSecret:  viper.GetString(keycloakClientSecret),
		BaseURL:       viper.GetString(keycloakBaseURL),
		Realm:         viper.GetString(keycloakRealm),
		AdminUsername: viper.GetString(keycloakAdmin),
		AdminPwd:      viper.GetString(keycloakAdminPwd),
	}
}

func (s KeycloakConfig) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "ClientID: %s\n", s.ClientID)
	fmt.Fprintf(&sb, "ClientSecret: %v\n", len(s.ClientSecret) != 0)
	fmt.Fprintf(&sb, "BaseURL: %s\n", s.BaseURL)
	fmt.Fprintf(&sb, "Realm: %s\n", s.Realm)

	return sb.String()
}

func GetMinioConfig() MinioConfig {
	return MinioConfig{
		Url:      viper.GetString(minioUrl),
		User:     viper.GetString(minioUser),
		Password: viper.GetString(minioPwd),
	}
}

func (m MinioConfig) String() string {
	return fmt.Sprintf("Url: %s, User: %s, Password: %v", m.Url, m.User, len(m.Password) > 0)
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

func GetGracefulShutdownDuration() time.Duration {
	if !viper.IsSet(gracefulShutdown) {
		return defaultGracefulShutdown
	}

	return viper.GetDuration(gracefulShutdown)
}

func GetHttpRequestTimeout() time.Duration {
	return defaultHttpTimeout
}

func GetRepoCacheConfig() (ttl time.Duration, interval time.Duration) {
	logger := logutil.GetDefaultLogger()

	ttl = viper.GetDuration(repoCacheTTL)
	interval = viper.GetDuration(repoCacheCleanInterval)

	logger.Infof("repo cache config is: ttl=%s cleanInterval=%s", ttl, interval)

	return ttl, interval
}
