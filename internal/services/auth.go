package services

import (
	"site/internal/datastruct"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Validate(user *datastruct.User) error
	Sanitize(user *datastruct.User)
	ComparePassword(user *datastruct.User, password string) bool
	BeforeCreate(user *datastruct.User) error
	EncryptPassword(unencryptedPassword string) (string, error)
}

type AuthServiceImpl struct {}

func NewAuthService(opts ...AuthServiceOption) AuthService {
	svc := &AuthServiceImpl{}
	for _, v := range(opts) {
		v(svc)
	}
	return svc
}

func requiredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}

		return nil
	}
}

func (a AuthServiceImpl) Validate(user *datastruct.User) error {
	return validation.ValidateStruct(
		user,
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.By(requiredIf(user.EncryptedPassword == "")), validation.Length(6, 100)),
	)
}

func (a AuthServiceImpl) Sanitize(user *datastruct.User) {
	user.Password = ""
}

func (a AuthServiceImpl) ComparePassword(user *datastruct.User, password string) bool {
	result := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password)) == nil
	return result
}

func (a AuthServiceImpl) BeforeCreate(user *datastruct.User) error {
	if len(user.Password) > 0 {
		encryptedPassword, err := a.EncryptPassword(user.Password)
		if err != nil {
			return err
		}
		user.EncryptedPassword = encryptedPassword
	}
	return nil
}

func (a AuthServiceImpl) EncryptPassword(unencryptedPassword string) (string, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(unencryptedPassword), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(encryptedPassword), nil
}