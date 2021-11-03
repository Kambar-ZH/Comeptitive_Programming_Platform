package services

import (
	"site/internal/grpc/api"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

func Validate(user *api.User) error {
	return validation.ValidateStruct(
		user,
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.By(requiredIf(user.EncryptedPassword == "")), validation.Length(6, 100)),
	)
}

func Sanitize(user *api.User) {
	user.Password = ""
}

func ComparePassword(user *api.User, password string) bool {
	result := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password)) == nil
	return result
}

func BeforeCreate(user *api.User) error {
	if len(user.Password) > 0 {
		encPassword, err := EncryptString(&api.UnencryptedPassword{Password: user.Password})
		if err != nil {
			return err
		}
		user.EncryptedPassword = encPassword
	}
	return nil
}

func EncryptString(unEncpassword *api.UnencryptedPassword) (string, error) {
	encPassword, err := bcrypt.GenerateFromPassword([]byte(unEncpassword.Password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(encPassword), nil
}