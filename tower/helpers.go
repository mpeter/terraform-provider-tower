package tower

import (
	"encoding/json"
	"fmt"

	"github.com/go-yaml/yaml"
)

func normalizeJson(s interface{}) string {
	if s == nil || s == "" {
		return ""
	}
	var j interface{}
	err := json.Unmarshal([]byte(s.(string)), &j)
	if err != nil {
		return fmt.Sprintf("Error parsing JSON: %s", err)
	}
	b, _ := json.Marshal(j)
	return string(b[:])
}

func normalizeYaml(s interface{}) string {
	if s == nil || s == "" {
		return ""
	}
	var j interface{}
	err := yaml.Unmarshal([]byte(s.(string)), &j)
	if err != nil {
		return fmt.Sprintf("Error parsing YAML: %s", err)
	}
	b, _ := yaml.Marshal(j)
	return string(b[:])
}
