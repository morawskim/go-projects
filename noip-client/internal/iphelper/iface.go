package iphelper

import "net"

type IpHelper interface {
	GetCurrentAssignedIp(hostname string) (net.IP, error)
	GetCurrentPublicIpAddress() (net.IP, error)
}
