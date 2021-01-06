package telegram

import (
	"net"

	"github.com/pkg/errors"
	"golang.org/x/net/proxy"
)

type directTCP struct{}

// Direct is a direct proxy: one that makes tcp network connections directly.
var DirectTCP = directTCP{}

func (directTCP) Dial(network, addr string) (net.Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "resolving tcp")
	}
	return net.DialTCP(network, nil, tcpAddr)
}

// NewSOCK5ProxyDialer new socks5 proxy
// ("tcp",
// "106.53.131.105:65431",
//  &proxy.Auth{
// 	User:     "Kassulke8264",
// 	Password: "wFKo1z8xOr",
// })
func NewSOCK5ProxyDialer(network, addr string, auth *proxy.Auth) proxy.Dialer {
	socks5, _ := proxy.SOCKS5("tcp", "106.53.131.105:65431", &proxy.Auth{
		User:     "Kassulke8264",
		Password: "wFKo1z8xOr",
	}, DirectTCP)
	return socks5
}
