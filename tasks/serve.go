package tasks

import (
	"fmt"
	"os"
	"os/exec"
)

type serve struct {
	Name string
}

func (p serve) Init() error {
	return nil
}

func (p serve) Run(args Args, context Context, exchanger ResultExchanger) {
	var manifest string
	var vars string
	var plugin string

	if param, ok := args.Kwargs["manifest"]; ok {
		manifest = param.(string)
	} else {
		manifest = ""
	}

	if param, ok := args.Kwargs["plugin"]; ok {
		plugin = param.(string)
	} else {
		plugin = ""
	}

	if param, ok := args.Kwargs["vars"]; ok {
		for k, v := range param.(map[string]interface{}) {
			vars = fmt.Sprintf("%v --vars %v=%v", vars, k, v.(string))
		}
	} else {
		vars = ""
	}

	cmd := exec.Command(fmt.Sprintf("serve %v %v --manifest=%v", plugin, vars, manifest))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		exchanger.Dispatch(Result{Code: 256, Args: Args{}, Name: p.Name, Err: err})
	} else {
		exchanger.Dispatch(Result{Code: 0, Args: Args{}, Name: p.Name, Err: nil})
	}
}
