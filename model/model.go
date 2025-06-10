package model

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"xorm.io/xorm"
)

// model模块
var ModelModule = fx.Module("modelModule", fx.Provide(NewModel))

type ModelResult struct {
	fx.Out
	Models *Models
}
type ModelParams struct {
	fx.In
	*Models
}
type Models struct {
	Test            *MacVod
	PermissionModel *PermissionModel
	UserModel       *UserModel
	RoleModel       *RoleModel
	RuleModel       *RuleModel
	SysConfigModel  *SysConfigModel
	MediasModel     *MediasModel
	UdoObjectModel  *UdoObjectModel
	UdoFieldModel   *UdoFieldModel
	UdoDataModel    *UdoDataModel
}

func NewModel(db *xorm.Engine, config *viper.Viper) (ModelResult, error) {
	return ModelResult{Models: &Models{
		Test:            NewMacVod(db),
		PermissionModel: NewPermission(db),
		UserModel:       NewUser(db, config),
		RoleModel:       NewRole(db),
		RuleModel:       NewRule(db),
		SysConfigModel:  NewSysConfig(db),
		MediasModel:     NewMedias(db),
		UdoObjectModel:  NewUdoObject(db),
		UdoFieldModel:   NewUdoField(db),
		UdoDataModel:    NewUdoData(db),
	}}, nil
}
