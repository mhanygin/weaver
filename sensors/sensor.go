package sensors

import (
	"sync"

	"github.com/InnovaCo/serve-runner/logger"
)

var SensorsPool = map[string]Sensor{}
var instance TriggerExchanger
var once sync.Once
var channel chan Trigger
var TriggerExchangerCtrl = "default"
var exchangerMap = map[string]TriggerExchanger{}

func registry(name string, exchanger TriggerExchanger) {
	exchangerMap[name] = exchanger
}

func GetExchanger() TriggerExchanger {
	once.Do(func() {
		if _, ok := exchangerMap[TriggerExchangerCtrl]; !ok {
			logger.Log.Info("Create default sensor exchange controller")
			TriggerExchangerCtrl = "default"
		}
		instance = exchangerMap[TriggerExchangerCtrl]
		instance.Init()
	})
	return instance
}

type Trigger struct {
	Type    string
	Payload map[string]interface{}
}

type TriggerExchanger interface {
	Init() error
	Dispatch(trigger Trigger) error
	Reception() (Trigger, error)
}

type Sensor interface {
	Setup() error
	Run() error
}
