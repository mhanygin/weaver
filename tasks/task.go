package tasks

import (
	"github.com/InnovaCo/serve-runner/logger"
)

var WorkerPool map[string]Worker
var channel chan Result
var TaskExchangerCtrl = "default"
var exchangerMap = map[string]ResultExchanger{}

func registry(name string, exchnger ResultExchanger) {
	exchangerMap[name] = exchnger
}

func init() {
	WorkerPool = make(map[string]Worker)
	RegistryTask("serve", Task(&serve{Name: "serve"}))
	RegistryTask("sleep", Task(&sleep{Name: "sleep"}))
}

func RegistryTask(name string, task Task) {
	if err := task.Init(); err != nil {
		logger.Log.Errorf("Error init task: %v", task)
		return
	}
	WorkerPool[name] = Worker{Name: name, Task: task}
}

func GetExchanger() ResultExchanger {
	var instance ResultExchanger
	if e, ok := exchangerMap[TaskExchangerCtrl]; ok {
		instance = e
	} else {
		logger.Log.Info("Create default store controller")
		instance = exchangerMap["default"]
	}
	instance.Init()
	return instance
}

type ResultExchanger interface {
	Init() error
	Dispatch(result Result) error
	Reception() (Result, error)
}

type Args struct {
	Kwargs map[string]interface{}
}

type Context struct {
	Kwargs map[string]interface{}
}

type Result struct {
	Name string
	Code int
	Err  error
	Args Args
}

type Worker struct {
	Name string
	Task Task
}

type Task interface {
	Init() error
	Run(args Args, context Context, exchanger ResultExchanger)
}
