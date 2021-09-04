package entity

type Bucket struct {
	// ID - id of the bucket
	ID int32
	// Urn - name of the bucket on store
	Urn string
	// AlbumID - id of the album associated withe the bucket.
	AlbumID int32
}
