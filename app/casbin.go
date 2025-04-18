package app

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	xormadapter "github.com/casbin/xorm-adapter/v3"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

var casbinModule = fx.Module("casbinModule", fx.Provide(NewCasbin))

func NewCasbin(lc fx.Lifecycle, config *viper.Viper, log *zap.Logger, db *xorm.Engine) (*casbin.Enforcer, error) {
	a, err := xormadapter.NewAdapterByEngineWithTableName(db, "rule", config.GetString("db.tablePrefix"))
	if err != nil {
		return nil, err
	}
	cm, err := model.NewModelFromString(config.GetString("casbin.model"))
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(cm, a)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("casbin start")
			e.LoadPolicy()
			e.SavePolicy()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("casbin stop")
			return nil
		},
	})
	return e, nil
}
