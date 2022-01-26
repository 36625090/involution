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
}
