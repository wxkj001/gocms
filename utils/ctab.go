package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// TableNameGetter 定义获取表名的接口
type TableNameGetter interface {
	TableName() string
}

// 创建一个实现了TableNameGetter接口的结构体
type dynamicTable struct {
	value     any
	tableName string
}

// 实现TableName方法
func (d *dynamicTable) TableName() string {
	return d.tableName
}
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = string([]rune(part)[0]-32) + part[1:] // 首字母大写
		}
	}
	return strings.Join(parts, "")
}

// createDynamicStruct 创建动态结构体
func CreateDynamicStruct(tableName string, fieldDefs []map[string]any) (any, error) {
	// 创建字段列表
	fields := make([]reflect.StructField, 0, len(fieldDefs)+1)

	// 添加ID字段（通常作为主键）
	fields = append(fields, reflect.StructField{
		Name: "ID",
		Type: reflect.TypeOf(int64(0)),
		Tag:  reflect.StructTag(`json:"id" xorm:"id pk autoincr"`),
	})

	// 添加其他字段
	for _, fieldDef := range fieldDefs {
		fieldNameS, ok := fieldDef["name"].(string)
		if !ok || fieldNameS == "" {
			return nil, fmt.Errorf("field name is required and must be a string")
		}

		// 首字母大写以确保字段可导出
		// fieldName = fmt.Sprintf("%s%s", string([]rune(fieldName)[0]-32), fieldName[1:])
		fieldName := toCamelCase(fieldNameS)
		fieldType, ok := fieldDef["type"].(string)
		if !ok || fieldType == "" {
			return nil, fmt.Errorf("field type is required and must be a string")
		}

		// 根据字段类型确定反射类型
		var reflectType reflect.Type
		switch fieldType {
		case "string":
			reflectType = reflect.TypeOf("")
		case "int":
			reflectType = reflect.TypeOf(int(0))
		case "int64":
			reflectType = reflect.TypeOf(int64(0))
		case "float":
			reflectType = reflect.TypeOf(float64(0))
		case "bool":
			reflectType = reflect.TypeOf(bool(false))
		case "time":
			reflectType = reflect.TypeOf(time.Time{})
		case "date":
			reflectType = reflect.TypeOf(time.Time{})
		case "datetime":
			reflectType = reflect.TypeOf(time.Time{})
		default:
			return nil, fmt.Errorf("unsupported field type: %s", fieldType)
		}

		// 构建xorm标签
		xormTag := fieldNameS
		if constraints, ok := fieldDef["constraints"].(map[string]interface{}); ok {
			for k, v := range constraints {
				xormTag += " " + k
				if v != nil && v != "" {
					xormTag += fmt.Sprintf("(%v)", v)
				}
			}
		}

		// 创建字段
		fields = append(fields, reflect.StructField{
			Name: fieldName,
			Type: reflectType,
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s" xorm:"%s"`, fieldNameS, xormTag)),
		})
	}

	// 创建结构体类型
	structType := reflect.StructOf(fields)

	// 创建一个新的结构体实例
	structValue := reflect.New(structType).Interface()
	// 创建动态表实例
	dynamic := &dynamicTable{
		value:     structValue,
		tableName: tableName,
	}
	return dynamic, nil
}

// syncTable 同步表结构到数据库
func SyncTable(tableName string, fieldDefs []map[string]any) (string, any, error) {
	structValue, err := CreateDynamicStruct(tableName, fieldDefs)
	if err != nil {
		return "", nil, err
	}
	// 获取动态表实例
	dynamic, ok := structValue.(*dynamicTable)
	if !ok {
		return "", nil, errors.New("invalid dynamic table type")
	}
	return dynamic.tableName, dynamic.value, nil
}
