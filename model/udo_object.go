package model

import "time"

// UdoObject UDO Object Definition
type UdoObject struct {
	ID          int64     `json:"id" xorm:"id"`                   // Primary key
	TenantId    int64     `json:"tenant_id" xorm:"tenant_id"`     // Tenant ID for multi-tenant isolation
	Code        string    `json:"code" xorm:"code"`               // Unique object code within tenant
	Name        string    `json:"name" xorm:"name"`               // Display name of UDO object
	Description string    `json:"description" xorm:"description"` // Description of the object
	Status      int8      `json:"status" xorm:"status"`           // Status: 1-Active, 0-Inactive
	CreatedAt   time.Time `json:"created_at" xorm:"created_at"`   // Creation time
	UpdatedAt   time.Time `json:"updated_at" xorm:"updated_at"`   // Update time
	CreatedBy   int64     `json:"created_by" xorm:"created_by"`   // Creator user ID
	UpdatedBy   int64     `json:"updated_by" xorm:"updated_by"`   // Updater user ID
}

// TableName 表名称
func (*UdoObject) TableName() string {
	return "udo_object"
}
