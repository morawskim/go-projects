package main

import (
	"flag"
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/exp/slog"
	"net"
	"os"
	"time"
)

type contextKey string

func (c contextKey) String() string {
	return "_fake_ssh_ctx_" + string(c)
}

var (
	contextKeySlog = contextKey("slog")
)

func main() {
	var address string
	flag.StringVar(&address, "address", ":2222", "Address")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout))
	logger.Enabled(nil, slog.LevelInfo)

	// Create a new server instance
	s := &ssh.Server{
		Addr: address,
		// handler will never been called, because we deny all authentication request
		Handler:         handleSession,
		PasswordHandler: handleAuthentication,
		IdleTimeout:     60 * time.Second,
		ConnCallback: func(ctx ssh.Context, conn net.Conn) net.Conn {
			ctx.SetValue(contextKeySlog, logger)
			return conn
		},
	}

	logger.Info(fmt.Sprintf("Starting SSH server on address: %v", s.Addr))
	err := s.ListenAndServe()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func handleSession(s ssh.Session) {
	// this code will never execute because neither user is allowed to connect
	s.Write([]byte("Hello World!"))
	s.Close()
}

func handleAuthentication(ctx ssh.Context, passwd string) bool {
	logger := ctx.Value(contextKeySlog).(*slog.Logger)

	logger.Info(
		"New authentication request",
		slog.String("username", ctx.User()),
		slog.String("password", passwd),
		slog.String("remote", ctx.RemoteAddr().String()),
	)

	// deny all connections
	return false
}
