package unmarshaler

import (
	"encoding/json"
)

func init() {
	register(&JSON{})
}

// JSON json unmarshaler
type JSON struct{}

// Unmarshal unmarshal by json
func (j *JSON) Unmarshal(in []byte, out interface{}) error {
	in = []byte(expandEnv(string(in)))
	return json.Unmarshal(in, out)
}

// Name unmarshal name
func (j *JSON) Name() string {
	return "json"
}
