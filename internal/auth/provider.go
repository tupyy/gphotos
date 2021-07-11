package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/tupyy/gophoto/internal/conf"
	"golang.org/x/oauth2"
)

type OidcProvider struct {
	provider *oidc.Provider
	Config   oauth2.Config
	Issuer   string
}

func NewOidcProvider(conf conf.KeycloakConfig, authCallBack string) *OidcProvider {
	issuer := fmt.Sprintf("%s/auth/realms/%s", conf.BaseURL, conf.Realm)

	provider, err := oidc.NewProvider(context.Background(), issuer)
	if err != nil {
		log.Fatal(err)
	}

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		RedirectURL:  authCallBack,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	return &OidcProvider{provider, oauth2Config, issuer}
}

func (p *OidcProvider) Verifier() *oidc.IDTokenVerifier {
	config := p.Config
	oidcConfig := &oidc.Config{
		ClientID: config.ClientID,
	}

	return p.provider.Verifier(oidcConfig)
}

func (p *OidcProvider) GetLogoutURL() (string, error) {
	var emptyString string

	wellKnown := strings.TrimSuffix(p.Issuer, "/") + "/.well-known/openid-configuration"

	req, err := http.NewRequest("GET", wellKnown, nil)
	if err != nil {
		return emptyString, err
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return emptyString, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return emptyString, fmt.Errorf("unable to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return emptyString, fmt.Errorf("%s: %s", resp.Status, body)
	}

	var l = struct {
		LogoutURL string `json:"end_session_endpoint"`
	}{}

	err = json.Unmarshal(body, &l)
	if err != nil {
		return emptyString, fmt.Errorf("oidc: failed to decode provider discovery object: %v", err)
	}

	return l.LogoutURL, nil
}
