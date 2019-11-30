package tools

import (
	"fmt"
	"strings"
)

const networkSplit = "@"

// ParseNetwork ...
func ParseNetwork(str string) (network, addr string, err error) {
	idx := strings.Index(str, networkSplit)
	if idx == -1 {
		err = fmt.Errorf("addr: \"%s\" error, must be network@tcp:port or network@unixsocket", str)
		return
	}
	network = str[:idx]
	addr = str[idx+1:]
	return
}
