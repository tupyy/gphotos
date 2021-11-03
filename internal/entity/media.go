package entity

type MediaType int

const (
	Photo MediaType = iota
	Video
	Unknown
)

type Media struct {
	MediaType MediaType
	Filename  string
	Bucket    string
	Thumbnail string
}
