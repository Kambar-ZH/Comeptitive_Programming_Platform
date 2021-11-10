package services

import (
	"site/internal/grpc/api"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

func requiredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}

		return nil
	}
}

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
		res, err := EncryptPassword(&api.EncryptPasswordRequest{UnencryptedPassword: user.Password})
		if err != nil {
			return err
		}
		user.EncryptedPassword = res.Password
	}
	return nil
}

func EncryptPassword(req *api.EncryptPasswordRequest) (*api.EncryptPasswordResponse, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(req.UnencryptedPassword), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	res := &api.EncryptPasswordResponse{
		Password: string(encryptedPassword),
	}
	return res, nil
}
