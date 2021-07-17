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


Table: users
[ 0] id                                             INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
[ 1] username                                       VARCHAR(100)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 100     default: []
[ 2] role                                           USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []
[ 3] user_id                                        TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 4] can_share                                      BOOL                 null: true   primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: []


JSON Sample
-------------------------------------
{    "can_share": true,    "id": 3,    "username": "LydEupdsfAvbJyWvhmmmzcIjo",    "role": 52,    "user_id": "oZSwCSlmiaGWLoavkcaqkcZNh"}



*/

// Users struct is a row record of the users table in the gophoto database
type Users struct {
	//[ 0] id                                             INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
	ID int32 `gorm:"primary_key;column:id;type:INT4;"`
	//[ 1] username                                       VARCHAR(100)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 100     default: []
	Username string `gorm:"column:username;type:VARCHAR;size:100;"`
	//[ 2] role                                           USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []
	Role Role `gorm:"column:role;type:ROLE;"`
	//[ 3] user_id                                        TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	UserID string `gorm:"column:user_id;type:TEXT;"`
	//[ 4] can_share                                      BOOL                 null: true   primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: []
	CanShare *bool `gorm:"column:can_share;type:BOOL;"`
}

var usersTableInfo = &TableInfo{
	Name: "users",
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
			Name:               "username",
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
			GoFieldName:        "Username",
			GoFieldType:        "string",
			JSONFieldName:      "username",
			ProtobufFieldName:  "username",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "role",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "ROLE",
			DatabaseTypePretty: "USER_DEFINED",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "USER_DEFINED",
			ColumnLength:       -1,
			GoFieldName:        "Role",
			GoFieldType:        "Role",
			JSONFieldName:      "role",
			ProtobufFieldName:  "role",
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

		&ColumnInfo{
			Index:              4,
			Name:               "can_share",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "BOOL",
			DatabaseTypePretty: "BOOL",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "BOOL",
			ColumnLength:       -1,
			GoFieldName:        "CanShare",
			GoFieldType:        "*bool",
			JSONFieldName:      "can_share",
			ProtobufFieldName:  "can_share",
			ProtobufType:       "",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (u *Users) TableName() string {
	return "users"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (u *Users) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (u *Users) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (u *Users) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (u *Users) TableInfo() *TableInfo {
	return usersTableInfo
}
