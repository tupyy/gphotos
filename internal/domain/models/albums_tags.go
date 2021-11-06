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


Table: albums_tags
[ 0] user_id                                        TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] album_id                                       INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
[ 2] tag_id                                         INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []


JSON Sample
-------------------------------------
{    "album_id": 87,    "tag_id": 39,    "user_id": "dKFzLsYueBjcpsiRwxLsntzNu"}



*/

// AlbumsTags struct is a row record of the albums_tags table in the gophoto database
type AlbumsTags struct {
	//[ 0] user_id                                        TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	UserID string `gorm:"primary_key;column:user_id;type:TEXT;"`
	//[ 1] album_id                                       INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
	AlbumID int32 `gorm:"primary_key;column:album_id;type:INT4;"`
	//[ 2] tag_id                                         INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
	TagID int32 `gorm:"primary_key;column:tag_id;type:INT4;"`
}

var albums_tagsTableInfo = &TableInfo{
	Name: "albums_tags",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "user_id",
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
			GoFieldName:        "UserID",
			GoFieldType:        "string",
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
			Name:               "tag_id",
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
			GoFieldName:        "TagID",
			GoFieldType:        "int32",
			JSONFieldName:      "tag_id",
			ProtobufFieldName:  "tag_id",
			ProtobufType:       "",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (a *AlbumsTags) TableName() string {
	return "albums_tags"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (a *AlbumsTags) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (a *AlbumsTags) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (a *AlbumsTags) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (a *AlbumsTags) TableInfo() *TableInfo {
	return albums_tagsTableInfo
}
