package noip

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"noip-client/internal/config"
	"strings"
	"time"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ApiClientHttp struct {
	httpClient HttpClient
	config     *config.NoIpConfig
}

func NewApiClient(config *config.NoIpConfig) *ApiClientHttp {
	return &ApiClientHttp{
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
		config: config,
	}
}

func validateNoIpUpdateResponse(response string) (bool, error) {
	status := strings.TrimSpace(response)

	if strings.HasPrefix(status, "nochg") || strings.HasPrefix(status, "good") {
		return true, nil
	}

	apiError := ApiError(status)

	return false, errors.New(apiError.String())
}

func (n *ApiClientHttp) buildUpdateRequest(newIp net.IP) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, "https://dynupdate.no-ip.com/nic/update", nil)
	if err != nil {
		return nil, fmt.Errorf("Got error %s", err.Error())
	}

	req.SetBasicAuth(n.config.Username, n.config.Password)
	q := req.URL.Query()
	q.Add("hostname", n.config.Hostname)
	q.Add("myip", newIp.String())
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func (n *ApiClientHttp) UpdateAssignedIp(newIp net.IP) error {
	req, err := n.buildUpdateRequest(newIp)

	response, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("cannot build HTTP request %v", err.Error())
	}

	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	ok, err := validateNoIpUpdateResponse(string(content))

	if ok {
		return nil
	}

	return err
}
