package form

import (
	"html"
)

type Tag struct {
	Name  string `form:"name" binding:"required"`
	Color string `form:"color" binding:"required"`
}

func (t Tag) Sanitize() Tag {
	escapedTag := Tag{
		Name:  html.EscapeString(t.Name),
		Color: html.EscapeString(t.Color),
	}

	return escapedTag
}
