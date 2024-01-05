package authentication

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// パスワード
type Password struct {
	Password []byte
}

// 適格なパスワードであるか検証する
func (*Password) isValidPassword(pw string) error {
	return nil
}

// パスワードを検証します
func (p *Password) Match(pw string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(pw))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

// パスワードの暗号化
func (*Password) encrypt(pw string) ([]byte, error) {
	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost); err != nil {
		return nil, err
	} else {
		return hashedPassword, nil
	}
}

// パスワードの生成
func NewPassword(pw string) (*Password, error) {
	p := &Password{}
	if err := p.isValidPassword(pw); err != nil {
		return nil, err
	}
	if encrypted, err := p.encrypt(pw); err != nil {
		return nil, err
	} else {
		p.Password = encrypted
	}
	return p, nil
}
