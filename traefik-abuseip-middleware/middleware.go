package traefik_abuseip_middleware

import (
	"context"
	"net/http"
	"os"
)

// Config the plugin configuration.
type Config struct {
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// AbuseIp a plugin.
type AbuseIp struct {
	next http.Handler
	name string
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &AbuseIp{
		next: next,
		name: name,
	}, nil
}

func (e *AbuseIp) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	os.Stdout.WriteString("my first traefik middleware...\n")
	e.next.ServeHTTP(rw, req)
}
