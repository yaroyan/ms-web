package authentication

import (
	"net/mail"
)

type Email struct {
	Email string `json:"email"`
}

// 適格なEmailであるか検証します
func (*Email) isValidEmail(e string) error {
	_, err := mail.ParseAddress(e)
	return err
}

// メールを生成します
func NewEmail(email string) (*Email, error) {
	p := &Email{}
	if err := p.isValidEmail(email); err != nil {
		return nil, err
	}
	p.Email = email
	return p, nil
}
