package email

import (
	"akastra-access/internal/app/config"
	"crypto/tls"
	"net/smtp"
)

func SendEmail(to, subject, body string) error {
	cfg := config.Load().Email

	from := cfg.EmailUsername
	password := cfg.EmailPassword
	host := cfg.SmtpHost
	port := cfg.SmtpPort

	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
		body + "\r\n")

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Connect with TLS
	conn, err := tls.Dial("tcp", host+":"+port, tlsconfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", from, password, host)

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(from); err != nil {
		return err
	}

	if err = client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}
