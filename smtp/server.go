package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/containrrr/shoutrrr"
	"github.com/emersion/go-smtp"
)

var notificationURL string
var addr string

type Backend struct{}
type Session struct {
	from string
	rcpt []string
}

func (s *Session) Reset() {
}

func (s *Session) Logout() error {
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.rcpt = append(s.rcpt, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if b, err := io.ReadAll(r); err != nil {
		return err
	} else {
		err := shoutrrr.Send(notificationURL, string(b))
		if err != nil {
			slog.Default().Error(fmt.Sprintf("Cannot forward message: %v", err))
			return errors.New("cannot forward message")
		}
	}
	return nil
}

func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{
		rcpt: make([]string, 1),
	}, nil
}

func validateAddr(addr string) error {
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}

	return nil
}

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}),
	)
	slog.SetDefault(logger)

	flag.StringVar(&notificationURL, "notification-url", "", "Notification url")
	flag.StringVar(&addr, "addr", "127.0.0.1:25", "SMTP server listening address in host:port format")
	flag.Parse()

	if notificationURL == "" {
		slog.Default().Error("Notification url is required")
		os.Exit(1)
	}

	if validateAddr(addr) != nil {
		slog.Default().Error("Invalid address")
		os.Exit(1)
	}

	server := smtp.NewServer(&Backend{})
	defer server.Close()
	server.Addr = addr
	server.Domain = "localhost"
	server.WriteTimeout = 5 * time.Second
	server.ReadTimeout = 5 * time.Second
	server.MaxMessageBytes = 1024 * 1024
	server.MaxRecipients = 10
	server.AllowInsecureAuth = true

	logger.Info(fmt.Sprintf("Starting server at %s", server.Addr))
	if err := server.ListenAndServe(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
