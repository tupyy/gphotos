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


Table: tag
[ 0] id                                             INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
[ 1] name                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] color                                          TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 3] user_id                                        TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "color": "aZbHVaFScXgVCrKfTJNEPhmlL",    "user_id": "dRqpAyEGkQmVNpSuyWsDCudic",    "id": 5,    "name": "cBnJZaiaJaiDKBOTKBvVOenHD"}



*/

// Tag struct is a row record of the tag table in the gophoto database
type Tag struct {
	//[ 0] id                                             INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;"`
	//[ 1] name                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Name string `gorm:"column:name;type:TEXT;"`
	//[ 2] color                                          TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Color *string `gorm:"column:color;type:TEXT;"`
	//[ 3] user_id                                        TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	UserID string `gorm:"column:user_id;type:TEXT;"`
}

var tagTableInfo = &TableInfo{
	Name: "tag",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "id",
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
			GoFieldName:        "ID",
			GoFieldType:        "string",
			JSONFieldName:      "id",
			ProtobufFieldName:  "id",
			ProtobufType:       "",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "name",
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
			GoFieldName:        "Name",
			GoFieldType:        "string",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "color",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "Color",
			GoFieldType:        "*string",
			JSONFieldName:      "color",
			ProtobufFieldName:  "color",
			ProtobufType:       "",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "user_id",
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
			GoFieldName:        "UserID",
			GoFieldType:        "string",
			JSONFieldName:      "user_id",
			ProtobufFieldName:  "user_id",
			ProtobufType:       "",
			ProtobufPos:        4,
		},
	},
}

// TableName sets the insert table name for this struct type
func (t *Tag) TableName() string {
	return "tag"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (t *Tag) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (t *Tag) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (t *Tag) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (t *Tag) TableInfo() *TableInfo {
	return tagTableInfo
}
