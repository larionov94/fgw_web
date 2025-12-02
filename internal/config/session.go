package config

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

const (
	sessionName          = "fgw_session"
	SessionPerformerKey  = "performer_id"
	SessionRoleKey       = "role_id"
	SessionAuthPerformer = "authenticated"
	maxAge               = 86400 * 7
	pathToDefault        = "/"
)

var Store *sessions.CookieStore

func InitSessionStore() {
	secretKey := getSecretKey()

	Store = sessions.NewCookieStore([]byte(secretKey))
	Store.Options = &sessions.Options{
		Path:     pathToDefault,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	log.Println("Сессия создана")
}

func getSecretKey() string {
	if secret := os.Getenv("SESSION_SECRET"); secret != "" {
		return secret
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic("Не удалось сгенерировать секретный ключ: " + err.Error())
	}

	return base64.StdEncoding.EncodeToString(key)
}

func GetSessionName() string {
	return sessionName
}
