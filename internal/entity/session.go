package entity

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type Alert struct {
	Message string
	IsError bool
}

type Session struct {
	User       User   `json:"user"`
	TokenID    string `json:"-"`
	SessionID  string `json:"session_id"`
	Token      *oauth2.Token
	ExpireAt   time.Time
	IssueAt    time.Time
	Attributes map[string]interface{} `json:"-"`
}

func NewSession() *Session {
	return &Session{
		Attributes: make(map[string]interface{}),
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
