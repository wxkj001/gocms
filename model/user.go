package model

import (
	"time"

	"xorm.io/xorm"
)

type User struct {
	ID         int64     `json:"id" xorm:"id autoincr index pk"`
	Username   string    `json:"username" xorm:"username"`
	Password   string    `json:"password" xorm:"password"`
	Enable     bool      `json:"enable" xorm:"enable tinyint(1)"`
	Createtime time.Time `json:"createTime" xorm:"created"`
	Updatetime time.Time `json:"updateTime" xorm:"updated"`
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
	ID      int64 `json:"id" xorm:"id autoincr"`
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
func (m *UserModel) GetUserList(page, pageSize int) (int64, []UserProfileGroup, error) {
	var users []UserProfileGroup
	count, err := m.db.Join("LEFT", "profile", "user.id=profile.userId").Limit(pageSize, (page-1)*pageSize).FindAndCount(&users)
	for k, v := range users {
		roles := []RuleRoleGroup{}
		err = m.db.Where("ptype =? and v0=?", "g", v.ID).
			Join("INNER", "role", "role.id=rule.v1").
			Find(&roles)
		if err != nil {
			return 0, nil, err
		}
		users[k].Role = roles
	}
	return count, users, err
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
	session.Begin()
	defer session.Rollback()
	_, err := session.Insert(&user.User)
	if err != nil {
		return err
	}
	user.Profile.Userid = user.User.ID
	user.Profile.Nickname = user.User.Username
	_, err = session.Insert(user.Profile)
	if err != nil {
		return err
	}
	err = session.Commit()
	return err
}

// 更新用户
func (m *UserModel) UpdateUser(id int64, user *UserGroup) error {
	session := m.db.NewSession()
	session.Begin()
	defer session.Rollback()
	_, err := session.ID(id).Update(&user.User)
	if err != nil {
		return err
	}
	user.Profile.Userid = user.User.ID
	user.Profile.Nickname = user.User.Username
	_, err = session.Where("userId=?", id).Update(user.Profile)
	if err != nil {
		return err
	}
	err = session.Commit()
	return err
}

// 删除用户
func (m *UserModel) DeleteUser(id int64) error {
	session := m.db.NewSession()
	session.Begin()
	_, err := session.Delete(&User{ID: id})
	if err != nil {
		return err
	}
	_, err = session.Delete(&Profile{Userid: id})
	err = session.Commit()
	return err
}
