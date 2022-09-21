package unmarshaler

import (
	"github.com/BurntSushi/toml"
)

func init() {
	register(&TOML{})
}

// TOML toml unmarshaler
type TOML struct{}

// Unmarshal unmarshal by toml
func (t *TOML) Unmarshal(in []byte, out interface{}) error {
	in = []byte(expandEnv(string(in)))
	return toml.Unmarshal(in, out)
}

// Name unmarshal name
func (t *TOML) Name() string {
	return "toml"
}
