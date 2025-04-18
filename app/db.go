package app

import (
	"context"
	"errors"
	"gocms/model"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

var dbModule = fx.Module("dbModule", fx.Provide(NewDB))

func NewDB(config *viper.Viper, log *zap.Logger, lc fx.Lifecycle) (*xorm.Engine, error) {
	log.Debug("DB", zap.String("driver", config.GetString("db.driver")))
	var engine *xorm.Engine
	var err error
	switch config.GetString("db.driver") {
	case "mysql":
		engine, err = newMysqlDB(config, log, lc)
	case "postgres":
		engine, err = newPgsqlDB(config, log, lc)
	}
	if err != nil {
		return nil, err
	}
	if engine == nil {
		return nil, errors.New("db driver not support")
	}
	logger := model.NewLoggerAdapter(log)
	engine.SetLogger(logger)
	if strings.ToUpper(config.GetString("log.level")) == "DEBUG" {
		engine.ShowSQL(true)
	}
	tbMapper := names.NewPrefixMapper(names.SnakeMapper{}, config.GetString("db.tablePrefix"))
	engine.SetTableMapper(tbMapper)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error { return engine.Sync2() },
		OnStop:  func(ctx context.Context) error { return engine.Close() },
	})
	return engine, nil
}
func newMysqlDB(config *viper.Viper, log *zap.Logger, lc fx.Lifecycle) (*xorm.Engine, error) {
	mysqlurl := []string{
		config.GetString("db.user"),
		":",
		config.GetString("db.password"),
		"@tcp(",
		config.GetString("db.host"),
		":",
		config.GetString("db.port"),
		")/",
		config.GetString("db.dbname"),
		"?charset=utf8",
	}
	log.Debug("MysqlDB", zap.String("DSN", strings.Join(mysqlurl, "")))
	engine, err := xorm.NewEngine("mysql", strings.Join(mysqlurl, ""))
	return engine, err
}
func newPgsqlDB(config *viper.Viper, log *zap.Logger, lc fx.Lifecycle) (*xorm.Engine, error) {
	pgurl := []string{
		"postgres://",
		config.GetString("db.user"),
		":",
		config.GetString("db.password"),
		"@",
		config.GetString("db.host"),
		":",
		config.GetString("db.port"),
		"/",
		config.GetString("db.dbname"),
		"?sslmode=disable",
	}
	log.Debug("PgsqlDB", zap.String("DSN", strings.Join(pgurl, "")))
	engine, err := xorm.NewEngine("postgres", strings.Join(pgurl, ""))
	return engine, err
}
