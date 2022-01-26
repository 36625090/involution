package authorities

import (
	"encoding/json"
	"github.com/go-various/redisplus"
	"github.com/google/uuid"
)

type redisTokenHandler struct {
	settings *Settings
	redis redisplus.RedisCli
}

func NewRedisTokenHandler(settings *Settings, config *redisplus.Config)(TokenHandler, error){
	view, err := redisplus.NewRedisCli(config, config.KeyPrefix)
	if err != nil {
		return nil, err
	}
	return &redisTokenHandler{
		settings: settings,
		redis: view,
	}, nil
}

func (r *redisTokenHandler) GenerateToken(auth *Authorized) (string, error) {

	data, err := json.Marshal(auth)
	if err != nil {
		return "", err
	}

	token := uuid.New().String()
	if err := r.redis.Set(token, data, r.settings.Timeout.String()); err != nil {
		return "", err
	}
	return token, nil
}

func (r *redisTokenHandler) ParseToken(token string) (*Authorized, error) {
	data, err := r.redis.Get(token)
	if err != nil {
		return nil, err
	}
	var authorized Authorized
	if err := json.Unmarshal(data, &authorized); err != nil {
		return nil, err
	}

	return &authorized, nil
}

