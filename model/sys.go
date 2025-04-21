package model

import "xorm.io/xorm"

type SysConfig struct {
	ConfigName  string `json:"config_name" xorm:"config_name"`
	ConfigKey   string `json:"config_key" xorm:"config_key"`
	ConfigValue string `json:"config_value" xorm:"config_value"`
	ConfigType  string `json:"config_type" xorm:"config_type"` // 设置类型
}

// TableName 表名称
func (*SysConfig) TableName() string {
	return "sys_config"
}

type SysConfigModel struct {
	db *xorm.Engine
}

func NewSysConfig(db *xorm.Engine) *SysConfigModel {
	return &SysConfigModel{db: db}
}

// 根据key获取配置
func (m *SysConfigModel) GetConfigByKey(key string) (*SysConfig, error) {
	var config SysConfig
	_, err := m.db.Where("config_key = ?", key).Get(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// 获取所有配置
func (m *SysConfigModel) GetAllConfig() ([]SysConfig, error) {
	var config []SysConfig
	err := m.db.Find(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// 更新
func (m *SysConfigModel) UpdateConfig(config *SysConfig) error {
	_, err := m.db.MustCols("config_value").Update(config)
	if err != nil {
		return err
	}
	return nil
}
