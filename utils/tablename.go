package utils

import "github.com/spf13/viper"

func GetTableName(name string) string {
	return viper.GetString("db.tablePrefix") + name
}
