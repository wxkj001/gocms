package model

import (
	"time"

	"xorm.io/xorm"
)

type User struct {
	ID         int64     `json:"id" xorm:"id autoincr"`
	Username   string    `json:"username" xorm:"username"`
	Password   string    `json:"password" xorm:"password"`
	Enable     bool      `json:"enable" xorm:"enable tinyint(1)"`
	Createtime time.Time `json:"createTime" xorm:"createTime"`
	Updatetime time.Time `json:"updateTime" xorm:"updateTime"`
}

// TableName 表名称
func (*User) TableName() string {
	return "user"
}

type UserModel struct {
	db *xorm.Engine
	User
}
type UserGroup struct {
	User    `xorm:"extends"`
	Profile `xorm:"extends" json:"profile"`
	Role    []RuleRoleGroup `xorm:"-" json:"roles"`
}

func (UserGroup) TableName() string {
	return "user"
}

type UserProfileGroup struct {
	User    `xorm:"extends"`
	Profile `xorm:"extends"`
	Role    []RuleRoleGroup `xorm:"-" json:"roles"`
}

func (UserProfileGroup) TableName() string {
	return "user"
}
func NewUser(db *xorm.Engine) *UserModel {
	return &UserModel{db: db}
}

// 获取用户列表
func (m *UserModel) GetUserList() ([]UserProfileGroup, error) {
	var users []UserProfileGroup
	err := m.db.Join("LEFT", "profile", "user.id=profile.userId").Find(&users)
	for k, v := range users {
		roles := []RuleRoleGroup{}
		err = m.db.Where("ptype =? and v0=?", "g", v.User.ID).
			Join("INNER", "role", "role.id=rule.v1").
			Find(&roles)
		if err != nil {
			return nil, err
		}
		users[k].Role = roles
	}
	return users, err
}

// 获取用户信息
func (m *UserModel) GetUserByUsername(username string) (*UserGroup, error) {
	user := new(UserGroup)
	has, err := m.db.Where("username = ?", username).
		Join("LEFT", "profile", "user.id=profile.userId").
		Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	roles := []RuleRoleGroup{}
	err = m.db.Where("ptype =? and v0=?", "g", user.User.ID).
		Join("INNER", "role", "role.id=rule.v1").
		Find(&roles)
	if err != nil {
		return nil, err
	}
	user.Role = roles
	return user, nil
}

// 创建用户
func (m *UserModel) CreateUser(user *UserGroup) error {
	session := m.db.NewSession()
	defer m.db.Close()
	session.Begin()
	defer session.Rollback()
	_, err := session.Insert(user.User)
	if err != nil {
		return err
	}
	user.Profile.Userid = user.User.ID
	_, err = session.Insert(user.Profile)
	if err != nil {
		return err
	}
	err = session.Commit()
	return err
}

// 更新用户
func (m *UserModel) UpdateUser(id int64, user *User) error {
	_, err := m.db.ID(id).Update(user)
	return err
}

// 删除用户
func (m *UserModel) DeleteUser(id int64) error {
	_, err := m.db.ID(id).Delete(&User{})
	return err
}
