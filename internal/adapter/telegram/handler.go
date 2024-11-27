package telegram

import (
	"context"
	"log"
	"sync"

	"github.com/ZetoOfficial/telegram-smtp-bot/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type UserSession struct {
	Step    int
	Email   string
	Message string
	mu      sync.Mutex
}

type TelegramHandler struct {
	Bot      *tgbotapi.BotAPI
	UseCase  *usecase.BotUseCase
	Sessions map[int64]*UserSession
}

func NewTelegramHandler(bot *tgbotapi.BotAPI, useCase *usecase.BotUseCase) *TelegramHandler {
	return &TelegramHandler{
		Bot:      bot,
		UseCase:  useCase,
		Sessions: make(map[int64]*UserSession),
	}
}

func (h *TelegramHandler) HandleUpdates(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := h.Bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Ошибка получения обновлений: %v", err)
	}

	for {
		select {
		case update, ok := <-updates:
			if !ok {
				log.Println("Канал обновлений закрыт.")
				return
			}
			if update.Message == nil {
				continue
			}

			userID := update.Message.Chat.ID
			h.initSession(userID)

			session := h.Sessions[userID]
			session.mu.Lock()

			switch session.Step {
			case 0:
				h.handleStart(update.Message)
			case 1:
				h.handleEmail(update.Message)
			case 2:
				h.handleMessage(update.Message)
			}

			session.mu.Unlock()

		case <-ctx.Done():
			log.Println("Контекст отменён, завершение обработки обновлений.")
			return
		}
	}
}

func (h *TelegramHandler) initSession(userID int64) {
	if _, exists := h.Sessions[userID]; !exists {
		h.Sessions[userID] = &UserSession{Step: 0}
	}
}

func (h *TelegramHandler) handleStart(message *tgbotapi.Message) {
	response := "Здравствуйте! Пожалуйста, введите ваш email адрес."
	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	h.Bot.Send(msg)

	session := h.Sessions[message.Chat.ID]
	session.Step = 1
}

func (h *TelegramHandler) handleEmail(message *tgbotapi.Message) {
	emailInput := message.Text
	email, err := h.UseCase.ValidateEmail(emailInput)
	if err != nil {
		response := "Некорректный email адрес. Пожалуйста, попробуйте снова."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		h.Bot.Send(msg)
		return
	}

	session := h.Sessions[message.Chat.ID]
	session.Email = email.Address
	session.Step = 2

	response := "Введите текст сообщения, которое вы хотите отправить:"
	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	h.Bot.Send(msg)
}

func (h *TelegramHandler) handleMessage(message *tgbotapi.Message) {
	messageText := message.Text
	msgObj, err := h.UseCase.CreateMessage(messageText)
	if err != nil {
		log.Printf("Failed to create message %v", err)
		response := "Сообщение не может быть пустым. Пожалуйста, попробуйте снова."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		h.Bot.Send(msg)
		return
	}

	session := h.Sessions[message.Chat.ID]
	err = h.UseCase.SendEmail(session.Email, msgObj.Text)
	if err != nil {
		log.Printf("Failed to send email %v", err)
		response := "Не удалось отправить email. Пожалуйста, попробуйте позже."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		h.Bot.Send(msg)
		return
	}

	response := "Ваше сообщение успешно отправлено!"
	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	h.Bot.Send(msg)

	// Сбросить сессию
	session.Step = 0
	session.Email = ""
	session.Message = ""
}
