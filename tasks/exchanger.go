package tasks

import (
	"fmt"
)

func init() {
	registry("default", ResultExchanger(&defaultExchange{}))
}

type defaultExchange struct {
}

func (p defaultExchange) Init() error {
	channel = make(chan Result, 10)
	return nil
}

func (p defaultExchange) Dispatch(result Result) error {
	channel <- result
	return nil
}

func (p defaultExchange) Reception() (Result, error) {
	if len(channel) != 0 {
		result := <-channel
		return result, nil
	}
	return Result{}, fmt.Errorf("channel %v is empty", channel)
}
