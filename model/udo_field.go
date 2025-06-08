package model

import (
	"time"

	"xorm.io/xorm"
)

type UdoFieldEnum struct {
	Value string `json:"value"` // Value of the enum option
	Label string `json:"label"` // Label of the enum option
}

// UdoField UDO Field Definition
type UdoField struct {
	ID           int64          `json:"id" xorm:"id"`                       // Primary key
	ObjectId     int64          `json:"object_id" xorm:"object_id"`         // Reference to udo_object.id
	Code         string         `json:"code" xorm:"code"`                   // Field code, unique within object
	Name         string         `json:"name" xorm:"name"`                   // Display name of field
	Description  string         `json:"description" xorm:"description"`     // Field description
	FieldType    string         `json:"field_type" xorm:"field_type"`       // Field data type: string, number, boolean, date, datetime, enum, text, richtext, file, image
	IsRequired   int8           `json:"is_required" xorm:"is_required"`     // Whether field is required: 1-Required, 0-Optional
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
	IsUnique     int8           `json:"is_unique" xorm:"is_unique"`         // Whether field values must be unique: 1-Unique, 0-Not unique
	IsSearchable int8           `json:"is_searchable" xorm:"is_searchable"` // Whether field should be searchable: 1-Searchable, 0-Not searchable
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

// 创建字段
func (m *UdoFieldModel) CreatedField(data UdoField) error {
	// m.db.CreateTables()
	obj := new(UdoObject)
	_, err := m.db.ID(data.ObjectId).Get(obj)
	if err != nil {
		return err
	}
	return nil
}
