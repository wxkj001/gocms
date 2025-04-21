package plugin

import (
	"context"
	"errors"

	extism "github.com/extism/go-sdk"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

type Plugins struct {
	log           *zap.Logger
	config        *viper.Viper
	list          map[string]*extism.Plugin
	HostFunctions []extism.HostFunction
	isAdd         bool
}

func New(config *viper.Viper, log *zap.Logger, db *xorm.Engine) *Plugins {
	p := &Plugins{
		log:           log,
		config:        config,
		list:          make(map[string]*extism.Plugin),
		HostFunctions: []extism.HostFunction{},
	}
	p.AddHostFunction(DBQuery(db)).
		AddHostFunction(DBQueryOne(db)).
		AddHostFunction(DBExec(db))
	return p
}

// 增加HostFunction
func (c *Plugins) AddHostFunction(hostFunction extism.HostFunction) *Plugins {
	if c.isAdd {
		c.log.Error("插件已经注册，不能添加HostFunction")
		return c
	}
	c.HostFunctions = append(c.HostFunctions, hostFunction)
	return c
}

// 使用插件
func (c *Plugins) Use(name string) (*extism.Plugin, error) {
	if plugin, ok := c.list[name]; ok {
		return plugin, nil
	}
	return nil, errors.New("插件不存在")
}

// 注册插件
func (c *Plugins) Add(name string) error {
	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{Path: c.config.GetString("plugin.path") + "/" + name + ".wasm"},
		},
		Config: map[string]string{
			"vowels": "aeiouyAEIOUY",
		},
	}

	ctx := context.Background()
	config := extism.PluginConfig{
		EnableWasi: true,
	}

	plugin, err := extism.NewPlugin(ctx, manifest, config, c.HostFunctions)
	if err != nil {
		return err
	}
	plugin.SetLogger(func(ll extism.LogLevel, s string) {
		c.log.Info("插件日志", zap.String("name", name), zap.String("level", ll.String()), zap.String("message", s))
	})
	c.list[name] = plugin
	return nil
}

// 注销插件
func (c *Plugins) Remove(name string) {
	delete(c.list, name)
}
