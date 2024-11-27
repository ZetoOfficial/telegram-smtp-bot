package usecase

import (
	"fmt"

	"github.com/ZetoOfficial/telegram-smtp-bot/internal/adapter/smtp"
	"github.com/ZetoOfficial/telegram-smtp-bot/internal/domain"
)

// BotUseCase описывает бизнес-логику бота
type BotUseCase struct {
	SMTPClient smtp.SMTPClient
}

// NewBotUseCase создаёт новый экземпляр BotUseCase
func NewBotUseCase(smtpClient smtp.SMTPClient) *BotUseCase {
	return &BotUseCase{
		SMTPClient: smtpClient,
	}
}

// ValidateEmail валидирует email адрес
func (b *BotUseCase) ValidateEmail(address string) (*domain.Email, error) {
	email, err := domain.NewEmail(address)
	if err != nil {
		return nil, fmt.Errorf("email validation: %w", err)
	}
	return email, nil
}

// CreateMessage валидирует и создаёт сообщение
func (b *BotUseCase) CreateMessage(text string) (*domain.Message, error) {
	message, err := domain.NewMessage(text)
	if err != nil {
		return nil, fmt.Errorf("message validation: %w", err)
	}
	return message, nil
}

// SendEmail отправляет email с заданным сообщением
func (b *BotUseCase) SendEmail(to string, message string) error {
	err := b.SMTPClient.SendMail(to, message)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil
}
