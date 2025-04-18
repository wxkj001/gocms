package model

import "xorm.io/xorm"

type Rule struct {
	Ptype string `json:"ptype" xorm:"ptype"`
	V0    string `json:"v0" xorm:"v0"`
	V1    string `json:"v1" xorm:"v1"`
	V2    string `json:"v2" xorm:"v2"`
	V3    string `json:"v3" xorm:"v3"`
	V4    string `json:"v4" xorm:"v4"`
	V5    string `json:"v5" xorm:"v5"`
}

// TableName 表名称
func (*Rule) TableName() string {
	return "rule"
}

type RuleRoleGroup struct {
	Rule `xorm:"extends"`
	Role `xorm:"extends"`
}

func (*RuleRoleGroup) TableName() string {
	return "rule"
}

type RuleModel struct {
	db *xorm.Engine
	Rule
}

func NewRule(db *xorm.Engine) *RuleModel {
	return &RuleModel{db: db}
}

// 删除权限
func (this *RuleModel) Delete(rule *Rule) error {
	_, err := this.db.Delete(rule)
	return err
}

// 获取权限IDs
func (this *RuleModel) GetRuleIdsByRoleId(roleId int64) ([]int64, error) {
	var perIds []int64
	err := this.db.Table("rule").Cols("v1").Where("ptype = ? and v0=? and v2=?", "p", roleId, "role").Find(&perIds)
	if err != nil {
		return nil, err
	}
	return perIds, nil
}
