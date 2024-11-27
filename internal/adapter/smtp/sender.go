package smtp

import (
	"fmt"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

type SMTPClient interface {
	SendMail(to string, message string) error
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type SMTPSender struct {
	Config SMTPConfig
}

func NewSMTPSender(config SMTPConfig) *SMTPSender {
	return &SMTPSender{Config: config}
}

func (s *SMTPSender) SendMail(to string, message string) error {
	server := mail.NewSMTPClient()
	server.Host = s.Config.Host
	server.Port = s.Config.Port
	server.Username = s.Config.Username
	server.Password = s.Config.Password
	if s.Config.Port == 587 {
		server.Encryption = mail.EncryptionSTARTTLS
	} else if s.Config.Port == 465 {
		server.Encryption = mail.EncryptionSSLTLS
	}

	// Устанавливаем тайм-ауты
	server.ConnectTimeout = 5 * time.Second
	server.SendTimeout = 5 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return fmt.Errorf("не удалось подключиться к SMTP-серверу: %w", err)
	}

	email := mail.NewMSG()
	email.SetFrom(s.Config.Username).
		AddTo(to).
		SetSubject("Telegram Bot Message").
		SetBody(mail.TextPlain, message)

	err = email.Send(smtpClient)
	if err != nil {
		return fmt.Errorf("не удалось отправить письмо: %w", err)
	}

	return nil
}
