package plugin

import (
	"context"
	"errors"

	extism "github.com/extism/go-sdk"
	"github.com/spf13/viper"
)

type Plugins struct {
	config        *viper.Viper
	list          map[string]*extism.Plugin
	HostFunctions []extism.HostFunction
}

func New(config *viper.Viper) *Plugins {
	return &Plugins{
		config:        config,
		list:          make(map[string]*extism.Plugin),
		HostFunctions: []extism.HostFunction{},
	}
}

// 增加HostFunction
func (c *Plugins) AddHostFunction(hostFunction extism.HostFunction) {
	c.HostFunctions = append(c.HostFunctions, hostFunction)
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
	config := extism.PluginConfig{}

	plugin, err := extism.NewPlugin(ctx, manifest, config, c.HostFunctions)
	if err != nil {
		return err
	}
	c.list[name] = plugin
	return nil
}

// 注销插件
func (c *Plugins) Remove(name string) {
	delete(c.list, name)
}
