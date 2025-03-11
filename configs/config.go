package configs

import (
	"github.com/xframe-go/x/contracts"
	"github.com/xframe-go/x/utils/xmap"
)

type Config struct {
	cfg *xmap.Map[any, any]
}

func New() *Config {
	return &Config{
		cfg: xmap.NewMap[any, any](true),
	}
}

func (c *Config) Register(config contracts.ConfigValue) {
	c.cfg.Set(config.Symbol(), config)
}

func (c *Config) Get(key any) (any, bool) {
	return c.cfg.Get(key)
}
