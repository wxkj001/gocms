package model

import (
	"strings"
	"time"

	"xorm.io/builder"
	"xorm.io/xorm"
)

// UdoObject UDO Object Definition
type UdoObject struct {
	ID          int64     `json:"id" xorm:"id pk autoincr notnull unique index"` // Primary key
	Code        string    `json:"code" xorm:"code"`                              // Unique object code within tenant
	Name        string    `json:"name" xorm:"name"`                              // Display name of UDO object
	Description string    `json:"description" xorm:"description"`                // Description of the object
	Status      int8      `json:"status" xorm:"status"`                          // Status: 1-Active, 0-Inactive
	CreatedAt   time.Time `json:"created_at" xorm:"created"`                     // Creation time
	UpdatedAt   time.Time `json:"updated_at" xorm:"updated"`                     // Update time
}

// TableName 表名称
func (*UdoObject) TableName() string {
	return "udo_object"
}

type UdoObjectModel struct {
	db *xorm.Engine
	UdoObject
}

func NewUdoObject(db *xorm.Engine) *UdoObjectModel {
	return &UdoObjectModel{db: db}
}

// 根据code查询对象
func (m *UdoObjectModel) GetUdoObjectByCode(code string) (*UdoObject, error) {
	udoObject := &UdoObject{}
	has, err := m.db.Where("code = ?", code).Get(udoObject)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return udoObject, nil
}

// 根据ID查询对象
func (m *UdoObjectModel) GetUdoObjectByID(id int64) (*UdoObject, error) {
	udoObject := &UdoObject{}
	has, err := m.db.ID(id).Get(udoObject)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return udoObject, nil
}

// 创建对象
func (m *UdoObjectModel) CreatedObject(data *UdoObject) error {
	_, err := m.db.Insert(data)
	if err != nil {
		return err
	}
	return nil
}

// 更新对象
func (m *UdoObjectModel) UpdateObject(data *UdoObject) error {
	_, err := m.db.ID(data.ID).Update(data)
	if err != nil {
		return err
	}
	return nil
}

// 获取对象列表带分页
func (m *UdoObjectModel) ObjectPage(page, pageSize int, where map[string]string) (int64, []UdoObject, error) {
	var list []UdoObject
	wheres := builder.NewCond()
	for k, v := range where {
		wheres.And(builder.Like{k, v})
	}
	count, err := m.db.Where(wheres).Limit(pageSize, (page-1)*pageSize).FindAndCount(&list)

	return count, list, err
}

// 获取对象列表不带分页
func (m *UdoObjectModel) ObjectList(where map[string]string) ([]UdoObject, error) {
	var list []UdoObject
	wheres := builder.NewCond()
	for k, v := range where {
		wheres.And(builder.Like{k, v})
	}
	err := m.db.Where(wheres).Find(&list)

	return list, err
}

// 删除
func (m *UdoObjectModel) DeleteObject(id int64) error {
	obj, err := m.GetUdoObjectByID(id)
	if err != nil {
		return err
	}
	_, err = m.db.ID(id).Delete(new(UdoObject))
	if err != nil {
		return err
	}
	m.db.DropTables("udo_data_" + strings.ToLower(obj.Code))
	return nil
}
