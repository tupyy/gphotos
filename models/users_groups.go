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


Table: users_groups
[ 0] users_id                                       INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
[ 1] groups_id                                      INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []


JSON Sample
-------------------------------------
{    "users_id": 6,    "groups_id": 55}



*/

// UsersGroups struct is a row record of the users_groups table in the gophoto database
type UsersGroups struct {
	//[ 0] users_id                                       INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
	UsersID int32 `gorm:"primary_key;column:users_id;type:INT4;"`
	//[ 1] groups_id                                      INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
	GroupsID int32 `gorm:"primary_key;column:groups_id;type:INT4;"`
}

var users_groupsTableInfo = &TableInfo{
	Name: "users_groups",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "users_id",
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
			GoFieldName:        "UsersID",
			GoFieldType:        "int32",
			JSONFieldName:      "users_id",
			ProtobufFieldName:  "users_id",
			ProtobufType:       "",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "groups_id",
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
			GoFieldName:        "GroupsID",
			GoFieldType:        "int32",
			JSONFieldName:      "groups_id",
			ProtobufFieldName:  "groups_id",
			ProtobufType:       "",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (u *UsersGroups) TableName() string {
	return "users_groups"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (u *UsersGroups) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (u *UsersGroups) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (u *UsersGroups) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (u *UsersGroups) TableInfo() *TableInfo {
	return users_groupsTableInfo
}
