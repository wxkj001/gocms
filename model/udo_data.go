package model

import "time"

// UdoData UDO Data Storage
type UdoData struct {
	ID        int64     `json:"id" xorm:"id"`                 // Primary key
	TenantId  int64     `json:"tenant_id" xorm:"tenant_id"`   // Tenant ID for multi-tenant isolation
	ObjectId  int64     `json:"object_id" xorm:"object_id"`   // Reference to udo_object.id
	Status    int8      `json:"status" xorm:"status"`         // Status: 1-Active, 0-Inactive
	CreatedAt time.Time `json:"created_at" xorm:"created_at"` // Creation time
	UpdatedAt time.Time `json:"updated_at" xorm:"updated_at"` // Update time
	CreatedBy int64     `json:"created_by" xorm:"created_by"` // Creator user ID
	UpdatedBy int64     `json:"updated_by" xorm:"updated_by"` // Updater user ID
}

// TableName 表名称
func (*UdoData) TableName() string {
	return "udo_data"
}
