package models

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
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


Table: album_permissions
[ 0] owner_id                                       TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] owner_kind                                     USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []
[ 2] album_id                                       TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 3] permissions                                    USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []


JSON Sample
-------------------------------------
{    "owner_id": "mqaTFdnTERQbMoWVmDvbvdtlD",    "owner_kind": "HpNyGkntKuRJqCrecZiTfKUpA",    "album_id": "AywDRatDvIiyYKUBKxrTSbZSQ",    "permissions": "WosePVlQvfMgNdZHqeZkqUEGZ"}



*/

// AlbumPermissions struct is a row record of the album_permissions table in the gophoto database
type AlbumPermissions struct {
	//[ 0] owner_id                                       TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	OwnerID string `gorm:"primary_key;column:owner_id;type:TEXT;"`
	//[ 1] owner_kind                                     USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []
	OwnerKind string `gorm:"column:owner_kind;type:VARCHAR;"`
	//[ 2] album_id                                       TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	AlbumID string `gorm:"primary_key;column:album_id;type:TEXT;"`
	//[ 3] permissions                                    USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []
	Permissions []PermissionID `gorm:"column:permissions;type:VARCHAR;"`
}

var album_permissionsTableInfo = &TableInfo{
	Name: "album_permissions",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "owner_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "OwnerID",
			GoFieldType:        "string",
			JSONFieldName:      "owner_id",
			ProtobufFieldName:  "owner_id",
			ProtobufType:       "",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "owner_kind",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "USER_DEFINED",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "USER_DEFINED",
			ColumnLength:       -1,
			GoFieldName:        "OwnerKind",
			GoFieldType:        "string",
			JSONFieldName:      "owner_kind",
			ProtobufFieldName:  "owner_kind",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "album_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "AlbumID",
			GoFieldType:        "string",
			JSONFieldName:      "album_id",
			ProtobufFieldName:  "album_id",
			ProtobufType:       "",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "permissions",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "USER_DEFINED",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "USER_DEFINED",
			ColumnLength:       -1,
			GoFieldName:        "Permissions",
			GoFieldType:        "string",
			JSONFieldName:      "permissions",
			ProtobufFieldName:  "permissions",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},
	},
}

// TableName sets the insert table name for this struct type
func (a *AlbumPermissions) TableName() string {
	return "album_permissions"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (a *AlbumPermissions) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (a *AlbumPermissions) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (a *AlbumPermissions) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (a *AlbumPermissions) TableInfo() *TableInfo {
	return album_permissionsTableInfo
}
