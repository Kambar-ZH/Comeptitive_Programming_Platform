package middleware

import (
	"context"
	"fmt"
	"site/internal/datastruct"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type ctxKey int8

const (
	SessionName        = "gopherschool"
	CtxKeyUser  ctxKey = iota
	CtxKeyRequestID
)

var (
	ErrIncorrectEmailOrPassword = fmt.Errorf("incorrect email or password")
	ErrNotAuthenticated         = fmt.Errorf("not authenticated")
)

func UserFromCtx(ctx context.Context) *datastruct.User {
	return ctx.Value(CtxKeyUser).(*datastruct.User)
}

func RequiredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}

		return nil
	}
}

func Validate(user *datastruct.User) error {
	return validation.ValidateStruct(
		user,
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.By(RequiredIf(user.EncryptedPassword == "")), validation.Length(6, 100)),
	)
}

func Sanitize(user *datastruct.User) {
	user.Password = ""
}

func ComparePassword(user *datastruct.User, password string) bool {
	result := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password)) == nil
	return result
}

func BeforeCreate(user *datastruct.User) error {
	if len(user.Password) > 0 {
		encryptedPassword, err := EncryptPassword(user.Password)
		if err != nil {
			return err
		}
		user.EncryptedPassword = encryptedPassword
	}
	return nil
}

func EncryptPassword(unencryptedPassword string) (string, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(unencryptedPassword), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(encryptedPassword), nil
}