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

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repo"
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
	userRepo     UserRepo
	groupRepo    GroupRepo
	oidcProvider *OidcProvider
}

func NewKeyCloakAuthenticator(oidcProvider *OidcProvider, ur UserRepo, gr GroupRepo) Authenticator {
	return &keyCloakAuthenticator{oidcProvider: oidcProvider, userRepo: ur, groupRepo: gr}
}

type gophotoClaims struct {
	UserName     *string  `json:"preferred_username"`
	Role         *string  `json:"role"`
	CanShare     *bool    `json:"can_share"`
	Groups       []string `json:"groups"`
	SessionState string   `json:"session_state"`
	jwt.StandardClaims
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

		user, err := k.createOrUpdateUserFromClaims(c, normalizeGroupsFromClaims(claims))
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		}

		sessionData := entity.Session{
			User:      user,
			TokenID:   claims.StandardClaims.Id,
			SessionID: claims.SessionState,
			Token:     oauth2Token,
			ExpireAt:  time.Unix(claims.StandardClaims.ExpiresAt, 0),
			IssueAt:   time.Unix(claims.StandardClaims.IssuedAt, 0),
		}

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
		return errSessionExpired
	}

	ctx.Set("sessionData", sessionData)

	logutil.GetLogger(ctx).Debug("user logged in")

	return nil
}

// createOrUpdateUserFromClaims creates or updates an existing user from claims.
func (k *keyCloakAuthenticator) createOrUpdateUserFromClaims(ctx *gin.Context, claims gophotoClaims) (entity.User, error) {
	var noUser entity.User

	logger := logutil.GetLogger(ctx)

	username := getUsernameFromClaims(claims)

	loggedUser := entityFromClaims(*username, claims)
	// create or update user in db
	user, err := k.userRepo.GetByUsername(ctx, *username)
	if err != nil {
		if err != repo.ErrUserNotFound {
			logger.WithError(err).Error("failed to get user")
			return noUser, errInternalError
		}

		if groups, err := k.getGroupsFromClaims(ctx, claims.Groups); err != nil {
			logger.WithError(err).Error("cannot retrieve groups from claims")
		} else {
			loggedUser.Groups = groups
		}

		if id, err := k.userRepo.Create(ctx.Request.Context(), loggedUser); err != nil {
			logger.WithError(err).Error("failed to create user")
			return noUser, errInternalError
		} else {
			loggedUser.ID = &id

			logger.WithField("user id", id).WithField("username", *username).Debug("user created")
		}
	} else {
		// update user id
		loggedUser.ID = user.ID

		if groups, err := k.getGroupsFromClaims(ctx, claims.Groups); err != nil {
			logger.WithError(err).Error("cannot retrieve groups from claims")
		} else {
			loggedUser.Groups = groups
		}

		err := k.userRepo.Update(ctx.Request.Context(), loggedUser)
		if err != nil {
			logger.WithError(err).Error("failed to update user")
			return noUser, errInternalError
		}

		logger.WithField("user", fmt.Sprintf("%+v", loggedUser)).Debug("user updated")
	}

	return loggedUser, nil
}

func (k *keyCloakAuthenticator) getGroupsFromClaims(ctx context.Context, groups []string) ([]entity.Group, error) {
	grps := make([]entity.Group, 0, len(groups))
	for _, name := range groups {
		group, err := k.groupRepo.GetByName(ctx, name)
		if err != nil {
			if errors.Is(err, repo.ErrGroupNotFound) {
				id, err := k.groupRepo.Create(ctx, entity.Group{Name: name})
				if err != nil {
					return []entity.Group{}, err
				} else {
					logrus.WithField("group", name).WithField("id", id).Debug("group created")
					group = entity.Group{ID: &id, Name: name}
				}
			}
		}

		grps = append(grps, group)
	}

	return grps, nil
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

func getUsernameFromClaims(claims gophotoClaims) *string {
	if claims.UserName != nil && *claims.UserName != "" {
		return claims.UserName
	}

	return nil
}

func normalizeGroupsFromClaims(claims gophotoClaims) gophotoClaims {

	groups := make([]string, 0, len(claims.Groups))
	for _, g := range claims.Groups {
		groups = append(groups, strings.TrimLeft(g, "/"))
	}

	claims.Groups = groups

	return claims
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
		Username: username,
		UserID:   claims.StandardClaims.Subject,
		Role:     r,
		CanShare: canShare,
	}
}
