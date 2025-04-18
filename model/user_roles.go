package model

import "xorm.io/xorm"

type UserRolesRole struct {
	Userid int64 `json:"userId" xorm:"userId"`
	Roleid int64 `json:"roleId" xorm:"roleId"`
}

// TableName 表名称
func (*UserRolesRole) TableName() string {
	return "user_roles_role"
}

type UserRolesRoleModel struct {
	db *xorm.Engine
	UserRolesRole
}
type UserRolesRoleGroup struct {
	UserRolesRole `xorm:"extends"`
	Role          `xorm:"extends"`
}

func (UserRolesRoleGroup) TableName() string {
	return "user_roles_role"
}
func NewUserRolesRole(db *xorm.Engine) *UserRolesRoleModel {
	return &UserRolesRoleModel{db: db}
}
