package noip

import "net"

type NoIpApiClient interface {
	UpdateAssignedIp(newIp net.IP) error
}
