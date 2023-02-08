package iphelper

import (
	"fmt"
	"io"
	"net"
	"net/http"
)

var ipUrls = [2]string{
	"http://ip1.dynupdate.no-ip.com",
	"http://ip2.dynupdate.no-ip.com",
}

type IpHelperImpl struct {
}

func (i *IpHelperImpl) GetCurrentAssignedIp(hostname string) (net.IP, error) {
	addr, err := net.LookupIP(hostname)

	if err != nil {
		return nil, fmt.Errorf("cannot get IP address for host %v, reason %v", hostname, err)
	}

	return addr[0], nil
}

func (i *IpHelperImpl) GetCurrentPublicIpAddress() (net.IP, error) {
	res, err := http.Get(ipUrls[0])

	if err != nil {
		return nil, fmt.Errorf("unable to get current public IP address, %v", err)
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode < 200 || res.StatusCode > 299 {

		return nil, fmt.Errorf("cant get current public IP address. Got %d response status code. Body: %s", res.StatusCode, body)
	}

	return net.ParseIP(string(body)), nil
}

func NewIpHelper() IpHelper {
	return &IpHelperImpl{}
}
