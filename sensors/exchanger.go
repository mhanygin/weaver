package sensors

import (
	"fmt"

	"github.com/InnovaCo/serve-runner/logger"
)

func init() {
	registry("default", TriggerExchanger(&defaultExchange{}))
}

type defaultExchange struct {
}

func (p defaultExchange) Init() error {
	channel = make(chan Trigger, 10)
	return nil
}

func (p defaultExchange) Dispatch(trigger Trigger) error {
	logger.Log.Debugf("Dispatch: type -- %v, payload -- %v", trigger.Type, trigger.Payload)
	channel <- trigger
	return nil
}

func (p defaultExchange) Reception() (Trigger, error) {
	if len(channel) != 0 {
		trigger := <-channel
		logger.Log.Debugf("Reception: type -- %v, payload -- %v", trigger.Type, trigger.Payload)
		return trigger, nil
	}
	return Trigger{}, fmt.Errorf("channel %v is empty", channel)
}
