package smtp

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type EmailClient struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewEmailClient(host string, port int, username, password, from string) *EmailClient {
	return &EmailClient{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (c *EmailClient) SendEmail(to []string, subject, body string) error {
	// Set up authentication information
	auth := smtp.PlainAuth("", c.username, c.password, c.host)

	// Prepare message
	message := fmt.Sprintf("From: %s\r\n", c.from) +
		fmt.Sprintf("To: %s\r\n", to[0]) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"\r\n" +
		body

	// Connect to the server
	conn, err := smtp.Dial(fmt.Sprintf("%s:%d", c.host, c.port))
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Start TLS if available
	if ok, _ := conn.Extension("STARTTLS"); ok {
		config := &tls.Config{
			ServerName: c.host,
		}
		if err = conn.StartTLS(config); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	// Authenticate
	if err = conn.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set the sender and recipient
	if err = conn.Mail(c.from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	for _, recipient := range to {
		if err = conn.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}
	}

	// Send the email body
	w, err := conn.Data()
	if err != nil {
		return fmt.Errorf("failed to prepare data: %w", err)
	}
	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return conn.Quit()
}

func (c *EmailClient) SendHTMLEmail(to []string, subject, htmlBody string) error {
	// Similar to SendEmail but with HTML content type
	headers := fmt.Sprintf("From: %s\r\n", c.from) +
		fmt.Sprintf("To: %s\r\n", to[0]) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n"

	message := headers + htmlBody

	return c.sendRawEmail(to, message)
}

func (c *EmailClient) sendRawEmail(to []string, message string) error {
	auth := smtp.PlainAuth("", c.username, c.password, c.host)
	return smtp.SendMail(
		fmt.Sprintf("%s:%d", c.host, c.port),
		auth,
		c.from,
		to,
		[]byte(message),
	)
}