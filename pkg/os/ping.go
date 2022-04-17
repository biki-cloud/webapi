package os

import (
	"fmt"
	"strings"
)

func Ping(ip string) (bool, error) {
	stdout, _, err := SimpleExec("ping -c 1 -W 1")
	if err != nil {
		return false, fmt.Errorf("Ping: %v \n", err)
	}
	return strings.Contains(stdout, "1 packets received"), nil
}
