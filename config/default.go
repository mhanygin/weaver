package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/InnovaCo/serve-runner/logger"
	"github.com/Jeffail/gabs"
	"github.com/ghodss/yaml"
)

func init() {
	registry("default", Config(&defaultConfig{}))
}

type defaultConfig struct {
	config gabs.Container
	path   string
}

func (p *defaultConfig) Init(path string) error {
	p.path = path
	logger.Log.Debugf("config file: %v", p.path)

	data, err := ioutil.ReadFile(p.path)
	if err != nil {
		return fmt.Errorf("file `%s` not found: %v", p.path, err)
	}

	if jsonData, err := yaml.YAMLToJSON(data); err != nil {
		return fmt.Errorf("Error on parse file `%s`: %v!", p.path, err)
	} else {
		g, _ := gabs.ParseJSON(jsonData)
		p.config = *g
	}
	return nil
}

func (p *defaultConfig) Sensor(name string) map[string]interface{} {
	if p.config.Exists("sensors", name) {
		if val, ok := p.config.Path(fmt.Sprintf("sensors.%s", name)).Data().(map[string]interface{}); ok {
			return val
		}
	}
	return map[string]interface{}{}
}

func (p *defaultConfig) Task(name string) map[string]interface{} {
	if p.config.Exists("tasks", name) {
		if val, ok := p.config.Path(fmt.Sprintf("tasks.%s", name)).Data().(map[string]interface{}); ok {
			return val
		}
	}
	return map[string]interface{}{}
}

func (p *defaultConfig) ActiveFlows(sensorType string) []FlowConfig {
	flowsConfig := []FlowConfig{}
	flows, _ := p.config.S("activeFlows").Children()
	for _, flow := range flows {
		if strings.Compare(flow.Path("sensor").Data().(string), sensorType) != 0 {
			continue
		}
		fc := configSection{Section: *flow, Key: sensorType}
		fc.Init()
		flowsConfig = append(flowsConfig, FlowConfig(&fc))
	}
	return flowsConfig
}

type configSection struct {
	Section     gabs.Container
	Key         string
	tasksConfig []TaskConfig
}

func (p *configSection) Init() {
	p.tasksConfig = []TaskConfig{}
	taskSections, _ := p.Section.S("tasks").Children()
	for _, task := range taskSections {
		names, _ := task.S().ChildrenMap()
		for name, _ := range names {
			p.tasksConfig = append(p.tasksConfig, TaskConfig(&configSection{Section: *task, Key: name}))
			break
		}
	}
}

func (p configSection) Sensor() string {
	return p.Key
}

func (p configSection) Tasks() []TaskConfig {
	return p.tasksConfig
}

func (p configSection) Rules() []map[string]interface{} {
	logger.Log.Debugf("%v", p.Section.Path("rules"))

	if p.Section.Exists("rules") {
		if val, ok := p.Section.Path("rules").Data().([]map[string]interface{}); ok {
			return val
		}
	}
	return []map[string]interface{}{}
}

func (p configSection) Context() map[string]interface{} {
	if val, ok := p.Section.Path(p.Key).Data().(map[string]interface{}); ok {
		return val
	}
	return map[string]interface{}{}
}

func (p configSection) Task() string {
	return p.Key
}
