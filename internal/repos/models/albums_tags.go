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


Table: albums_tags
[ 0] album_id                                       TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] tag_id                                         TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "album_id": 1,    "tag_id": 81}



*/

// AlbumsTags struct is a row record of the albums_tags table in the gophoto database
type AlbumsTags struct {
	//[ 0] album_id                                       TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	AlbumID string `gorm:"primary_key;column:album_id;type:TEXT;"`
	//[ 1] tag_id                                         TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	TagID string `gorm:"primary_key;column:tag_id;type:TEXT;"`
}

var albums_tagsTableInfo = &TableInfo{
	Name: "albums_tags",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
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
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "tag_id",
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
			GoFieldName:        "TagID",
			GoFieldType:        "string",
			JSONFieldName:      "tag_id",
			ProtobufFieldName:  "tag_id",
			ProtobufType:       "",
			ProtobufPos:        2,
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
