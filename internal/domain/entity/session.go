package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type Alert struct {
	Message string
	IsError bool
}

type Session struct {
	User       User
	TokenID    string
	SessionID  string
	Token      *oauth2.Token
	ExpireAt   time.Time
	IssueAt    time.Time
	Attributes map[string]interface{}
	Alerts     map[string]Alert
}

func NewSession() *Session {
	return &Session{
		Attributes: make(map[string]interface{}),
		Alerts:     make(map[string]Alert),
	}
}

func (s *Session) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "username: %s ", s.User.Username)
	fmt.Fprintf(&sb, "user_id: %s ", s.User.ID)
	fmt.Fprintf(&sb, "role: %s ", s.User.Role.String())
	fmt.Fprintf(&sb, "groups: %+v ", s.User.Groups)
	fmt.Fprintf(&sb, "ExpireAt: %s ", s.ExpireAt.Format("Mon Jan 2 15:04:05 MST 2006"))
	fmt.Fprintf(&sb, "IssueAt: %s ", s.ExpireAt.Format("Mon Jan 2 15:04:05 MST 2006"))

	return sb.String()
}

func (s *Session) AddAlert(a Alert) {
	id := uuid.New()
	s.Alerts[id.String()] = a
}

func (s *Session) ClearAlerts() {
	for k := range s.Alerts {
		delete(s.Alerts, k)
	}
}
