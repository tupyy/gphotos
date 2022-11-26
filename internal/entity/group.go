package entity

type Group struct {
	Name  string `json:"name"`
	Users []User `json:"-"`
}
