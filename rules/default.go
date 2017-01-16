package rules

import (
	"encoding/json"

	"github.com/hopkinsth/go-ruler"

	"github.com/InnovaCo/serve-runner/logger"
)

func init() {
	registry("default", Ruler(&defaultRule{}))
}

type defaultRule struct {
}

func (p *defaultRule) Init() error {
	return nil
}

func (p *defaultRule) Assert(rule map[string]interface{}, data interface{}) bool {
	r, err := json.Marshal([]map[string]interface{}{rule})
	if err != nil {
		logger.Log.Errorf("Invalid rule %v", rule)
		return false
	}
	logger.Log.Debugf("Rule %v", string(r))
	logger.Log.Debugf("Data for assert %v", data)

	if engine, err := ruler.NewRulerWithJson(r); err != nil {
		logger.Log.Errorf("Rule error %v", err)
		return false
	} else {
		result := engine.Test(data.(map[string]interface{}))
		logger.Log.Debugf("Rule result %v", result)
		return result
	}

}
