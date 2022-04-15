package common

import (
	"net"
	"strings"
)

func Connect(protoAddr string) (net.Conn, error) {
	parts := strings.SplitN(protoAddr, "://", 2)
	proto, address := parts[0], parts[1]
	conn, err := net.Dial(proto, address)
	return conn, err
}
