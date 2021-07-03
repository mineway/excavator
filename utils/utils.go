package utils

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func InArrayInt(arr []int, find int) bool {
	for _, i := range arr {
		if i == find {
			return true
		}
	}
	return false
}

func GetAvailablePort() (string, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}
	return  strconv.Itoa(listener.Addr().(*net.TCPAddr).Port), nil
}

func IsAvailablePort(port string) bool {
	conn, _ := net.DialTimeout("tcp", fmt.Sprintf(":%s", port), 3*time.Second)
	if conn != nil {
		_ = conn.Close()
		return false
	}
	return true
}