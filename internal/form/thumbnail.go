package form

import (
	"html"
)

type AlbumThumbnail struct {
	Image string `form:"image" binding:"required"`
}

func (a AlbumThumbnail) Sanitize() AlbumThumbnail {
	escapeAlbumThumbnail := AlbumThumbnail{
		Image: html.EscapeString(a.Image),
	}

	return escapeAlbumThumbnail
}
