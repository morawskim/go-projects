package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

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
		chunks := splitStringIntoChunks(string(b))
		for _, chunk := range chunks {
			err := shoutrrr.Send(notificationURL, chunk)
			if err != nil {
				slog.Default().Error(fmt.Sprintf("Cannot forward message: %v", err))
				return errors.New("cannot forward message")
			}
		}
	}
	return nil
}

func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{
		rcpt: make([]string, 1),
	}, nil
}

const (
	maxChunkBodySize = 3800
	longLineSize     = 300
	maxSubjectLength = 190
	maxChunkSize     = maxChunkBodySize + maxSubjectLength
)

func extractSubject(raw string, maxSubjectLength int) string {
	msg, err := mail.ReadMessage(strings.NewReader(raw))
	if err != nil {
		return ""
	}
	subject := msg.Header.Get("Subject")

	return subject[0:min(len(subject), maxSubjectLength)]
}

func splitStringIntoChunks(s string) []string {
	var chunks []string
	isFirst := true
	subject := extractSubject(s, maxSubjectLength)

	if len(subject) > 0 {
		subject = "Subject: " + subject + "\n\n"
	}

	prependSubjectFn := func(isFirst bool, subject, txt string) string {
		if !isFirst && len(subject) > 0 {
			return subject + txt
		}

		return txt
	}

	for len(s) > maxChunkBodySize {
		cut, skip := findSplitPoint(s)
		txt := s[:cut]
		chunks = append(chunks, prependSubjectFn(isFirst, subject, txt))
		s = s[cut+skip:]
		isFirst = false
	}

	if len(s) > 0 {
		chunks = append(chunks, prependSubjectFn(isFirst, subject, s))
	}

	return chunks
}

func findSplitPoint(s string) (cut, skip int) {
	window := s[:maxChunkBodySize+1]

	nl := strings.LastIndexByte(window, '\n')
	if nl > 0 {
		lineStart := nl + 1
		lineEnd := len(s)
		if i := strings.IndexByte(s[lineStart:], '\n'); i >= 0 {
			lineEnd = lineStart + i
		}
		if lineEnd-lineStart <= longLineSize {
			return nl, 1
		}
	}

	from := 0
	if nl >= 0 {
		from = nl + 1
	}
	if sp := strings.LastIndexByte(window[from:], ' '); sp >= 0 && from+sp > 0 {
		return from + sp, 1
	}
	if nl > 0 {
		return nl, 1
	}

	i := maxChunkBodySize
	for i > maxChunkBodySize-utf8.UTFMax && !utf8.RuneStart(s[i]) {
		i--
	}
	return i, 0
}

func validateAddr(addr string) error {
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}

	return nil
}

func readCredential(name string) (string, error) {
	if v := os.Getenv(name); v != "" {
		return v, nil
	}

	if f := os.Getenv(name + "_FILE"); f != "" {
		b, err := os.ReadFile(f)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(b)), nil
	}

	if dir := os.Getenv("CREDENTIALS_DIRECTORY"); dir != "" {
		p := filepath.Join(dir, strings.ToLower(name))
		if b, err := os.ReadFile(p); err == nil {
			return strings.TrimSpace(string(b)), nil
		}
	}

	return "", fmt.Errorf("credential %s not found", name)
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
		val, err := readCredential("NOTIFICATION_URL")
		if err == nil {
			notificationURL = val
		}
	}

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
