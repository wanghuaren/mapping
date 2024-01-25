package uts

import (
	"gopkg.in/ini.v1"
)

type IConfig interface {
	String(field string) string
	Int(field string) int64
	Bool(field string) bool
}
type config struct {
	ini *ini.File
}

var _config config
var _key string

func InitConf(confPath string, key string) IConfig {
	_config = config{}
	_key = key
	var err error
	_config.ini, err = ini.Load(confPath)
	if err != nil {
		LogF("项目 "+key+" 参数文件读取错误，请检查文件路径:", err.Error())
	}
	return _config
}

func (c config) String(field string) string {
	_cv := _config.ini.Section(_key).Key(field).String()
	return _cv
}

func (c config) Int(field string) int64 {
	_cv, err := _config.ini.Section(_key).Key(field).Int64()
	if err != nil {
		_cv = 0
		LogF(err.Error())
	}
	return _cv
}

func (c config) Bool(field string) bool {
	_cv, err := _config.ini.Section(_key).Key(field).Bool()
	if err != nil {
		_cv = false
		LogF(err.Error())
	}
	return _cv
}
