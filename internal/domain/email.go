package domain

import (
	"errors"
	"net/mail"
)

type Email struct {
	Address string
}

func NewEmail(address string) (*Email, error) {
	_, err := mail.ParseAddress(address)
	if err != nil {
		return nil, errors.New("incorrect email")
	}
	return &Email{Address: address}, nil
}
