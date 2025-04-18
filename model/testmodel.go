package model

import "xorm.io/xorm"

type MacVod struct {
	db      *xorm.Engine
	VodId   int    `json:"vod_id"`
	VodName string `json:"vod_name"`
}

func NewMacVod(db *xorm.Engine) *MacVod {
	return &MacVod{db: db}
}
func (this *MacVod) TableName() string {
	return "mac_vod"
}
func (this *MacVod) GetOne() (*MacVod, error) {
	m := &MacVod{}
	_, err := this.db.Get(m)
	return m, err
}
