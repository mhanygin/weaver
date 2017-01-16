package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	config := []byte(`---
sensors:
  timeSensor:
    interval: 30
  hookSensor:
    port: 8081
    url: "/github/st2"
    api-key: "123456789"

tasks:
  github:
    re: "https://github.com/*"
    host: "https://api.github.com"
    name: "test.yml"
    sys: "github"
    id: "username"
    url_field: "url"

  gitlab:
    re: "https://gitlab.ru/*"
    host: "https://gitlab.ru/api/v3/projects/"
    name: "test.yml"
    sys: "gitlab"
    id: "username"
    url_field: "git_http_url"

activeFlows:
  - sensor: "timeSensor"
    rules:
      - comparator: "eq"
        path: "name"
        value: "value"
      - comparator: "neq"
        path: "name1"
        value: "value1"
    tasks:
      - sleep:
          interval: 10
          times: 100
      - serve:
          purge: false

  - sensor: "timeSensor"
    rules: []
    tasks:
      - sleep: {}
      - serve:
          purge: false

  - sensor: "hookSensor"
    rules: []
    tasks:
      - sleep: {}
      - serve: {}
      - sleep: {}`)

	if err := ioutil.WriteFile("/tmp/test", config, 0644); err != nil {
		t.Error("Error file not create")
		t.Fail()
	}

	defer os.Remove("/tmp/test")

	cfg := defaultConfig{}
	if err := cfg.Init("/tmp/test"); err != nil {
		t.Error(err)
		t.Fail()
	}
	fmt.Printf("%v\n", cfg.Sensor("hookSensor"))
	if !reflect.DeepEqual(cfg.Sensor("hookSensor"), map[string]interface{}{
		"api-key": "123456789",
		"port":    8081,
		"url":     "/github/st2"}) {
		t.Fail()
	}

	fmt.Printf("%v\n", cfg.Task("gitlab"))
	if !reflect.DeepEqual(cfg.Task("gitlab"), map[string]interface{}{
		"re":        "https://gitlab.ru/*",
		"host":      "https://gitlab.ru/api/v3/projects/",
		"name":      "manifest.yml",
		"sys":       "gitlab",
		"id":        "username",
		"url_field": "git_http_url"}) {
		t.Fail()
	}
	for _, flow := range cfg.ActiveFlows("timeSensor") {
		fmt.Printf("%v\n", flow.Rules())
		for _, task := range flow.Tasks() {
			fmt.Printf("%v\n", task.Task())
			fmt.Printf("%v\n", task.Context())
		}
	}
}
