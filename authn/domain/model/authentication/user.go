package authentication

import (
	"fmt"
	"time"
)

// ユーザ
type User struct {
	ID        int
	Email     Email
	FirstName string
	LastName  string
	Password  Password
	Active    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ユーザを生成します。
func NewUser(
	id int,
	email string,
	firstName string,
	lastName string,
	password string,
	active int,
	createdAt time.Time,
	updatedAt time.Time,
) (*User, error) {
	u := User{}
	if err := u.setId(id); err != nil {
		return nil, err
	}
	if err := u.setPassword(password); err != nil {
		return nil, err
	}
	if err := u.setEmail(email); err != nil {
		return nil, err
	}
	return &u, nil
}

// パスワードをセットします
func (u *User) setPassword(pw string) error {
	if u == nil {
		return fmt.Errorf("receiver is nil")
	}
	p, err := NewPassword(pw)
	if err != nil {
		return err
	}
	u.Password = *p
	return nil
}

// IDをセットします
func (u *User) setId(i int) error {
	if u == nil {
		return fmt.Errorf("receiver is nil")
	}
	if !u.isValidId(i) {
		return fmt.Errorf("invalid id")
	}
	u.ID = i
	return nil
}

// メールアドレスをセットします
func (u *User) setEmail(e string) error {
	if u == nil {
		return fmt.Errorf("receiver is nil")
	}
	em, err := NewEmail(e)
	if err != nil {
		return err
	}
	u.Email = *em
	return nil
}

// IDの妥当性を検証します。
func (*User) isValidId(i int) bool {
	return true
}

// メールアドレスの妥当性を検証します。
func (*User) isValidEmail(e string) bool {
	return true
}

// 姓の妥当性を検証します。
func (*User) isValidFirstName(f string) bool {
	return true
}

// 名の妥当性を検証します。
func (*User) isValidLastName(l string) bool {
	return true
}

// 活性の妥当性を検証します。
func (*User) isValidActive(e string) bool {
	return true
}

// 登録日の妥当性を検証します。
func (*User) isValidCreatedAt(e string) bool {
	return true
}

// 更新日の妥当性を検証します。
func (*User) isValidUpdatedAt(e string) bool {
	return true
}
