package sensors

import (
	"time"

	"github.com/InnovaCo/serve-runner/config"
	"github.com/InnovaCo/serve-runner/logger"
)

//func init() {
//	SensorsPool["timeSensor"] = Sensor(&timeSensor{})
//}

type timeSensor struct {
	Name      string
	interval  time.Duration
	exchanger TriggerExchanger
}

func (p *timeSensor) Setup() error {
	p.Name = "timeSensor"
	p.interval = 1 * time.Second
	if interval, ok := config.GetInstance().Sensor(p.Name)["interval"]; ok {
		p.interval = time.Duration(interval.(float64)) * time.Second
	}
	p.exchanger = GetExchanger()
	logger.Log.Debugf("timeSensor setup: %v, %v", p.Name, p.interval)
	return nil
}

func (p *timeSensor) Run() error {
	for {
		time.Sleep(p.interval)
		p.exchanger.Dispatch(Trigger{Type: p.Name, Payload: map[string]interface{}{"time": time.Now()}})
	}
	logger.Log.Debug("timeSensor Complete")
	return nil
}
