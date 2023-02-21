package noip

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net"
	"net/http"
	"noip-client/internal/config"
	"testing"
)

type ClientMock struct {
	response *http.Response
}

func (c *ClientMock) setResponse(response *http.Response) {
	c.response = response
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	return c.response, nil
}

func TestValidateNoIpUpdateResponseSuccess(t *testing.T) {

	status := []string{
		"nochg",
		"good",
		" good ",
	}

	for _, val := range status {
		ok, err := validateNoIpUpdateResponse(val)

		assert.NoError(t, err)
		assert.True(t, ok)
	}
}

func TestValidateNoIpUpdateResponseError(t *testing.T) {
	status := []string{
		"foo",
		"",
		"nohost",
		"badauth",
		"badagent",
		"!donato",
		"abuse",
		"911",
	}

	for _, val := range status {
		ok, err := validateNoIpUpdateResponse(val)

		assert.Error(t, err)
		assert.False(t, ok)
		assert.NotEmpty(t, err.Error())
	}
}

func TestMockHttp(t *testing.T) {
	mock := &ClientMock{}
	mock.setResponse(&http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(`  good`))),
	})

	client := &ApiClientHttp{
		httpClient: mock,
		config: &config.NoIpConfig{
			Username: "foo",
			Password: "pass",
			Hostname: "example.com",
		},
	}

	err := client.UpdateAssignedIp(net.ParseIP("192.168.1.1"))

	assert.NoError(t, err)
}
