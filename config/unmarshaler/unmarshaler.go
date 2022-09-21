package unmarshaler

import (
	"os"
	"sync"
)

var (
	unmarshalerMap   = map[string]Unmarshaler{}
	unmarshalerMapRW sync.RWMutex
)

// register register unmarshaler
func register(unmarshaler Unmarshaler) {
	unmarshalerMapRW.Lock()
	defer unmarshalerMapRW.Unlock()
	unmarshalerMap[unmarshaler.Name()] = unmarshaler
}

// Get get unmarshaler by name
func Get(name string) Unmarshaler {
	unmarshalerMapRW.RLock()
	defer unmarshalerMapRW.RUnlock()
	return unmarshalerMap[name]
}

// Unmarshaler ...
type Unmarshaler interface {

	// Unmarshal ...
	Unmarshal([]byte, interface{}) error

	// Name unmarshaler name
	Name() string
}

// expandEnv 寻找 ${var} 并替换为环境变量的值，没有则替换为空，不解析 $var
//
// os.ExpandEnv 会同时处理${var}和$var，配置文件中可能包含一些含特殊字符$的配置项，
// 如redisClient、mysqlClient的连接密码。
func expandEnv(s string) string {
	var buf []byte
	i := 0
	for j := 0; j < len(s); j++ {
		if s[j] == '$' && j+2 < len(s) && s[j+1] == '{' {
			if buf == nil {
				buf = make([]byte, 0, 2*len(s))
			}
			buf = append(buf, s[i:j]...)
			name, w := getEnvName(s[j+1:])
			if name == "" && w > 0 {
			} else if name == "" {
				buf = append(buf, s[j])
			} else {
				buf = append(buf, os.Getenv(name)...)
			}
			j += w
			i = j + 1
		}
	}
	if buf == nil {
		return s
	}
	return string(buf) + s[i:]
}

// getEnvName 获取环境变量名，即${var}里面的var内容，返回var内容及其长度
func getEnvName(s string) (string, int) {
	// 匹配右括号 }
	// 输入已经保证第一个字符是{，并且至少两个字符以上
	for i := 1; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '\n' || s[i] == '"' { // "xx${xxx"
			return "", 0 // 遇到上面这些字符认为没有匹配中，保留$
		}
		if s[i] == '}' {
			if i == 1 { // ${}
				return "", 2 // 去掉${}
			}
			return s[1:i], i + 1
		}
	}
	return "", 0 // 没有右括号，保留$
}
