package domain

import "errors"

type Message struct {
	Text string
}

func NewMessage(text string) (*Message, error) {
	if len(text) == 0 {
		return nil, errors.New("message can't be empty")
	}
	return &Message{Text: text}, nil
}
