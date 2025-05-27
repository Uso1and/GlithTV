package session

import "github.com/gorilla/sessions"

var Store *sessions.CookieStore

func InitSessionStore(secretKey string) {
	Store = sessions.NewCookieStore([]byte(secretKey))
}
