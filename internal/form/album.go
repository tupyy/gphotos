package form

import "html"

type Album struct {
	Name             string `form:"name" binding:"required"`
	Description      string `form:"description" binding:"required"`
	Location         string `form:"location"`
	UserPermissions  string `form:"user_permissions"`
	GroupPermissions string `form:"group_permissions"`
}

func (a Album) Sanitize() Album {
	escapedAlbum := Album{
		Name:             html.EscapeString(a.Name),
		Description:      html.EscapeString(a.Description),
		Location:         html.EscapeString(a.Location),
		UserPermissions:  html.EscapeString(a.UserPermissions),
		GroupPermissions: html.EscapeString(a.GroupPermissions),
	}

	return escapedAlbum
}
