package application

import "github.com/obnahsgnaw/sockethandler"

type Option func()

func Rps(rps *sockethandler.ManagedRpc) Option {
	return func() {
		rpsIns = rps
	}
}

func IpPort(ip string, port int) Option {
	return func() {
		rpsIp = ip
		rpsPort = port
	}
}
