package nsdyndns

import "net"

type NameService interface {
	SetAddress(domain, host string, address net.IP) error
}
