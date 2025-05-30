package model

import "time"

// UdoLog UDO Change Log
type UdoLog struct {
	ID        int64     `json:"id" xorm:"id"`                 // Primary key
	TenantId  int64     `json:"tenant_id" xorm:"tenant_id"`   // Tenant ID for multi-tenant isolation
	DataId    int64     `json:"data_id" xorm:"data_id"`       // Reference to udo_data.id
	ObjectId  int64     `json:"object_id" xorm:"object_id"`   // Reference to udo_object.id
	Action    string    `json:"action" xorm:"action"`         // Action: create, update, delete
	CreatedAt time.Time `json:"created_at" xorm:"created_at"` // Log creation time
	CreatedBy int64     `json:"created_by" xorm:"created_by"` // User who performed the action
}

// TableName 表名称
func (*UdoLog) TableName() string {
	return "udo_log"
}
