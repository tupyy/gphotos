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


Table: groups
[ 0] id                                             INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
[ 1] name                                           VARCHAR(100)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 100     default: []


JSON Sample
-------------------------------------
{    "id": 19,    "name": "dzLWxMLRwtCPCytjYBVkxwRbr"}



*/

// Groups struct is a row record of the groups table in the gophoto database
type Groups struct {
	//[ 0] id                                             INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
	ID int32 `gorm:"primary_key;column:id;type:INT4;"`
	//[ 1] name                                           VARCHAR(100)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 100     default: []
	Name string `gorm:"column:name;type:VARCHAR;size:100;"`
}

var groupsTableInfo = &TableInfo{
	Name: "groups",
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
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(100)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       100,
			GoFieldName:        "Name",
			GoFieldType:        "string",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (g *Groups) TableName() string {
	return "groups"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (g *Groups) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (g *Groups) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (g *Groups) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (g *Groups) TableInfo() *TableInfo {
	return groupsTableInfo
}
