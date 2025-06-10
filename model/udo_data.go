package model

import (
	"time"

	"xorm.io/builder"
	"xorm.io/xorm"
)

// UdoData UDO Data Storage
type UdoData struct {
	ID        int64     `json:"id" xorm:"id pk autoincr notnull unique index"` // Primary key
	TenantId  int64     `json:"tenant_id" xorm:"tenant_id"`                    // Tenant ID for multi-tenant isolation
	ObjectId  int64     `json:"object_id" xorm:"object_id"`                    // Reference to udo_object.id
	Data      string    `json:"data" xorm:"data"`                              // Data in JSON format
	Status    int8      `json:"status" xorm:"status"`                          // Status: 1-Active, 0-Inactive
	CreatedAt time.Time `json:"created_at" xorm:"created_at"`                  // Creation time
	UpdatedAt time.Time `json:"updated_at" xorm:"updated_at"`                  // Update time
	CreatedBy int64     `json:"created_by" xorm:"created_by"`                  // Creator user ID
	UpdatedBy int64     `json:"updated_by" xorm:"updated_by"`                  // Updater user ID
}

// TableName 表名称
func (*UdoData) TableName() string {
	return "udo_data"
}

type UdoDataModel struct {
	db *xorm.Engine
}

// NewUdoData 创建 UdoData 实例
func NewUdoData(db *xorm.Engine) *UdoDataModel {
	return &UdoDataModel{db: db}
}

// 通过code获取表数据
func (m *UdoDataModel) GetListByCode(code string, page, pageSize int, where map[string]string) (int64, []map[string]any, error) {
	udoData := make([]map[string]any, 0)
	wheres := builder.NewCond()
	for k, v := range where {
		wheres.And(builder.Like{k, v})
	}
	count, err := m.db.Table("udo_data_"+code).
		Where(wheres).
		Limit(pageSize, (page-1)*pageSize).
		FindAndCount(udoData)
	if err != nil {
		return count, nil, err
	}

	return count, udoData, nil
}
