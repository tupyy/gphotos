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


Table: token_blacklist
[ 0] id                                             INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
[ 1] token                                          TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": 80,    "token": "YykxEYZZZSLWdkcLPZTtjVIJr"}



*/

// TokenBlacklist struct is a row record of the token_blacklist table in the gophoto database
type TokenBlacklist struct {
	//[ 0] id                                             INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
	ID int32 `gorm:"primary_key;column:id;type:INT4;"`
	//[ 1] token                                          TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Token string `gorm:"column:token;type:TEXT;"`
}

var token_blacklistTableInfo = &TableInfo{
	Name: "token_blacklist",
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
			Name:               "token",
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
			GoFieldName:        "Token",
			GoFieldType:        "string",
			JSONFieldName:      "token",
			ProtobufFieldName:  "token",
			ProtobufType:       "",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (t *TokenBlacklist) TableName() string {
	return "token_blacklist"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (t *TokenBlacklist) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (t *TokenBlacklist) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (t *TokenBlacklist) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (t *TokenBlacklist) TableInfo() *TableInfo {
	return token_blacklistTableInfo
}
