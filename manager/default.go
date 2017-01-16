package manager

import (
	"github.com/InnovaCo/serve-runner/config"
	"github.com/InnovaCo/serve-runner/logger"
	"github.com/InnovaCo/serve-runner/rules"
	"github.com/InnovaCo/serve-runner/sensors"
	"github.com/InnovaCo/serve-runner/store"
	"github.com/InnovaCo/serve-runner/tasks"
)

func init() {
	registry("default", Manager(&defaultManage{}))
}

type defaultManage struct {
	triggerExchanger sensors.TriggerExchanger
}

func (p *defaultManage) Init() error {
	p.triggerExchanger = sensors.GetExchanger()
	return nil
}

func (p *defaultManage) validateResult(result tasks.Result) bool {
	if result.Err != nil {
		logger.Log.Debugf("Flow error: %v (return code: %v), error %v", result.Name, result.Code, result.Err)
	}
	return true
}

func (p *defaultManage) createArgs(trigger sensors.Trigger, results []tasks.Result) tasks.Args {
	return results[len(results)-1].Args
}

func (p *defaultManage) createFlows(trigger sensors.Trigger) error {
	for _, flow := range config.GetInstance().ActiveFlows(trigger.Type) {
		create := true
		for _, rule := range flow.Rules() {
			if create = rules.GetInstance().Assert(rule, trigger.Payload); !create {
				break
			}
		}
		if !create {
			continue
		}
		_, err := store.GetInstance().AddFlow(trigger, flow)
		return err
	}
	return nil
}

func (p *defaultManage) runTask(flow store.InstFlow) error {
	if flow.CurrentTask >= len(flow.Results) {
		return nil
	}
	if flow.CurrentTask == store.TaskInit {
		flow.CurrentTask = flow.CurrentTask + 1
		go flow.TaskPool[flow.CurrentTask].Run(tasks.Args{flow.Trigger.Payload},
			flow.ContextPool[flow.CurrentTask],
			flow.Exchanger)
	} else {
		flow.CurrentTask = flow.CurrentTask + 1
		go flow.TaskPool[flow.CurrentTask].Run(p.createArgs(flow.Trigger, flow.Results),
			flow.ContextPool[flow.CurrentTask],
			flow.Exchanger)
	}
	logger.Log.Debugf("Run task %v in flow %v", flow.CurrentTask, flow.Uuid)
	store.GetInstance().SaveFlow(flow)
	return nil
}

func (p *defaultManage) retryTask(flow store.InstFlow) {
	logger.Log.Debugf("Retry task %d in flow %v", flow.CurrentTask, flow.Uuid)
	if flow.CurrentTask > store.TaskInit {
		flow.CurrentTask = flow.CurrentTask - 1
		store.GetInstance().SaveFlow(flow)
	}
}

func (p *defaultManage) Run() error {
	for sensorType, sensor := range sensors.SensorsPool {
		if len(config.GetInstance().ActiveFlows(sensorType)) != 0 {
			if err := sensor.Setup(); err != nil {
				logger.Log.Errorf("Error: %v", err)
				continue
			}
			go sensor.Run()
		}
	}
	for {
		if trigger, err := p.triggerExchanger.Reception(); err == nil {
			err := p.createFlows(trigger)
			if err != nil {
				logger.Log.Errorf("Create flow error: %v", err)
			}
		}
		for _, flow := range store.GetInstance().GetFlows() {
			result, err := flow.Exchanger.Reception()
			if err != nil {
				continue
			}
			if p.validateResult(result) {
				flow.Results = append(flow.Results, result)
				err := store.GetInstance().SaveFlow(flow)
				if err != nil {
					logger.Log.Errorf("Result %v save error: %v", flow, err)
				}
			} else {
				p.retryTask(flow)
			}
		}
		for _, flow := range store.GetInstance().GetFlows() {
			if err := p.runTask(flow); err != nil {
				logger.Log.Errorf("Error: %v", err)
			}
		}
	}

	return nil
}
