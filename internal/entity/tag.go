package entity

import "fmt"

type Tag struct {
	// Name - name of the tag
	Name string
	// Color - color of the tag (optional)
	Color *string
}

func (t Tag) String() string {
	if t.Color == nil {
		return fmt.Sprintf("Name: %s", t.Name)
	}

	return fmt.Sprintf("Name: %s Color: %s", t.Name, *t.Color)
}
