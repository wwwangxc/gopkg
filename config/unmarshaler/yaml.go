package unmarshaler

import (
	"gopkg.in/yaml.v3"
)

func init() {
	register(&YAML{})
}

// YAML yaml unmarshaler
type YAML struct{}

// Unmarshal unmarshal by yaml
func (y *YAML) Unmarshal(in []byte, out interface{}) error {
	in = []byte(expandEnv(string(in)))
	return yaml.Unmarshal(in, out)
}

// Name unmarshal name
func (y *YAML) Name() string {
	return "yaml"
}
