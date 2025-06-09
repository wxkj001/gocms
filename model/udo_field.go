package model

import (
	"gocms/utils"
	"strings"
	"time"

	"xorm.io/builder"
	"xorm.io/xorm"
)

type UdoFieldEnum struct {
	Value string `json:"value"` // Value of the enum option
	Label string `json:"label"` // Label of the enum option
}

// UdoField UDO Field Definition
type UdoField struct {
	ID           int64          `json:"id" xorm:"id pk"`                    // Primary key
	ObjectId     int64          `json:"object_id" xorm:"object_id"`         // Reference to udo_object.id
	Code         string         `json:"code" xorm:"code"`                   // Field code, unique within object
	Name         string         `json:"name" xorm:"name"`                   // Display name of field
	Description  string         `json:"description" xorm:"description"`     // Field description
	FieldType    string         `json:"field_type" xorm:"field_type"`       // Field data type: string, number, boolean, date, datetime, enum, text, richtext, file, image
	IsRequired   bool           `json:"is_required" xorm:"is_required"`     // Whether field is required: 1-Required, 0-Optional
	MinLength    int64          `json:"min_length" xorm:"min_length"`       // Minimum length (for string/text types)
	MaxLength    int64          `json:"max_length" xorm:"max_length"`       // Maximum length (for string/text types)
	RegexPattern string         `json:"regex_pattern" xorm:"regex_pattern"` // Regular expression for validation
	RegexMessage string         `json:"regex_message" xorm:"regex_message"` // Message to show when regex validation fails
	MinValue     float64        `json:"min_value" xorm:"min_value"`         // Minimum value (for number type)
	MaxValue     float64        `json:"max_value" xorm:"max_value"`         // Maximum value (for number type)
	EnumOptions  []UdoFieldEnum `json:"enum_options" xorm:"enum_options"`   // JSON array of options for enum type: [{"value": "red", "label": "Red"}, ...]
	DefaultValue string         `json:"default_value" xorm:"default_value"` // Default value in string format
	Placeholder  string         `json:"placeholder" xorm:"placeholder"`     // Input placeholder text
	HelpText     string         `json:"help_text" xorm:"help_text"`         // Help text for this field
	IsUnique     bool           `json:"is_unique" xorm:"is_unique"`         // Whether field values must be unique: 1-Unique, 0-Not unique
	IsSearchable bool           `json:"is_searchable" xorm:"is_searchable"` // Whether field should be searchable: 1-Searchable, 0-Not searchable
	IsSystem     int8           `json:"is_system" xorm:"is_system"`         // Whether field is system field: 1-System, 0-Custom
	SortOrder    int64          `json:"sort_order" xorm:"sort_order"`       // Display order of the field
	Status       int8           `json:"status" xorm:"status"`               // Status: 1-Active, 0-Inactive
	CreatedAt    time.Time      `json:"created_at" xorm:"created"`          // Creation time
	UpdatedAt    time.Time      `json:"updated_at" xorm:"updated"`          // Update time
}

// TableName 表名称
func (*UdoField) TableName() string {
	return "udo_field"
}

type UdoFieldModel struct {
	db *xorm.Engine
	UdoObject
}

func NewUdoField(db *xorm.Engine) *UdoFieldModel {
	return &UdoFieldModel{db: db}
}

// 根据对象ID查询字段
func (m *UdoFieldModel) GetFieldListByObjectId(objectId int64) ([]UdoField, error) {
	var fields []UdoField
	err := m.db.Where("object_id = ?", objectId).Find(&fields)
	return fields, err
}
func (m *UdoFieldModel) GetFieldListPageByObjectId(objectId int64, page, pageSize int, where map[string]string) (int64, []UdoField, error) {
	var fields []UdoField
	wheres := builder.NewCond()
	for k, v := range where {
		wheres.And(builder.Like{k, v})
	}
	count, err := m.db.Where("object_id = ?", objectId).Where(wheres).
		Limit(pageSize, (page-1)*pageSize).
		FindAndCount(&fields)
	return count, fields, err
}

// 创建字段
func (m *UdoFieldModel) CreatedField(data UdoField) error {

	session := m.db.NewSession()
	defer session.Close()
	session.Begin()
	defer session.Rollback()
	obj := new(UdoObject)
	_, err := session.ID(data.ObjectId).Get(obj)
	if err != nil {
		return err
	}
	// 创建字段
	_, err = session.Insert(&data)
	if err != nil {
		return err
	}
	list := []UdoField{}
	err = session.Where("object_id = ?", data.ObjectId).Find(&list)
	if err != nil {
		return err
	}
	fields := []map[string]any{}
	for _, v := range list {
		t := ""
		constraints := map[string]any{
			"comment": v.Name,
		}
		switch v.FieldType {
		case "string":
			t = "string"
			if v.MaxLength == 0 {
				v.MaxLength = 255
			}
			constraints["varchar"] = v.MaxLength
		case "number":
			t = "int"
		case "boolean":
			t = "bool"
		case "date":
			t = "date"
		case "datetime":
			t = "datetime"
		case "enum":
			t = "string"
		case "text":
			t = "string"
		case "richtext":
			t = "string"
		case "file":
			t = "string"
			constraints["varchar"] = 255
		case "image":
			t = "string"
			constraints["varchar"] = 255
		}
		if v.IsUnique {
			constraints["unique"] = ""
		}
		if v.IsSearchable {
			constraints["index"] = ""
		}
		if v.DefaultValue != "" {
			constraints["default"] = v.DefaultValue
		}
		if v.IsRequired {
			constraints["not null"] = ""
		}
		field := map[string]any{
			"name":        v.Code,
			"type":        t,
			"constraints": constraints,
		}
		fields = append(fields, field)
	}
	tablename, table, err := utils.SyncTable("udo_data_"+strings.ToLower(obj.Code), fields)
	if err != nil {
		return err
	}
	session.Table(tablename).Sync2(table)
	session.Commit()
	return nil
}

// 更新
func (m *UdoFieldModel) UpdateField(data UdoField) error {
	session := m.db.NewSession()
	defer session.Close()
	session.Begin()
	defer session.Rollback()
	field := new(UdoField)
	_, err := session.ID(data.ID).Get(field)
	if err != nil {
		return err
	}
	_, err = session.ID(data.ID).Update(&data)
	if err != nil {
		return err
	}
	session.Commit()
	return nil
}

// 删除
func (m *UdoFieldModel) DeleteField(id int64) error {
	session := m.db.NewSession()
	defer session.Close()
	session.Begin()
	defer session.Rollback()
	field := new(UdoField)
	_, err := session.ID(id).Get(field)
	if err != nil {
		return err
	}
	obj := new(UdoObject)
	_, err = session.ID(field.ObjectId).Get(obj)
	if err != nil {
		return err
	}
	_, err = session.ID(id).Delete(new(UdoField))
	if err != nil {
		return err
	}
	list := []UdoField{}
	err = session.Where("object_id = ?", obj.ID).Find(&list)
	if err != nil {
		return err
	}
	fields := []map[string]any{}
	for _, v := range list {
		t := ""
		constraints := map[string]any{
			"comment": v.Name,
		}
		switch v.FieldType {
		case "string":
			t = "string"
			if v.MaxLength == 0 {
				v.MaxLength = 255
			}
			constraints["varchar"] = v.MaxLength
		case "number":
			t = "int"
		case "boolean":
			t = "bool"
		case "date":
			t = "date"
		case "datetime":
			t = "datetime"
		case "enum":
			t = "string"
		case "text":
			t = "string"
		case "richtext":
			t = "string"
		case "file":
			t = "string"
			constraints["varchar"] = 255
		case "image":
			t = "string"
			constraints["varchar"] = 255
		}
		if v.IsUnique {
			constraints["unique"] = ""
		}
		if v.IsSearchable {
			constraints["index"] = ""
		}
		if v.DefaultValue != "" {
			constraints["default"] = v.DefaultValue
		}
		if v.IsRequired {
			constraints["not null"] = ""
		}
		field := map[string]any{
			"name":        v.Code,
			"type":        t,
			"constraints": constraints,
		}
		fields = append(fields, field)
	}
	tablename, table, err := utils.SyncTable("udo_data_"+strings.ToLower(obj.Code), fields)
	if err != nil {
		return err
	}
	session.Table(tablename).SyncWithOptions(xorm.SyncOptions{
		WarnIfDatabaseColumnMissed: false,
		IgnoreConstrains:           false,
		IgnoreIndices:              false,
		IgnoreDropIndices:          false,
	}, table)
	session.Commit()
	return nil
}
