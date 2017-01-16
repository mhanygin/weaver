package sensors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/InnovaCo/serve-runner/config"
	"github.com/InnovaCo/serve-runner/logger"
	"strings"
)

func init() {
	SensorsPool["hookSensor"] = Sensor(&hookSensor{})
}

type hookSensor struct {
	Name      string
	Port      int
	exchanger TriggerExchanger
	apiKey    string
}

func (p *hookSensor) handler(w http.ResponseWriter, r *http.Request) {
	if strings.Compare(p.apiKey, r.FormValue("api-key")) != 0 {
		logger.Log.Debugf("not valid api-key: \"%v\"", r.FormValue("api-key"))
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Log.Errorf("Error in body: %v", err)
		return
	}
	data := map[string]interface{}{}
	if err := json.Unmarshal(body, &data); err != nil {
		logger.Log.Errorf("Parse json: %v", err)
		return
	}
	p.exchanger.Dispatch(Trigger{Type: p.Name, Payload: data})
}

func (p *hookSensor) Setup() error {
	p.Name = "hookSensor"
	p.exchanger = GetExchanger()
	if val, ok := config.GetInstance().Sensor(p.Name)["port"]; ok {
		p.Port = int(val.(float64))
	} else {
		p.Port = 8080
	}

	if val, ok := config.GetInstance().Sensor(p.Name)["url"]; ok {
		http.HandleFunc(val.(string), p.handler)
	} else {
		http.HandleFunc("/", p.handler)
	}

	if val, ok := config.GetInstance().Sensor(p.Name)["api-key"]; ok {
		p.apiKey = val.(string)
	} else {
		p.apiKey = ""
	}
	logger.Log.Debug("hookSensor setup for port ", p.Port)
	return nil
}

func (p *hookSensor) Run() error {
	logger.Log.Debug("hookSensor Run")
	err := http.ListenAndServe(fmt.Sprintf(":%d", p.Port), nil)
	if err != nil {
		logger.Log.Errorf("%v %v", p.Name, err)
	}
	logger.Log.Debug("hookSensor Complete")
	return nil
}
