package main

import (
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/InnovaCo/serve-runner/config"
	"github.com/InnovaCo/serve-runner/manager"
	"github.com/InnovaCo/serve-runner/rules"
	"github.com/InnovaCo/serve-runner/sensors"
	"github.com/InnovaCo/serve-runner/tasks"
)

var version = "0.2.0"

func main() {
	cfgPath := kingpin.Flag("config", "Path to config.yml file.").Default("config.yml").String()
	configCtrl := kingpin.Flag("config-ctrl", "Config controller").Default("default").String()
	manageCtrl := kingpin.Flag("manage-ctrl", "Manage controller").Default("default").String()
	ruleCtrl := kingpin.Flag("rule-ctrl", "Rule controller").Default("default").String()
	sensorCtrl := kingpin.Flag("sensor-ex-ctrl", "Sensor exchange controller").Default("default").String()
	taskCtrl := kingpin.Flag("task-ex-ctrl", "Task exchange controller").Default("default").String()

	kingpin.Version(version)
	kingpin.Parse()

	config.ConfigPath = *cfgPath
	config.ConfigCtrl = *configCtrl
	manager.ManagerCtrl = *manageCtrl
	sensors.TriggerExchangerCtrl = *sensorCtrl
	tasks.TaskExchangerCtrl = *taskCtrl
	rules.RulerCtrl = *ruleCtrl

	m := manager.GetInstance()
	m.Run()
}
