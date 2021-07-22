package form

import "html"

type Album struct {
<<<<<<< Updated upstream
	Name             string              `json:"name" binding:"required"`
	Description      string              `json:"description" binding:"required"`
	Location         string              `json:"location"`
	UserPermissions  map[string][]string `json:"user_permissions"`
	GroupPermissions map[string][]string `json:"group_permissions"`
=======
	Name             string `form:"name" binding:"required"`
	Description      string `form:"description" binding:"required"`
	Location         string `form:"location"`
	UserPermissions  string `form:"user_permissions"`
	GroupPermissions string `form:"group_permissions"`
>>>>>>> Stashed changes
}

func (a Album) Sanitize() Album {
	escapedAlbum := Album{
<<<<<<< Updated upstream
		Name:        html.EscapeString(a.Name),
		Description: html.EscapeString(a.Description),
		Location:    html.EscapeString(a.Location),
	}

	if len(a.UserPermissions) > 0 {
		escapedAlbum.UserPermissions = make(map[string][]string)

		for k, v := range a.UserPermissions {
			vals := make([]string, 0, len(v))

			for _, vv := range v {
				vals = append(vals, html.EscapeString(vv))
			}

			escapedAlbum.UserPermissions[html.EscapeString(k)] = vals
		}
	}

	if len(a.GroupPermissions) > 0 {
		escapedAlbum.GroupPermissions = make(map[string][]string)

		for k, v := range a.GroupPermissions {
			vals := make([]string, 0, len(v))

			for _, vv := range v {
				vals = append(vals, html.EscapeString(vv))
			}

			escapedAlbum.GroupPermissions[html.EscapeString(k)] = vals
		}
=======
		Name:             html.EscapeString(a.Name),
		Description:      html.EscapeString(a.Description),
		Location:         html.EscapeString(a.Location),
		UserPermissions:  html.EscapeString(a.UserPermissions),
		GroupPermissions: html.EscapeString(a.GroupPermissions),
>>>>>>> Stashed changes
	}

	return escapedAlbum
}
