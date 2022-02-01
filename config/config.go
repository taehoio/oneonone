package config

import (
	"github.com/sirupsen/logrus"
)

type Config interface {
	Setting() Setting
	Logger() *logrus.Logger
}

type DefaultConfig struct {
	setting Setting
	logger  *logrus.Logger
}

func (c *DefaultConfig) Setting() Setting {
	return c.setting
}

func (c *DefaultConfig) Logger() *logrus.Logger {
	return c.logger
}

func NewConfig(setting Setting, logger *logrus.Logger) Config {
	return &DefaultConfig{
		setting: setting,
		logger:  logger,
	}
}
