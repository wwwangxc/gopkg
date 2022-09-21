package redis

import (
	redigo "github.com/gomodule/redigo/redis"
)

// Scan copies from src to the values pointed at by dest.
//
// Scan uses RedisScan if available otherwise:
//
// The values pointed at by dest must be an integer, float, boolean, string,
// []byte, interface{} or slices of these types. Scan uses the standard strconv
// package to convert bulk strings to numeric and boolean types.
//
// If a dest value is nil, then the corresponding src value is skipped.
//
// If a src element is nil, then the corresponding dest value is not modified.
//
// To enable easy use of Scan in a loop, Scan returns the slice of src
// following the copied values.
//
// See: https://github.com/gomodule/redigo/blob/master/redis/scan.go
func Scan(src []interface{}, dest ...interface{}) ([]interface{}, error) {
	return redigo.Scan(src, dest)
}

// ScanSlice scans src to the slice pointed to by dest.
//
// If the target is a slice of types which implement Scanner then the custom
// RedisScan method is used otherwise the following rules apply:
//
// The elements in the dest slice must be integer, float, boolean, string, struct
// or pointer to struct values.
//
// Struct fields must be integer, float, boolean or string values. All struct
// fields are used unless a subset is specified using fieldNames.
//
// See: https://github.com/gomodule/redigo/blob/master/redis/scan.go
func ScanSlice(src []interface{}, dest interface{}, fieldNames ...string) error {
	return redigo.ScanSlice(src, dest, fieldNames...)
}

// ScanStruct scans alternating names and values from src to a struct. The
// HGETALL and CONFIG GET commands return replies in this format.
//
// ScanStruct uses exported field names to match values in the response. Use
// 'redis' field tag to override the name:
//
//      Field int `redis:"myName"`
//
// Fields with the tag redis:"-" are ignored.
//
// Each field uses RedisScan if available otherwise:
// Integer, float, boolean, string and []byte fields are supported. Scan uses the
// standard strconv package to convert bulk string values to numeric and
// boolean types.
//
// If a src element is nil, then the corresponding field is not modified.
//
// See: https://github.com/gomodule/redigo/blob/master/redis/scan.go
func ScanStruct(src []interface{}, dest interface{}) error {
	return redigo.ScanStruct(src, dest)
}
