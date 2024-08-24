package traefik_abuseip_middleware

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Config the plugin configuration.
type Config struct {
	AbuseIpFile string `json:"abuse_ip_file"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		AbuseIpFile: "https://github.com/borestad/blocklist-abuseipdb/raw/main/abuseipdb-s100-14d.ipv4",
	}
}

// AbuseIp a plugin.
type AbuseIp struct {
	next     http.Handler
	name     string
	abuseIps map[string]bool
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	ips, err := getAbuseIp(config.AbuseIpFile)

	if err != nil {
		return nil, err
	}

	return &AbuseIp{
		next:     next,
		name:     name,
		abuseIps: ips,
	}, nil
}

func (e *AbuseIp) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	os.Stdout.WriteString("my first traefik middleware...\n")
	//what about X-Forwarded-For Header?
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		os.Stderr.WriteString("error parsing remote address: " + req.RemoteAddr + "\n")
		e.next.ServeHTTP(rw, req)
		return
	}

	if _, exists := e.abuseIps[host]; exists {
		http.Error(rw, fmt.Sprintf(`Your IP "%s" is blocked`, host), http.StatusForbidden)
	} else {
		e.next.ServeHTTP(rw, req)
	}
}

func getAbuseIp(abuseIpDbUrl string) (map[string]bool, error) {
	get, err := http.Get(abuseIpDbUrl)

	if err != nil {
		return nil, fmt.Errorf("unable to download abuse db: %w", err)
	}

	defer get.Body.Close()

	scanner := bufio.NewScanner(get.Body)
	var ipSlice map[string]bool

	for scanner.Scan() {
		text := scanner.Text()

		if strings.HasPrefix(text, "# Number of ips:") {
			b := strings.TrimSpace(strings.TrimPrefix(text, "# Number of ips:"))
			nbOfIpsInFile, err := strconv.Atoi(b)

			if err != nil {
				return nil, fmt.Errorf("unable to convert nbOfIpsInFile to int: %w", err)
			}

			ipSlice = make(map[string]bool, nbOfIpsInFile)
			continue
		}

		if strings.HasPrefix(text, "#") {
			continue
		}

		ipSlice[text] = true
	}

	if scanner.Err() != nil {
		return nil, fmt.Errorf("db parse error: %w", scanner.Err())
	}

	return ipSlice, nil
}
