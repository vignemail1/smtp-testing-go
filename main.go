package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type list []string

func (l *list) String() string { return strings.Join(*l, ",") }
func (l *list) Set(v string) error { *l = append(*l, v); return nil }

type cfg struct {
	server, from, to, subject, body, helo, user, pass, auth string
	headers list
	starttls, tls, insecure, verbose bool
	timeout time.Duration
}

func main() {
	var c cfg
	flag.StringVar(&c.server, "server", "", "SMTP server host:port")
	flag.StringVar(&c.from, "from", "", "envelope sender (MAIL FROM)")
	flag.StringVar(&c.to, "to", "", "comma-separated recipients")
	flag.StringVar(&c.subject, "subject", "SMTP test", "message subject")
	flag.StringVar(&c.body, "body", "SMTP test message", "message body")
	flag.Var(&c.headers, "header", "custom message header (Name: value); repeatable")
	flag.StringVar(&c.helo, "helo", "localhost", "hostname used in EHLO/HELO")
	flag.BoolVar(&c.starttls, "starttls", false, "enable STARTTLS")
	flag.BoolVar(&c.tls, "tls", false, "use implicit TLS")
	flag.BoolVar(&c.insecure, "insecure", false, "skip TLS certificate verification")
	flag.StringVar(&c.user, "username", "", "SMTP username")
	flag.StringVar(&c.pass, "password", "", "SMTP password (falls back to SMTP_PASSWORD)")
	flag.StringVar(&c.auth, "auth", "", "authentication mechanism: PLAIN or CRAM-MD5")
	flag.DurationVar(&c.timeout, "timeout", 30*time.Second, "connection timeout")
	flag.BoolVar(&c.verbose, "verbose", false, "print connection information")
	flag.Parse()

	if c.server == "" || c.from == "" || c.to == "" {
		flag.Usage()
		os.Exit(2)
	}
	if c.pass == "" {
		c.pass = os.Getenv("SMTP_PASSWORD")
	}
	recipients := addresses(c.to)
	if len(recipients) == 0 {
		fmt.Fprintln(os.Stderr, "no recipients provided")
		os.Exit(2)
	}
	if err := send(c, recipients); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println("message sent")
}

func addresses(s string) []string {
	var result []string
	for _, address := range strings.Split(s, ",") {
		if address = strings.TrimSpace(address); address != "" {
			result = append(result, address)
		}
	}
	return result
}

func send(c cfg, recipients []string) error {
	host, _, err := net.SplitHostPort(c.server)
	if err != nil {
		host = c.server
	}

	var client *smtp.Client
	if c.tls {
		conn, err := tls.DialWithDialer(&net.Dialer{Timeout: c.timeout}, "tcp", c.server, &tls.Config{ServerName: host, InsecureSkipVerify: c.insecure})
		if err != nil { return err }
		client, err = smtp.NewClient(conn, host)
	} else {
		conn, err := (&net.Dialer{Timeout: c.timeout}).Dial("tcp", c.server)
		if err != nil { return err }
		client, err = smtp.NewClient(conn, host)
	}
	if err != nil { return err }
	defer client.Close()
	if c.verbose { fmt.Println("connected to", c.server) }
	if err = client.Hello(c.helo); err != nil { return err }
	if c.starttls && !c.tls {
		ok, _ := client.Extension("STARTTLS")
		if !ok { return fmt.Errorf("STARTTLS is not advertised") }
		if err = client.StartTLS(&tls.Config{ServerName: host, InsecureSkipVerify: c.insecure}); err != nil { return err }
	}
	if c.user != "" {
		var auth smtp.Auth
		switch strings.ToUpper(c.auth) {
		case "", "PLAIN": auth = smtp.PlainAuth("", c.user, c.pass, host)
		case "CRAM-MD5": auth = smtp.CRAMMD5Auth(c.user, c.pass)
		default: return fmt.Errorf("unsupported authentication mechanism: %s", c.auth)
		}
		if err = client.Auth(auth); err != nil { return err }
	}
	if err = client.Mail(c.from); err != nil { return err }
	for _, recipient := range recipients {
		if err = client.Rcpt(recipient); err != nil { return err }
	}
	writer, err := client.Data()
	if err != nil { return err }
	if _, err = writer.Write([]byte(message(c, recipients))); err != nil { writer.Close(); return err }
	if err = writer.Close(); err != nil { return err }
	return client.Quit()
}

func message(c cfg, recipients []string) string {
	headers := []string{
		"Date: " + time.Now().Format(time.RFC1123Z),
		"From: " + c.from,
		"To: " + strings.Join(recipients, ", "),
		"Subject: " + c.subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
	}
	for _, header := range c.headers {
		if strings.Contains(header, ":") { headers = append(headers, header) }
	}
	return strings.Join(headers, "\r\n") + "\r\n\r\n" + strings.ReplaceAll(c.body, "\n", "\r\n") + "\r\n"
}
