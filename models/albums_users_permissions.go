package models

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	"github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

/*
DB Table Details
-------------------------------------


Table: albums_users_permissions
[ 0] user_id                                        INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
[ 1] album_id                                       INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
[ 2] permissions                                    USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []


JSON Sample
-------------------------------------
{    "user_id": 68,    "album_id": 4,    "permissions": 75}



*/

// AlbumsUsersPermissions struct is a row record of the albums_users_permissions table in the gophoto database
type AlbumsUsersPermissions struct {
	//[ 0] user_id                                        INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
	UserID int32 `gorm:"primary_key;column:user_id;type:INT4;"`
	//[ 1] album_id                                       INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
	AlbumID int32 `gorm:"primary_key;column:album_id;type:INT4;"`
	//[ 2] permissions                                    USER_DEFINED         null: false  primary: false  isArray: true   auto: false  col: USER_DEFINED    len: -1      default: []
	Permissions PermissionIDs `gorm:"column:permissions;type:_PERMISSION_ID;"`
}

var albums_users_permissionsTableInfo = &TableInfo{
	Name: "albums_users_permissions",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "user_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "INT4",
			DatabaseTypePretty: "INT4",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "INT4",
			ColumnLength:       -1,
			GoFieldName:        "UserID",
			GoFieldType:        "int32",
			JSONFieldName:      "user_id",
			ProtobufFieldName:  "user_id",
			ProtobufType:       "",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "album_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "INT4",
			DatabaseTypePretty: "INT4",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "INT4",
			ColumnLength:       -1,
			GoFieldName:        "AlbumID",
			GoFieldType:        "int32",
			JSONFieldName:      "album_id",
			ProtobufFieldName:  "album_id",
			ProtobufType:       "",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "permissions",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "_PERMISSION_ID",
			DatabaseTypePretty: "USER_DEFINED",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            true,
			ColumnType:         "USER_DEFINED",
			ColumnLength:       -1,
			GoFieldName:        "Permissions",
			GoFieldType:        "PermissionIDs",
			JSONFieldName:      "permissions",
			ProtobufFieldName:  "permissions",
			ProtobufType:       "",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (a *AlbumsUsersPermissions) TableName() string {
	return "albums_users_permissions"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (a *AlbumsUsersPermissions) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (a *AlbumsUsersPermissions) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (a *AlbumsUsersPermissions) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (a *AlbumsUsersPermissions) TableInfo() *TableInfo {
	return albums_users_permissionsTableInfo
}
