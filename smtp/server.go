package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/containrrr/shoutrrr"
	"github.com/emersion/go-smtp"
)

var notificationURL string

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

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}),
	)
	slog.SetDefault(logger)

	flag.StringVar(&notificationURL, "notification-url", "", "Notification url")
	flag.Parse()

	if notificationURL == "" {
		slog.Default().Error("Notification url is required")
		os.Exit(1)
	}

	server := smtp.NewServer(&Backend{})
	defer server.Close()
	server.Addr = "localhost:1025"
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
