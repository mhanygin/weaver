package config

import (
	"sync"

	"github.com/InnovaCo/serve-runner/logger"
)

var instance Config
var once sync.Once
var ConfigPath = "config.yml"
var ConfigCtrl = "default"
var configMap = map[string]Config{}

func registry(name string, config Config) {
	configMap[name] = config
}

func GetInstance() Config {
	once.Do(func() {
		if _, ok := configMap[ConfigCtrl]; !ok {
			logger.Log.Info("Create default config controller")
			ConfigCtrl = "default"
		}
		instance = configMap[ConfigCtrl]
		instance.Init(ConfigPath)
	})
	return instance
}

type Config interface {
	Init(path string) error
	Sensor(name string) map[string]interface{}
	Task(name string) map[string]interface{}
	ActiveFlows(sensorType string) []FlowConfig
}

type FlowConfig interface {
	Sensor() string
	Rules() []map[string]interface{}
	Tasks() []TaskConfig
}

type TaskConfig interface {
	Task() string
	Context() map[string]interface{}
}
