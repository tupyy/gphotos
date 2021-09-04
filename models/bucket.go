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


Table: bucket
[ 0] id                                             INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
[ 1] urn                                            TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] album_id                                       INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []


JSON Sample
-------------------------------------
{    "urn": "pBHWOygeIbHqRvbLApQkonCuJ",    "album_id": 33,    "id": 98}



*/

// Bucket struct is a row record of the bucket table in the gophoto database
type Bucket struct {
	//[ 0] id                                             INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
	ID int32 `gorm:"primary_key;column:id;type:INT4;"`
	//[ 1] urn                                            TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Urn string `gorm:"column:urn;type:TEXT;"`
	//[ 2] album_id                                       INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
	AlbumID int32 `gorm:"column:album_id;type:INT4;"`
}

var bucketTableInfo = &TableInfo{
	Name: "bucket",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "id",
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
			GoFieldName:        "ID",
			GoFieldType:        "int32",
			JSONFieldName:      "id",
			ProtobufFieldName:  "id",
			ProtobufType:       "",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "urn",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "Urn",
			GoFieldType:        "string",
			JSONFieldName:      "urn",
			ProtobufFieldName:  "urn",
			ProtobufType:       "",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "album_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "INT4",
			DatabaseTypePretty: "INT4",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "INT4",
			ColumnLength:       -1,
			GoFieldName:        "AlbumID",
			GoFieldType:        "int32",
			JSONFieldName:      "album_id",
			ProtobufFieldName:  "album_id",
			ProtobufType:       "",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (b *Bucket) TableName() string {
	return "bucket"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (b *Bucket) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (b *Bucket) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (b *Bucket) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (b *Bucket) TableInfo() *TableInfo {
	return bucketTableInfo
}
