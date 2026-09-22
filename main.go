package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/smtp"
	"os"
	"strings"
)

func main() {
	server := flag.String("server", "localhost:25", "SMTP server address")
	from := flag.String("from", "sender@example.com", "Envelope sender")
	to := flag.String("to", "recipient@example.com", "Envelope recipient(s), comma-separated")
	subject := flag.String("subject", "SMTP test", "Message subject")
	body := flag.String("body", "This is an SMTP test message.", "Message body")
	username := flag.String("username", "", "SMTP username")
	password := flag.String("password", "", "SMTP password")
	startTLS := flag.Bool("starttls", false, "Use STARTTLS")
	insecure := flag.Bool("insecure", false, "Skip TLS certificate verification")
	flag.Parse()

	recipients := splitAddresses(*to)
	if len(recipients) == 0 {
		fmt.Fprintln(os.Stderr, "at least one recipient is required")
		os.Exit(2)
	}

	host, _, err := splitHostPort(*server)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	var auth smtp.Auth
	if *username != "" {
		auth = smtp.PlainAuth("", *username, *password, host)
	}

	message := buildMessage(recipients, *subject, *body)
	if err := sendMail(*server, host, *from, recipients, message, auth, *startTLS, *insecure); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func buildMessage(recipients []string, subject, body string) []byte {
	return []byte("To: " + strings.Join(recipients, ", ") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")
}

func sendMail(address, host, from string, recipients []string, message []byte, auth smtp.Auth, startTLS, insecure bool) error {
	if !startTLS {
		return smtp.SendMail(address, auth, from, recipients, message)
	}

	conn, err := smtp.Dial(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.StartTLS(&tls.Config{ServerName: host, InsecureSkipVerify: insecure}); err != nil {
		return err
	}
	if auth != nil {
		if err := conn.Auth(auth); err != nil {
			return err
		}
	}
	if err := conn.Mail(from); err != nil {
		return err
	}
	for _, recipient := range recipients {
		if err := conn.Rcpt(recipient); err != nil {
			return err
		}
	}
	writer, err := conn.Data()
	if err != nil {
		return err
	}
	if _, err := writer.Write(message); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return conn.Quit()
}

func splitAddresses(value string) []string {
	parts := strings.Split(value, ",")
	addresses := make([]string, 0, len(parts))
	for _, part := range parts {
		if address := strings.TrimSpace(part); address != "" {
			addresses = append(addresses, address)
		}
	}
	return addresses
}

func splitHostPort(address string) (string, string, error) {
	host, port, ok := strings.Cut(address, ":")
	if !ok || host == "" || port == "" {
		return "", "", fmt.Errorf("invalid SMTP server address %q; expected host:port", address)
	}
	return host, port, nil
}

var _ io.Writer
