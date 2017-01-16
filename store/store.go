package store

import (
	"sync"
	"time"

	"github.com/InnovaCo/serve-runner/config"
	"github.com/InnovaCo/serve-runner/logger"
	"github.com/InnovaCo/serve-runner/sensors"
	"github.com/InnovaCo/serve-runner/tasks"
)

const TaskInit = -1

var instance Store
var once sync.Once
var StoreCtrl = "default"
var storeMap = map[string]Store{}

func registry(name string, store Store) {
	storeMap[name] = store
}

func GetInstance() Store {
	once.Do(func() {
		if _, ok := storeMap[StoreCtrl]; !ok {
			logger.Log.Info("Create default store controller")
			StoreCtrl = "default"
		}
		instance = storeMap[StoreCtrl]
		instance.Init()
	})
	return instance
}

type InstFlow struct {
	Uuid        string
	Trigger     sensors.Trigger
	Exchanger   tasks.ResultExchanger
	CurrentTask int
	TaskPool    []tasks.Task
	ContextPool []tasks.Context
	Results     []tasks.Result
	lastModify  time.Time
}

type Store interface {
	Init() error
	AddFlow(trigger sensors.Trigger, cfg config.FlowConfig) (InstFlow, error)
	GetFlows() []InstFlow
	SaveFlow(flow InstFlow) error
}
