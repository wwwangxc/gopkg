package redis

import (
	redigo "github.com/gomodule/redigo/redis"
)

/*************************************************************
 ***************************  Int  ***************************
 *************************************************************/

// Int is a helper that converts a command reply to an integer. If err is not
// equal to nil, then Int returns 0, err. Otherwise, Int converts the
// reply to an int as follows:
//
//  Reply type    Result
//  integer       int(reply), nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func Int(reply interface{}, err error) (int, error) {
	return redigo.Int(reply, err)
}

// Ints is a helper that converts an array command reply to a []int.
// If err is not equal to nil, then Ints returns nil, err. Nil array
// items are stay nil. Ints returns an error if an array item is not a
// bulk string or nil.
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func Ints(reply interface{}, err error) ([]int, error) {
	return redigo.Ints(reply, err)
}

// IntMap is a helper that converts an array of strings (alternating key, value)
// into a map[string]int. The HGETALL commands return replies in this format.
// Requires an even number of values in result.
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func IntMap(result interface{}, err error) (map[string]int, error) {
	return redigo.IntMap(result, err)
}

/*************************************************************
 **************************  Int64  **************************
 *************************************************************/

// Int64 is a helper that converts a command reply to 64 bit integer. If err is
// not equal to nil, then Int64 returns 0, err. Otherwise, Int64 converts the
// reply to an int64 as follows:
//
//  Reply type    Result
//  integer       reply, nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func Int64(reply interface{}, err error) (int64, error) {
	return redigo.Int64(reply, err)
}

// Int64s is a helper that converts an array command reply to a []int64.
// If err is not equal to nil, then Int64s returns nil, err. Nil array
// items are stay nil. Int64s returns an error if an array item is not a
// bulk string or nil.
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func Int64s(reply interface{}, err error) ([]int64, error) {
	return redigo.Int64s(reply, err)
}

// Int64Map is a helper that converts an array of strings (alternating key, value)
// into a map[string]int64. The HGETALL commands return replies in this format.
// Requires an even number of values in result.
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func Int64Map(result interface{}, err error) (map[string]int64, error) {
	return redigo.Int64Map(result, err)
}

/*************************************************************
 *************************  Uint64  **************************
 *************************************************************/

// Uint64 is a helper that converts a command reply to 64 bit unsigned integer.
// If err is not equal to nil, then Uint64 returns 0, err. Otherwise, Uint64 converts the
// reply to an uint64 as follows:
//
//  Reply type    Result
//  +integer      reply, nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func Uint64(reply interface{}, err error) (uint64, error) {
	return redigo.Uint64(reply, err)
}

// Uint64s is a helper that converts an array command reply to a []uint64.
// If err is not equal to nil, then Uint64s returns nil, err. Nil array
// items are stay nil. Uint64s returns an error if an array item is not a
// bulk string or nil.
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func Uint64s(reply interface{}, err error) ([]uint64, error) {
	return redigo.Uint64s(reply, err)
}

// Uint64Map is a helper that converts an array of strings (alternating key, value)
// into a map[string]uint64. The HGETALL commands return replies in this format.
// Requires an even number of values in result.
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func Uint64Map(result interface{}, err error) (map[string]uint64, error) {
	return redigo.Uint64Map(result, err)
}

/*************************************************************
 *************************  Float64  *************************
 *************************************************************/

// Float64 is a helper that converts a command reply to 64 bit float. If err is
// not equal to nil, then Float64 returns 0, err. Otherwise, Float64 converts
// the reply to a float64 as follows:
//
//  Reply type    Result
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func Float64(reply interface{}, err error) (float64, error) {
	return redigo.Float64(reply, err)
}

// Float64s is a helper that converts an array command reply to a []float64. If
// err is not equal to nil, then Float64s returns nil, err. Nil array items are
// converted to 0 in the output slice. Floats64 returns an error if an array
// item is not a bulk string or nil.
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func Float64s(reply interface{}, err error) ([]float64, error) {
	return redigo.Float64s(reply, err)
}

/*************************************************************
 *************************  String  **************************
 *************************************************************/

// String is a helper that converts a command reply to a string. If err is not
// equal to nil, then String returns "", err. Otherwise String converts the
// reply to a string as follows:
//
//  Reply type      Result
//  bulk string     string(reply), nil
//  simple string   reply, nil
//  nil             "",  ErrNil
//  other           "",  error
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func String(reply interface{}, err error) (string, error) {
	return redigo.String(reply, err)
}

// Strings is a helper that converts an array command reply to a []string. If
// err is not equal to nil, then Strings returns nil, err. Nil array items are
// converted to "" in the output slice. Strings returns an error if an array
// item is not a bulk string or nil.
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func Strings(reply interface{}, err error) ([]string, error) {
	return redigo.Strings(reply, err)
}

// StringMap is a helper that converts an array of strings (alternating key, value)
// into a map[string]string. The HGETALL and CONFIG GET commands return replies in this format.
// Requires an even number of values in result.
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func StringMap(result interface{}, err error) (map[string]string, error) {
	return redigo.StringMap(result, err)
}

/*************************************************************
 *************************  Bytes  ***************************
 *************************************************************/

// Bytes is a helper that converts a command reply to a slice of bytes. If err
// is not equal to nil, then Bytes returns nil, err. Otherwise Bytes converts
// the reply to a slice of bytes as follows:
//
//  Reply type      Result
//  bulk string     reply, nil
//  simple string   []byte(reply), nil
//  nil             nil, ErrNil
//  other           nil, error
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func Bytes(reply interface{}, err error) ([]byte, error) {
	return redigo.Bytes(reply, err)
}

// ByteSlices is a helper that converts an array command reply to a [][]byte.
// If err is not equal to nil, then ByteSlices returns nil, err. Nil array
// items are stay nil. ByteSlices returns an error if an array item is not a
// bulk string or nil.
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func ByteSlices(reply interface{}, err error) ([][]byte, error) {
	return redigo.ByteSlices(reply, err)
}

/*************************************************************
 *************************  Others  **************************
 *************************************************************/

// Bool is a helper that converts a command reply to a boolean. If err is not
// equal to nil, then Bool returns false, err. Otherwise Bool converts the
// reply to boolean as follows:
//
//  Reply type      Result
//  integer         value != 0, nil
//  bulk string     strconv.ParseBool(reply)
//  nil             false, ErrNil
//  other           false, error
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func Bool(reply interface{}, err error) (bool, error) {
	return redigo.Bool(reply, err)
}

// Values is a helper that converts an array command reply to a []interface{}.
// If err is not equal to nil, then Values returns nil, err. Otherwise, Values
// converts the reply as follows:
//
//  Reply type      Result
//  array           reply, nil
//  nil             nil, ErrNil
//  other           nil, error
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func Values(reply interface{}, err error) ([]interface{}, error) {
	return redigo.Values(reply, err)
}

// Positions is a helper that converts an array of positions (lat, long)
// into a [][2]float64. The GEOPOS command returns replies in this format.
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func Positions(result interface{}, err error) ([]*[2]float64, error) {
	return redigo.Positions(result, err)
}

// SlowLogs is a helper that parse the SLOWLOG GET command output and
// return the array of SlowLog
//
// See: https://github.com/gomodule/redigo/blob/master/redis/reply.go
func SlowLogs(result interface{}, err error) ([]redigo.SlowLog, error) {
	return redigo.SlowLogs(result, err)
}
