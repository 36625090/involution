package config

import (
	"github.com/36625090/involution/authorities"
	"github.com/36625090/involution/transport"
	"github.com/go-various/redisplus"
	"github.com/go-various/xorm"
)

type GlobalConfig struct {
	XormConfig    *xorm.Config          `json:"xorm" hcl:"xorm,block"`
	RedisConfig   *redisplus.Config     `json:"redis" hcl:"redis,block"`
	Authorization *authorities.Settings `json:"authorization" hcl:"authorization,block"`
	Transport     *transport.Settings   `json:"transport" hcl:"transport"`
	Extras        Extras                `json:"extras" hcl:"extras,block"`
}

type Extra map[string]interface{}
type Extras map[string]Extra

func (e Extras) GetExtra(key string) Extra {
	return e[key]
}

func (e Extra) GetValue(key string) interface{} {
	return e[key]
}
func (e Extra) GetString(key string) string {
	if val, ok := e[key]; ok {
		return val.(string)
	}
	return ""
}

func (e Extra) GetInt(key string) (int, bool) {
	if val, ok := e[key]; ok {
		return val.(int), true
	}
	return 0, false
}

func (e Extra) GetInt64(key string) (int64, bool) {
	if val, ok := e[key]; ok {
		return val.(int64), true
	}
	return 0, false
}
func (e Extra) GetUInt(key string) (uint, bool) {
	if val, ok := e[key]; ok {
		return val.(uint), true
	}
	return 0, false
}

func (e Extra) GetUInt64(key string) (uint64, bool) {
	if val, ok := e[key]; ok {
		return val.(uint64), true
	}
	return 0, false
}
