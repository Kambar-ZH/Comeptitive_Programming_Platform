package dto

import "site/internal/consts"

type Cridentials struct {
	Email        string              `json:"email"`
	Password     string              `json:"password"`
	LanguageCode consts.LanguageCode `json:"language_code"`
}
