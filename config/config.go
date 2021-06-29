package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
)

func init() {
	Load(true)
}

func Load(isLocal bool) {
	if isLocal { // 使用本地配置
		viper.SetConfigName("local")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			panic(errors.Wrap(err, "ReadInConfig failed"))
		}
	}
}

// Get 获取字符串环境变量
func Get(key string) string {
	if envVar, ok := os.LookupEnv(key); ok {
		return envVar
	}
	return viper.GetString(key)
}

// GetInt 获取环境变量
func GetInt(key string) (int, error) {
	result, err := strconv.Atoi(Get(key))
	if err != nil {
		return 0, errors.Wrap(err, "GetInt failed")
	}
	return result, nil
}

func MustGetInt(key string) int {
	res, err := GetInt(key)
	if err != nil {
		panic(err)
	}
	return res
}

func GetUint32(key string) (uint32, error) {
	result, err := strconv.ParseUint(Get(key), 10, 32)
	if err != nil {
		return 0, errors.Wrap(err, "GetUint32 failed")
	}
	return uint32(result), nil
}

func MustGetUint32(key string) uint32 {
	res, err := GetUint32(key)
	if err != nil {
		panic(err)
	}
	return res
}

func GetBool(key string) (bool, error) {
	res := Get(key)
	switch strings.ToLower(res) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, errors.New("GetBool failed")
	}
}

func MustGetBool(key string) bool {
	res, err := GetBool(key)
	if err != nil {
		panic(err)
	}
	return res
}
