package model

import (
	"xorm.io/xorm"
)

type Profile struct {
	ID       int64  `json:"id" xorm:"id autoincr"`
	Gender   int64  `json:"gender" xorm:"gender"`
	Avatar   string `json:"avatar" xorm:"avatar"`
	Address  string `json:"address" xorm:"address"`
	Email    string `json:"email" xorm:"email"`
	Userid   int64  `json:"userId" xorm:"userId index"`
	Nickname string `json:"nickName" xorm:"nickName"`
}

// TableName 表名称
func (*Profile) TableName() string {
	return "profile"
}

type ProfileModel struct {
	db *xorm.Engine
	Profile
}

func NewProfile(db *xorm.Engine) *ProfileModel {
	return &ProfileModel{db: db}
}
