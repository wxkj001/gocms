package model

import (
	"errors"
	"gocms/utils"
	"time"

	"github.com/spf13/viper"
	"xorm.io/xorm"
)

type User struct {
	ID         int64     `json:"id" xorm:"id autoincr index pk notnull unique"`
	Username   string    `json:"username" xorm:"username"`
	Password   string    `json:"password" xorm:"password"`
	Enable     bool      `json:"enable" xorm:"enable"`
	Createtime time.Time `json:"createTime" xorm:"created"`
	Updatetime time.Time `json:"updateTime" xorm:"updated"`
}

// TableName 表名称
func (*User) TableName() string {
	return "users"
}

type UserModel struct {
	db     *xorm.Engine
	config *viper.Viper
	User
}
type UserGroup struct {
	User    `xorm:"extends"`
	Profile `xorm:"extends" json:"profile"`
	Role    []RuleRoleGroup `xorm:"-" json:"roles"`
}

func (UserGroup) TableName() string {
	return "users"
}

type UserProfileGroup struct {
	ID      int64 `json:"id" xorm:"id autoincr"`
	User    `xorm:"extends"`
	Profile `xorm:"extends"`
	Role    []RuleRoleGroup `xorm:"-" json:"roles"`
}

func (UserProfileGroup) TableName() string {
	return "users"
}
func NewUser(db *xorm.Engine, config *viper.Viper) *UserModel {
	return &UserModel{db: db, config: config}
}

// 获取用户列表
func (m *UserModel) GetUserList(page, pageSize int) (int64, []UserProfileGroup, error) {
	var users []UserProfileGroup
	count, err := m.db.Join("LEFT", "profile", "users.id=profile.user_id").Limit(pageSize, (page-1)*pageSize).FindAndCount(&users)
	rw := "role.id=rule.v1"
	if m.config.GetString("db.driver") == "postgres" {
		rw = "role.id=CAST(rule.v1 AS INTEGER)"
	}
	for k, v := range users {
		roles := []RuleRoleGroup{}
		err = m.db.Where("ptype =? and v0=?", "g", utils.ToString(v.ID)).
			Join("INNER", "role", rw).
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
		Join("LEFT", "profile", "users.id=profile.user_id").
		Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("no user")
	}
	roles := []RuleRoleGroup{}
	rw := "role.id=rule.v1"
	if m.config.GetString("db.driver") == "postgres" {
		rw = "role.id=CAST(rule.v1 AS INTEGER)"
	}
	err = m.db.Where("ptype =? and v0=?", "g", user.User.ID).
		Join("INNER", "role", rw).
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
	defer session.Close()
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
	defer session.Close()
	session.Begin()
	defer session.Rollback()
	_, err := session.ID(id).Update(&user.User)
	if err != nil {
		return err
	}
	user.Profile.Userid = user.User.ID
	user.Profile.Nickname = user.User.Username
	_, err = session.Where("user_id=?", id).Update(user.Profile)
	if err != nil {
		return err
	}
	err = session.Commit()
	return err
}

// 删除用户
func (m *UserModel) DeleteUser(id int64) error {
	session := m.db.NewSession()
	defer session.Close()
	session.Begin()
	_, err := session.Delete(&User{ID: id})
	if err != nil {
		return err
	}
	_, err = session.Delete(&Profile{Userid: id})
	err = session.Commit()
	return err
}
