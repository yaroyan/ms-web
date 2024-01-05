package usecase

import (
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/yaroyan/ms/authn/domain/model/authentication"
)

type AuthenticateUserUsecase struct {
	repository authentication.Repository
}

type RequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ユーザを認証する
func (a *AuthenticateUserUsecase) Authenticate(p RequestPayload) (string, error) {
	u, err := a.repository.GetByEmail(p.Email)
	if err != nil {
		return "", err
	}

	isValid, err := u.Password.Match(p.Password)
	if err != nil || !isValid {
		return "", err
	}

	t := a.toJWT(u)

	st, err := a.toSignedJWTToken(t)

	if err != nil {
		return "", err
	}

	return st, nil
}

// JWTを発行する
func (*AuthenticateUserUsecase) toJWT(u *authentication.User) *jwt.Token {
	t := jwt.NewWithClaims(jwt.SigningMethodES512, jwt.MapClaims{
		"email": u.Email.Email,
	})
	return t
}

// JWTを署名する
func (*AuthenticateUserUsecase) toSignedJWTToken(t *jwt.Token) (string, error) {
	// JWT_SECRET=`openssl rand -base64 172 | tr -d '\n'`
	if st, err := t.SignedString(os.Getenv("JWT_SECRET")); err != nil {
		return "", err
	} else {
		return st, nil
	}
}
