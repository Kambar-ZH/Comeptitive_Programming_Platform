package middleware

type ctxKey int8

const (
	CtxKeyUser ctxKey = iota
	CtxKeyLanguageCode
	CtxKeyPage
	CtxKeyFilter
)
