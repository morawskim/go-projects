package internal

import (
	"github.com/pkg/errors"
	"noip-client/internal/config"
	"noip-client/internal/iphelper"
	"noip-client/internal/noip"
)

var ErrTheSameIpAddr = errors.New("current assigned IP address is the same as public address")

func UpdateNoIpDnsRecord(noIpConfig *config.NoIpConfig, ipHelper iphelper.IpHelper) error {
	currentAssignedIp, err := ipHelper.GetCurrentAssignedIp(noIpConfig.Hostname)
	if err != nil {
		return errors.Wrapf(err, "unable to get current assigned IP address for host %v", noIpConfig.Hostname)
	}

	myCurrentPublicIp, err := ipHelper.GetCurrentPublicIpAddress()
	if err != nil {
		return errors.Wrapf(err, "unable to get current public IP address")
	}

	if !currentAssignedIp.Equal(myCurrentPublicIp) {
		noipApiClient := noip.NewApiClient(noIpConfig)
		updateApiErr := noipApiClient.UpdateAssignedIp(myCurrentPublicIp)
		if updateApiErr != nil {
			return errors.Wrapf(updateApiErr, "unable to update assigned noip address for host %v using username %v", noIpConfig.Hostname, noIpConfig.Username)
		}
	} else {
		return ErrTheSameIpAddr
	}

	return nil
}
