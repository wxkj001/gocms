package model

import (
	"time"

	"xorm.io/xorm"
)

type Medias struct {
	ID          int64     `json:"id" xorm:"id"`
	MediaUrl    string    `json:"media_url" xorm:"media_url"`
	MediaName   string    `json:"media_name" xorm:"media_name"`
	ContentType string    `json:"content_type" xorm:"content_type"`
	Size        int64     `json:"size" xorm:"size"`
	Createdat   time.Time `json:"CreatedAt" xorm:"created"`
}

// TableName 表名称
func (*Medias) TableName() string {
	return "medias"
}

type MediasModel struct {
	db *xorm.Engine
}

func NewMedias(db *xorm.Engine) *MediasModel {
	return &MediasModel{db: db}
}

// 获取列表带分页
func (m *MediasModel) GetList(page, limit int) ([]Medias, int64, error) {
	var list []Medias
	count, err := m.db.Limit(limit, (page-1)*limit).FindAndCount(&list)
	return list, count, err
}

// 新增
func (m *MediasModel) Add(role *Medias) (int64, error) {
	return m.db.Insert(role)
}

// 根据id获取
func (m *MediasModel) GetByID(id int64) (Medias, error) {
	var role Medias
	_, err := m.db.ID(id).Get(&role)
	return role, err
}
