# smtp-testing-go

A lightweight SMTP testing tool written in Go, inspired by [Swaks](https://github.com/jetmore/swaks). It helps developers and operators test SMTP servers, authentication, TLS, recipients, and message delivery with minimal dependencies.

## Features

- Plain SMTP, STARTTLS, and implicit TLS
- SMTP authentication with PLAIN and CRAM-MD5
- Multiple recipients
- Custom headers
- Configurable timeout
- Cross-platform binaries released with GoReleaser

## Installation

Download a binary from the [latest release](https://github.com/vignemail1/smtp-testing-go/releases), or build from source:

```sh
go build -o smtp-testing-go .
```

## Usage

```sh
./smtp-testing-go \
  -server smtp.example.com:587 \
  -from sender@example.com \
  -to recipient@example.com \
  -subject "SMTP test" \
  -body "This is a test message" \
  -starttls
```

For authentication, provide `-username` and `-password`. If `-password` is omitted, the program uses the `SMTP_PASSWORD` environment variable when it is set:

```sh
SMTP_PASSWORD='secret' ./smtp-testing-go \
  -server smtp.example.com:587 \
  -from sender@example.com \
  -to recipient@example.com \
  -username user@example.com \
  -starttls
```

Use `-help` to display all available options.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
