package model

import (
	"xorm.io/xorm"
)

type Role struct {
	ID     int64  `json:"id" xorm:"id autoincr pk"`
	Code   string `json:"code" xorm:"code"`
	Name   string `json:"name" xorm:"name"`
	Enable bool   `json:"enable" xorm:"enable default 1 notnull tinyint(1)"`
}

// TableName 表名称
func (*Role) TableName() string {
	return "role"
}

type RoleModel struct {
	db *xorm.Engine
	Role
}

func NewRole(db *xorm.Engine) *RoleModel {
	return &RoleModel{db: db}
}

type RolePerGroup struct {
	Role          `xorm:"extends"`
	PermissionIds []int64 `json:"permissionIds" xorm:"-"`
}

func (*RolePerGroup) TableName() string {
	return "role"
}

// 获取列表
func (m *RoleModel) GetList(pageSize, pageNo int, name string) ([]RolePerGroup, int64, error) {
	list := make([]RolePerGroup, 0)
	db := m.db.NewSession()
	db1 := m.db.NewSession()
	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
		db1 = db1.Where("name like ?", "%"+name+"%")
	}

	count, err := db.Table("role").Count()
	if err != nil {
		return nil, 0, err
	}
	err = db1.Limit(pageSize, (pageNo-1)*pageSize).Find(&list)
	for k, v := range list {
		var perIds []int64
		err := m.db.Table("rule").Cols("v1").Where("ptype = ? and v0=? and v2=?", "p", v.ID, "role").Find(&perIds)
		if err != nil {
			return nil, 0, err
		}
		list[k].PermissionIds = perIds
	}
	return list, count, err
}

// 根据code获取列表
func (m *RoleModel) GetListByCode(code string) ([]Role, error) {
	var list []Role
	err := m.db.Where("code = ?", code).Find(&list)
	return list, err
}

// 根据id获取
func (m *RoleModel) GetByID(id int64) (Role, error) {
	var role Role
	_, err := m.db.ID(id).Get(&role)
	return role, err

}

// 新增
func (m *RoleModel) CreateRole(role *Role) (int64, error) {
	return m.db.Insert(role)
}

// 修改
func (m *RoleModel) UpdateRole(role *Role) (int64, error) {
	return m.db.ID(role.ID).MustCols("enable").Update(role)
}

// 删除
func (m *RoleModel) DeleteRole(id int64) (int64, error) {
	return m.db.ID(id).Delete(&Role{})
}
