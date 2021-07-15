package entity

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type Session struct {
	User      User
	TokenID   string
	SessionID string
	Token     *oauth2.Token
	ExpireAt  time.Time
	IssueAt   time.Time
}

func (s Session) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "username: %s ", s.User.Username)
	fmt.Fprintf(&sb, "user_id: %s ", s.User.UserID)
	fmt.Fprintf(&sb, "role: %s ", s.User.Role.String())
	fmt.Fprintf(&sb, "groups: %+v ", s.User.Groups)
	fmt.Fprintf(&sb, "AccessToken; %s ", s.Token.AccessToken)
	fmt.Fprintf(&sb, "RefreshToken; %s ", s.Token.RefreshToken)
	fmt.Fprintf(&sb, "ExpireAt: %s ", s.ExpireAt.Format("Mon Jan 2 15:04:05 MST 2006"))
	fmt.Fprintf(&sb, "IssueAt: %s ", s.ExpireAt.Format("Mon Jan 2 15:04:05 MST 2006"))

	return sb.String()
}
