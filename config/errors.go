package config

import "fmt"

var (
	// ErrUnmarshalerNotExist unmarshaler not exist
	ErrUnmarshalerNotExist = fmt.Errorf("%s: unmarshaler not exist", packageName)

	// ErrConfigNotExist config not exist
	ErrConfigNotExist = fmt.Errorf("%s: config not exist", packageName)
)
