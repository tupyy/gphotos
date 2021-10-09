package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Nerzal/gocloak/v8"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/utils/logutil"
	"golang.org/x/oauth2"
)

const (
	sessionID = "sessionID"
)

var (
	errSessionExpired = errors.New("session expired")
)

// Authenticator is the high level interface that authenticates a request based on a header given as parameter.
type Authenticator interface {
	// AuthMiddleware is the authentication middleware.
	AuthMiddleware() gin.HandlerFunc
	// Callback return an endpoint which will be called by keycloak after a successful authentication.
	Callback() gin.HandlerFunc
	// Logout logs out the user from keycloak only. It is up to the controller to clean up any remaining sessions.
	Logout(c *gin.Context, username, refreshToken string) error
}

type keyCloakAuthenticator struct {
	oidcProvider *oidcProvider
	client       gocloak.GoCloak
	conf         conf.KeycloakConfig
}

func NewKeyCloakAuthenticator(c conf.KeycloakConfig, authCallback string) Authenticator {

	// initialize oidc provier
	oidcProvider := newOidcProvider(c, authCallback)
	keycloakClient := gocloak.NewClient(c.BaseURL)

	return &keyCloakAuthenticator{oidcProvider: oidcProvider, client: keycloakClient, conf: c}
}

// Middleware returns the authentication middleware used for private routes.
func (k *keyCloakAuthenticator) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		cookie, err := c.Request.Cookie(sessionID)
		if err != nil {
			c.Abort()
			redirectToLogin(c, k.oidcProvider.Config)
			return
		}

		logger := logutil.GetLogger(c)

		s := session.Get(cookie.Value)
		if s == nil {
			logger.WithField("sessionID", cookie.Value).Warn("no session found with this id")
			c.Abort()
			redirectToLogin(c, k.oidcProvider.Config)
			return
		}

		logger.WithField("sessionID", cookie.Value).Debug("new request with session id")
		sessionData, _ := s.(entity.Session)

		if err := k.authenticate(c, sessionData); err != nil {
			logger.WithError(err).Debug("failed to authenticate")
			c.Abort()
			redirectToLogin(c, k.oidcProvider.Config)
			return
		}

		session.Set(cookie.Value, sessionData)
		session.Save()

		c.Next()
	}
}

// Callback returns a handler called after a successful authentication.
func (k *keyCloakAuthenticator) Callback() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		// generate a session ID
		uuid := uuid.New()

		logger := logutil.GetLogger(c)

		logger.WithField("uuid", uuid.String()).Info("session created")

		state := session.Get("state")
		if state == nil {
			http.Error(c.Writer, "state not found", http.StatusBadRequest)
			return
		}

		if c.Request.URL.Query().Get("state") != state.(string) {
			http.Error(c.Writer, "state did not match", http.StatusBadRequest)
			return
		}

		oauth2Token, err := k.oidcProvider.Config.Exchange(context.Background(), c.Request.URL.Query().Get("code"))
		if err != nil {
			http.Error(c.Writer, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			http.Error(c.Writer, "No id_token field in oauth2 token.", http.StatusInternalServerError)
			return
		}

		idToken, err := k.oidcProvider.Verifier().Verify(context.Background(), rawIDToken)
		if err != nil {
			http.Error(c.Writer, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		nonce := session.Get("nonce")
		if nonce == nil {
			http.Error(c.Writer, "nonce not found", http.StatusBadRequest)
			return
		}
		if idToken.Nonce != nonce.(string) {
			http.Error(c.Writer, "nonce did not match", http.StatusBadRequest)
			return
		}

		var claims gophotoClaims
		if err := idToken.Claims(&claims); err != nil {
			http.Error(c.Writer, "nonce did not match", http.StatusBadRequest)
		}

		username := getUsernameFromClaims(claims)

		loggedUser := entityFromClaims(*username, claims)
		loggedUser.Groups = getGroupsFromClaims(claims)

		sessionData := entity.NewSession()

		sessionData.User = loggedUser
		sessionData.TokenID = claims.StandardClaims.Id
		sessionData.SessionID = claims.SessionState
		sessionData.Token = oauth2Token
		sessionData.ExpireAt = time.Unix(claims.StandardClaims.ExpiresAt, 0)
		sessionData.IssueAt = time.Unix(claims.StandardClaims.IssuedAt, 0)

		session.Set(uuid.String(), sessionData)
		session.Save()

		logger.WithField("session data", fmt.Sprintf("%+v", sessionData)).Trace("session data for logged user")

		// save uuid to cookie
		c.SetCookie(sessionID, uuid.String(), 3600, "/", c.Request.Host, true, true)

		next := session.Get("next")
		if next != nil {
			session.Delete("next")
			http.Redirect(c.Writer, c.Request, next.(string), http.StatusFound)
		}
	}
}

func (k *keyCloakAuthenticator) Logout(c *gin.Context, username, refreshToken string) error {
	client := http.DefaultClient

	logoutUrl, err := k.oidcProvider.GetLogoutURL()
	if err != nil {
		return err
	}

	var formData = make(url.Values)
	formData.Add("refresh_token", refreshToken)

	req, err := http.NewRequest(http.MethodPost, logoutUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil
	}

	// set headers
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	basicAuth := fmt.Sprintf("%s:%s", conf.GetKeycloakConfig().ClientID, conf.GetKeycloakConfig().ClientSecret)
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(basicAuth))))

	res, err := client.Do(req.WithContext(c.Request.Context()))
	if err != nil {
		return err
	}

	if res.StatusCode != 204 {
		return fmt.Errorf("%w error logging out user %s", err, username)
	}

	return nil
}

func (k *keyCloakAuthenticator) authenticate(ctx *gin.Context, sessionData entity.Session) error {
	// check if session not expired
	if time.Now().After(sessionData.ExpireAt) {
		// try to refresh the token
		newJwt, err := k.client.RefreshToken(ctx.Request.Context(), sessionData.Token.RefreshToken, k.conf.ClientID, k.conf.ClientSecret, k.conf.Realm)
		if err != nil {
			return errSessionExpired
		}

		sessionData.Token.AccessToken = newJwt.AccessToken
		sessionData.Token.RefreshToken = newJwt.RefreshToken
		sessionData.ExpireAt = time.Now().Add(time.Duration(int64(newJwt.ExpiresIn)) * time.Second)

		logutil.GetLogger(ctx).WithField("username", sessionData.User.Username).WithField("token expire at", sessionData.ExpireAt).Info("session has expired. Token refreshed.")
	}

	ctx.Set("sessionData", sessionData)

	logutil.GetLogger(ctx).WithField("user", fmt.Sprintf("%+v", sessionData)).Info("user logged in")

	return nil
}

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func redirectToLogin(c *gin.Context, config oauth2.Config) {
	session := sessions.Default(c)

	// generate state and nonce
	state, err := randString(16)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Set("state", state)

	// save the current url to redirect back to it if auth is ok.
	session.Set("next", c.Request.URL.String())

	nonce, err := randString(16)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Set("nonce", nonce)
	session.Save()

	http.Redirect(c.Writer, c.Request, config.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)
}

type gophotoClaims struct {
	UserName     *string  `json:"preferred_username"`
	Role         *string  `json:"role"`
	CanShare     *bool    `json:"can_share"`
	Groups       []string `json:"groups"`
	SessionState string   `json:"session_state"`
	FirstName    string   `json:"given_name"`
	LastName     string   `json:"family_name"`
	jwt.StandardClaims
}

func getUsernameFromClaims(claims gophotoClaims) *string {
	if claims.UserName != nil && *claims.UserName != "" {
		return claims.UserName
	}

	return nil
}

func getGroupsFromClaims(claims gophotoClaims) []entity.Group {
	grps := make([]entity.Group, 0, len(claims.Groups))
	for _, name := range claims.Groups {
		name = strings.TrimLeft(name, "/")
		grps = append(grps, entity.Group{Name: name})
	}

	return grps
}

func entityFromClaims(username string, claims gophotoClaims) entity.User {
	var r entity.Role

	if claims.Role != nil {
		switch *claims.Role {
		case "[admin]":
			r = entity.RoleAdmin
		case "[editor]":
			r = entity.RoleEditor
		case "[user]":
			r = entity.RoleUser
		default:
			r = entity.RoleUser
		}
	} else {
		r = entity.RoleUser
	}

	canShare := false
	if claims.CanShare != nil {
		canShare = *claims.CanShare
	}

	return entity.User{
		Username:  username,
		ID:        claims.StandardClaims.Subject,
		Role:      r,
		CanShare:  canShare,
		FirstName: claims.FirstName,
		LastName:  claims.LastName,
	}
}
