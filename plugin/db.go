package plugin

import (
	"context"
	"encoding/json"
	"log"

	extism "github.com/extism/go-sdk"
	"xorm.io/xorm"
)

type DBConfig struct {
	Sql string `json:"sql"`
}

func DBQuery(db *xorm.Engine) extism.HostFunction {
	functions := extism.NewHostFunctionWithStack("db_query", func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
		// 从栈中获取参数
		ptr := stack[0]
		// 从指针中获取字符串
		str, err := p.ReadString(ptr)
		if err != nil {
			log.Panicf("ReadString: %v", err)
			return
		}
		// 打印字符串
		data, err := db.QueryInterface(str)
		if err != nil {
			log.Panicf("查询错误: %v", err)
			return
		}
		jdata, err := json.Marshal(data)
		if err != nil {
			log.Panicf("json编码错误: %v", err)
			return
		}
		// 将结果写入指针
		stack[0], err = p.WriteBytes(jdata)
	}, []extism.ValueType{extism.ValueTypePTR}, []extism.ValueType{extism.ValueTypePTR})
	return functions
}
func DBQueryOne(db *xorm.Engine) extism.HostFunction {
	functions := extism.NewHostFunctionWithStack("db_query_one", func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
		// 从栈中获取参数
		ptr := stack[0]
		// 从指针中获取字符串
		str, err := p.ReadString(ptr)
		if err != nil {
			log.Panicf("ReadString: %v", err)
			return
		}
		// 打印字符串
		data, err := db.QueryInterface(str)
		if err != nil {
			log.Panicf("查询错误: %v", err)
			return
		}
		d := map[string]any{}
		if len(data) > 0 {
			d = data[0]
		}
		jdata, err := json.Marshal(d)
		if err != nil {
			log.Panicf("json编码错误: %v", err)
			return
		}
		// 将结果写入指针
		stack[0], err = p.WriteBytes(jdata)
	}, []extism.ValueType{extism.ValueTypePTR}, []extism.ValueType{extism.ValueTypePTR})
	return functions
}
func DBExec(db *xorm.Engine) extism.HostFunction {
	functions := extism.NewHostFunctionWithStack("db_exec", func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
		// 从栈中获取参数
		ptr := stack[0]
		// 从指针中获取字符串
		str, err := p.ReadString(ptr)
		if err != nil {
			log.Panicf("ReadString: %v", err)
			return
		}
		// 打印字符串
		data, err := db.QueryInterface(str)
		if err != nil {
			log.Panicf("查询错误: %v", err)
			return
		}
		d := map[string]any{}
		if len(data) > 0 {
			d = data[0]
		}
		jdata, err := json.Marshal(d)
		if err != nil {
			log.Panicf("json编码错误: %v", err)
			return
		}
		// 将结果写入指针
		stack[0], err = p.WriteBytes(jdata)
	}, []extism.ValueType{extism.ValueTypePTR}, []extism.ValueType{extism.ValueTypePTR})
	return functions
}
