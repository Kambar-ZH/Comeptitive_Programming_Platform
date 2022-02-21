package consts

type LanguageCode int

const (
	EN LanguageCode = iota
	RU
	KZ
)

func (l LanguageCode) String() string {
	return []string{"EN", "RU", "KZ"}[l]
}
