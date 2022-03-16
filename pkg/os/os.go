package os

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func ListEnvToSlice(listEnv string) []string {
	return strings.Split(listEnv, ",")
}

func SetEnvIfNotExists(k, v string) error {
	_, ok := os.LookupEnv(k)
	if !ok {
		err := os.Setenv(k, v)
		if err != nil {
			return fmt.Errorf("SetEnvIfNotExists: %v ", err)
		}
	}
	return nil
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
