package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	postgres "github.com/tupyy/gophoto/internal/utils/pgclient"
)

const (
	prefix = "GPHOTOS"

	gracefulShutdown        = "GRACEFUL_SHUTDOWN"
	defaultGracefulShutdown = 5 * time.Second

	defaultHttpTimeout = 5 * time.Second
)

var (
	configuration Configuration
)

type KeycloakConfig struct {
	ClientID      string `json:"client_id" yaml:"client_id"`
	ClientSecret  string `json:"client_secret" yaml:"client_secret"`
	BaseURL       string `json:"base_url" yaml:"base_url"`
	AdminURL      string `json:"admin_url" yaml:"admin_url"`
	Realm         string `json:"realm" yaml:"realm"`
	AdminUsername string `json:"admin_username" yaml:"admin_username"`
	AdminPwd      string `json:"admin_password" yaml:"admin_password"`
}

func (k KeycloakConfig) String() string {
	kk := KeycloakConfig{
		ClientID:      k.ClientID,
		BaseURL:       k.BaseURL,
		AdminURL:      k.AdminURL,
		Realm:         k.Realm,
		AdminUsername: k.AdminUsername,
		ClientSecret:  shadePassword(k.ClientSecret),
		AdminPwd:      shadePassword(k.AdminPwd),
	}
	j, _ := json.Marshal(kk)
	return string(j)
}

type MinioConfig struct {
	Url             string `json:"url" yaml:"url"`
	AccessID        string `json:"access_id" yaml:"access_id"`
	AccessSecretKey string `json:"access_secret_key" yaml:"access_secret_key"`
}

func (m MinioConfig) String() string {
	mm := MinioConfig{
		Url:             m.Url,
		AccessID:        m.AccessID,
		AccessSecretKey: shadePassword(m.AccessSecretKey),
	}
	j, _ := json.Marshal(mm)
	return string(j)
}

type PostgresConfig struct {
	// params for postgresql.
	Host     string `json:"host" yaml:"host"`
	Port     uint   `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database"`
}

func (p PostgresConfig) String() string {
	pp := PostgresConfig{
		Host:     p.Host,
		Port:     p.Port,
		User:     p.User,
		Database: p.Database,
		Password: shadePassword(p.Password),
	}
	j, _ := json.Marshal(pp)
	return string(j)
}

type Configuration struct {
	LogLevel        string `json:"log_level" yaml:"log_level"`
	AuthCallbackURL string `json:"auth_callback_url" yaml:"auth_callback_url"`
	SecretKey       string `json:"secret_key" yaml:"secret_key"`
	EncryptionKey   string `json:"encryption_key" yaml:"encryption_key"`
	NoAuth          bool   `json:"no_auth" yaml:"no_auth"`

	Keycloak KeycloakConfig `json:"keycloak" yaml:"keycloak"`
	Minio    MinioConfig    `json:"minio" yaml:"minio"`
	Postgres PostgresConfig `json:"postgres" yaml:"postgres"`
}

func (c Configuration) String() string {
	cc := Configuration{
		LogLevel:        c.LogLevel,
		AuthCallbackURL: c.AuthCallbackURL,
		SecretKey:       shadePassword(c.SecretKey),
		EncryptionKey:   shadePassword(c.EncryptionKey),
		NoAuth:          c.NoAuth,
		Postgres: PostgresConfig{
			Host:     c.Postgres.Host,
			Port:     c.Postgres.Port,
			User:     c.Postgres.User,
			Database: c.Postgres.Database,
			Password: shadePassword(c.Postgres.Password),
		},
		Minio: MinioConfig{
			Url:             c.Minio.Url,
			AccessID:        c.Minio.AccessID,
			AccessSecretKey: shadePassword(c.Minio.AccessSecretKey),
		},
		Keycloak: KeycloakConfig{
			ClientID:      c.Keycloak.ClientID,
			BaseURL:       c.Keycloak.BaseURL,
			AdminURL:      c.Keycloak.AdminURL,
			Realm:         c.Keycloak.Realm,
			AdminUsername: c.Keycloak.AdminUsername,
			ClientSecret:  shadePassword(c.Keycloak.ClientSecret),
			AdminPwd:      shadePassword(c.Keycloak.AdminPwd),
		},
	}
	j, _ := json.Marshal(cc)
	return string(j)
}

func ParseConfiguration(confFile string) error {
	setDefaults()

	viper.SetEnvPrefix(prefix)
	viper.AutomaticEnv() // read in environment variables that match

	if len(confFile) == 0 {
		logrus.Info("no config file specified")
		return errors.New("no config file specified")
	}

	content, err := os.ReadFile(confFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, &configuration); err != nil {
		return err
	}
	logrus.Infof("using config file: %v", confFile)
	return nil
}

func setDefaults() {
	viper.SetDefault(gracefulShutdown, defaultGracefulShutdown)
}

func GetLogLevel() logrus.Level {
	switch configuration.LogLevel {
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

	return logrus.InfoLevel
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

func GetConfiguration() Configuration {
	return configuration
}

func GetKeycloakConfig() KeycloakConfig {
	return configuration.Keycloak
}

func GetMinioConfig() MinioConfig {
	return configuration.Minio
}

func GetPostgresConf() postgres.ClientParams {
	ret := postgres.ClientParams{
		Host:     configuration.Postgres.Host,
		Port:     configuration.Postgres.Port,
		User:     configuration.Postgres.User,
		Password: configuration.Postgres.Password,
		DBName:   configuration.Postgres.Database,
	}

	logrus.Infof("postgres conf: %+v", ret)

	return ret
}

func GetServerSecretKey() string {
	return configuration.SecretKey
}

func GetEncryptionKey() string {
	return configuration.EncryptionKey
}

// GetServerAuthCallback returns the url of the authentication callback.
func GetServerAuthCallback() string {
	return configuration.AuthCallbackURL
}

func GetStaticsFolder() string {
	return ""
}

func HasAuth() bool {
	return configuration.NoAuth
}
func shadePassword(pwd string) string {
	offset := len(pwd) - 3
	if offset > 0 {
		return fmt.Sprintf("******%s", pwd[offset:])
	}
	return "******"
}
