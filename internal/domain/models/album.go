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


Table: album
[ 0] id                                             INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
[ 1] name                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] created_at                                     TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [timezone('UTC']
[ 3] owner_id                                       TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 4] bucket                                         TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 5] description                                    TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 6] location                                       TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 7] thumbnail                                      VARCHAR(100)         null: true   primary: false  isArray: false  auto: false  col: VARCHAR         len: 100     default: []


JSON Sample
-------------------------------------
{    "description": "VVwFjkeGBbKvyCgzbYKOiCSdO",    "location": "HqipdySOyjJGhAXddMmojynpc",    "thumbnail": "uBgEDbmuYoPRHfmbYWNRpdvxL",    "id": 63,    "name": "FWhycuHcZaPgFbydTUWAimPUV",    "created_at": "2107-04-06T08:36:30.307147087+02:00",    "owner_id": "mIbIwHSPCgbEqbKoQNunBfLIu",    "bucket": "fbJmQrpTDcWMxwaQYZxvQCJAd"}



*/

// Album struct is a row record of the album table in the gophoto database
type Album struct {
	//[ 0] id                                             INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
	ID int32 `gorm:"primary_key;column:id;type:INT4;"`
	//[ 1] name                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Name string `gorm:"column:name;type:TEXT;"`
	//[ 2] created_at                                     TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [timezone('UTC']
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:timezone('UTC';"`
	//[ 3] owner_id                                       TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	OwnerID string `gorm:"column:owner_id;type:TEXT;"`
	//[ 4] bucket                                         TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Bucket string `gorm:"column:bucket;type:TEXT;"`
	//[ 5] description                                    TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Description *string `gorm:"column:description;type:TEXT;"`
	//[ 6] location                                       TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Location *string `gorm:"column:location;type:TEXT;"`
	//[ 7] thumbnail                                      VARCHAR(100)         null: true   primary: false  isArray: false  auto: false  col: VARCHAR         len: 100     default: []
	Thumbnail sql.NullString `gorm:"column:thumbnail;type:VARCHAR;size:100;"`
}

var albumTableInfo = &TableInfo{
	Name: "album",
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
			Name:               "created_at",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TIMESTAMP",
			DatabaseTypePretty: "TIMESTAMP",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TIMESTAMP",
			ColumnLength:       -1,
			GoFieldName:        "CreatedAt",
			GoFieldType:        "time.Time",
			JSONFieldName:      "created_at",
			ProtobufFieldName:  "created_at",
			ProtobufType:       "",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "owner_id",
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
			GoFieldName:        "OwnerID",
			GoFieldType:        "string",
			JSONFieldName:      "owner_id",
			ProtobufFieldName:  "owner_id",
			ProtobufType:       "",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "bucket",
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
			GoFieldName:        "Bucket",
			GoFieldType:        "string",
			JSONFieldName:      "bucket",
			ProtobufFieldName:  "bucket",
			ProtobufType:       "",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "description",
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
			GoFieldName:        "Description",
			GoFieldType:        "*string",
			JSONFieldName:      "description",
			ProtobufFieldName:  "description",
			ProtobufType:       "",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "location",
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
			GoFieldName:        "Location",
			GoFieldType:        "*string",
			JSONFieldName:      "location",
			ProtobufFieldName:  "location",
			ProtobufType:       "",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "thumbnail",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(100)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       100,
			GoFieldName:        "Thumbnail",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "thumbnail",
			ProtobufFieldName:  "thumbnail",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},
	},
}

// TableName sets the insert table name for this struct type
func (a *Album) TableName() string {
	return "album"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (a *Album) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (a *Album) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (a *Album) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (a *Album) TableInfo() *TableInfo {
	return albumTableInfo
}
