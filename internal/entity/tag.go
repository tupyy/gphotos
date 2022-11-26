package entity

import "fmt"

type Tag struct {
	ID string
	// UserID - id of the owner
	UserID string
	// Name - name of the tag
	Name string
	// Color - color of the tag (optional)
	Color *string
	// Albums -- list of associated albums
	Albums []string
}

func (t Tag) String() string {
	if t.Color == nil {
		return fmt.Sprintf("Name: %s", t.Name)
	}

	return fmt.Sprintf("UserID: %s Name: %s Color: %s", t.UserID, t.Name, *t.Color)
}
