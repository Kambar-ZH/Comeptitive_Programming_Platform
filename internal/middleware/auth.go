package middleware

import (
	"context"
	"fmt"
	"site/internal/consts"
	"site/internal/datastruct"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

const (
	SessionName = "gopherschool"
)

var (
	ErrIncorrectEmailOrPassword = fmt.Errorf("incorrect email or password")
	ErrNotAuthenticated         = fmt.Errorf("not authenticated")
)

func UserFromCtx(ctx context.Context) (*datastruct.User, bool) {
	user, ok := ctx.Value(CtxKeyUser).(*datastruct.User)
	return user, ok
}

func LanguageCodeFromCtx(ctx context.Context) consts.LanguageCode {
	languageCode, ok := ctx.Value(CtxKeyLanguageCode).(int)
	if !ok {
		return consts.EN
	}
	return consts.LanguageCode(languageCode)
}

func RequiredIf(required bool) validation.RuleFunc {
	return func(value interface{}) error {
		if required {
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
	if len(user.Password) != 0 {
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
