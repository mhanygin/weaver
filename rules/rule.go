package rules

import (
	"sync"

	"github.com/InnovaCo/serve-runner/logger"
)

func init() {
	GetInstance()
}

var instance Ruler
var once sync.Once
var RulerCtrl = "default"
var rulerMap = map[string]Ruler{}

func registry(name string, rule Ruler) {
	rulerMap[name] = rule
}

func GetInstance() Ruler {
	once.Do(func() {
		if _, ok := rulerMap[RulerCtrl]; !ok {
			logger.Log.Info("Create default rule controller")
			RulerCtrl = "default"
		}
		instance = rulerMap[RulerCtrl]
		instance.Init()
	})
	return instance
}

type Ruler interface {
	Init() error
	Assert(rule map[string]interface{}, data interface{}) bool
}
