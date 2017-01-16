package tasks

import (
	"time"

	"github.com/InnovaCo/serve-runner/config"
	"github.com/InnovaCo/serve-runner/logger"
)

type sleep struct {
	Name     string
	interval time.Duration
}

func (p *sleep) Init() error {
	p.interval = 1 * time.Second
	if interval, ok := config.GetInstance().Task(p.Name)["interval"]; ok {
		p.interval = time.Duration(interval.(float64)) * time.Second
	}
	return nil
}

func (p sleep) Run(args Args, context Context, exchanger ResultExchanger) {
	logger.Log.Debugf("%T: Run with args -- %v, context -- %v", p, args, context)
	if val, ok := context.Kwargs["interval"]; ok {
		p.interval = time.Duration(val.(float64)) * time.Second
	}
	time.Sleep(p.interval)
	exchanger.Dispatch(Result{Code: 0, Args: Args{}, Name: p.Name, Err: nil})
}
