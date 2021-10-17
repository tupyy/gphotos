package conf

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	// params for keycloak
	keycloakClientID     = "KEYCLOAK_CLIENT_ID"
	keycloakClientSecret = "KEYCLOAK_CLIENT_SECRET"
	keycloakBaseURL      = "KEYCLOAK_AUTH_URL"
	keycloakBackendURL =  "KEYCLOAK_BACKEND_URL"
	keycloakRealm        = "KEYCLOAK_REALM"
	keycloakAdmin        = "KEYCLOAK_ADMIN_USERNAME"
	keycloakAdminPwd     = "KEYCLOAK_ADMIN_PWD"
)

type KeycloakConfig struct {
	ClientID      string
	ClientSecret  string
	BaseURL       string
	Realm         string
	AdminUsername string
	AdminPwd      string
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
