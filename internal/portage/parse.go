package portage

import (
	"fmt"
	"strconv"
	"strings"
)

func key(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

func splitKey(k string) (string, int) {
	i := strings.LastIndex(k, ":")
	if i < 0 {
		return k, 0
	}
	host := k[:i]
	port, _ := strconv.Atoi(k[i+1:])
	return host, port
}
