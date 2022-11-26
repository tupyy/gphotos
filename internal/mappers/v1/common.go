package v1

type EncryptionService interface {
	// Encrypt data in a deterministic way.
	Encrypt(data string) (string, error)
	// Decrypt data.
	Decrypt(data string) (string, error)
}

const (
	AlbumKind            string = "Album"
	AlbumListKind        string = "AlbumList"
	AlbumPermissionsKind string = "AlbumPermissionsList"
	PhotoKind            string = "Photo"
	PhotoListKind        string = "PhotoList"
	UserKind             string = "User"
	GroupKind            string = "Group"
	TagKind              string = "Tag"
	TagListKind          string = "TagList"
)
