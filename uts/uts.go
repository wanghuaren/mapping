package uts

import (
	"crypto/md5"
	"encoding/hex"
	"net"
	"strings"
)

var IsDebug = true

func Debug(debug bool) bool {
	IsDebug = debug
	Log("PROJECT DEBUG MODE ", IsDebug)
	return IsDebug
}

// a_a -> aA
func UnderScoreCase2CamelCase(nameStr string) string {
	arr := strings.Split(nameStr, "_")
	result := arr[0]
	for i := 1; i < len(arr); i++ {
		result += strings.ToUpper(arr[i][:1]) + arr[i][1:]
	}
	return result
}

var big = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// aA -> a_a
func CamelCase2UnderScoreCase(nameStr string) string {
	// result := strings.ToLower(nameStr[:1])
	result := nameStr[:1]
	for i := 1; i < len(nameStr); i++ {
		w := nameStr[i : i+1]
		if strings.Contains(big, w) {
			result += "_" + strings.ToLower(w)
		} else {
			result += w
		}
	}
	return result
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if !ChkErr(err) {
		for _, address := range addrs {
			// 检查ip地址判断是否回环地址
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}
	return ""
}

func GetMd5(str string) string {
	d := []byte(str)
	m := md5.New()
	m.Write(d)
	_token := hex.EncodeToString(m.Sum(nil))
	return _token
}
