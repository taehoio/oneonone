package config

import "github.com/sirupsen/logrus"

func MockLogger() *logrus.Logger {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logger := logrus.StandardLogger()
	return logger
}

func MockConfig() Config {
	return NewConfig(MockSetting(), MockLogger())
}
