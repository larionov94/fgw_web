package config

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
)

const (
	sessionName          = "fgw_session"
	SessionPerformerKey  = "performer_id"
	SessionRoleKey       = "role_id"
	SessionAuthPerformer = "authenticated"
)

var Store *sessions.CookieStore

func InitSessionStore() {
	secretKey := getSecretKey()

	Store = sessions.NewCookieStore([]byte(secretKey))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 дней по умолчанию
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode, // Важно: защита от CSRF
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

// Новые функции для безопасной работы с сессиями

// CreateSecureSession создает новую сессию с защитой от кэширования
func CreateSecureSession(w http.ResponseWriter, r *http.Request, performerID, roleID string) (*sessions.Session, error) {
	session, err := Store.Get(r, sessionName)
	if err != nil {
		return nil, err
	}

	// Очищаем старую сессию полностью
	session.Options.MaxAge = -1
	session.Save(r, w)

	// Создаем новую сессию
	session, err = Store.New(r, sessionName)
	if err != nil {
		return nil, err
	}

	// Устанавливаем значения
	session.Values[SessionAuthPerformer] = true
	session.Values[SessionPerformerKey] = performerID
	session.Values[SessionRoleKey] = roleID

	// Генерируем уникальный токен сессии для отслеживания
	sessionToken := generateSessionToken()
	session.Values["session_token"] = sessionToken

	// Время создания сессии
	session.Values["created_at"] = time.Now().Unix()

	// Устанавливаем опции безопасности
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   1800, // 30 минут активности
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	// Устанавливаем заголовки для предотвращения кэширования
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Устанавливаем дополнительную куку для управления историей
	setHistoryControlCookie(w, r)

	return session, session.Save(r, w)
}

// DestroySession полностью уничтожает сессию
func DestroySession(w http.ResponseWriter, r *http.Request) error {
	session, err := Store.Get(r, sessionName)
	if err != nil {
		return err
	}

	// Очищаем все значения
	for key := range session.Values {
		delete(session.Values, key)
	}

	// Устанавливаем отрицательный MaxAge для удаления куки
	session.Options.MaxAge = -1

	// Очищаем куку управления историей
	clearHistoryControlCookie(w, r)

	return session.Save(r, w)
}

// GetSession безопасно получает сессию с проверкой
func GetSession(r *http.Request) (*sessions.Session, error) {
	session, err := Store.Get(r, sessionName)
	if err != nil {
		return nil, err
	}

	// Проверяем, аутентифицирована ли сессия
	if auth, ok := session.Values[SessionAuthPerformer].(bool); !ok || !auth {
		return nil, nil
	}

	// Проверяем время жизни сессии (опционально)
	if createdAt, ok := session.Values["created_at"].(int64); ok {
		createTime := time.Unix(createdAt, 0)
		if time.Since(createTime) > 24*time.Hour { // Максимум 24 часа
			return nil, nil
		}
	}

	return session, nil
}

// RefreshSession обновляет сессию при активности
func RefreshSession(w http.ResponseWriter, r *http.Request) error {
	session, err := GetSession(r)
	if err != nil || session == nil {
		return err
	}

	// Обновляем время жизни при активности (если прошло больше 5 минут)
	if lastActivity, ok := session.Values["last_activity"].(int64); ok {
		lastActivityTime := time.Unix(lastActivity, 0)
		if time.Since(lastActivityTime) > 5*time.Minute {
			session.Values["last_activity"] = time.Now().Unix()
			session.Options.MaxAge = 1800 // Сбрасываем на 30 минут
			return session.Save(r, w)
		}
	} else {
		session.Values["last_activity"] = time.Now().Unix()
		return session.Save(r, w)
	}

	return nil
}

// Генерация уникального токена сессии
func generateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// Кука для управления историей браузера
func setHistoryControlCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "history_control",
		Value:    "no-cache",
		Path:     "/",
		MaxAge:   0,     // Сессионная кука
		HttpOnly: false, // Должна быть доступна в JS
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func clearHistoryControlCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "history_control",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

// GetSessionToken возвращает токен сессии для проверки
func GetSessionToken(session *sessions.Session) string {
	if token, ok := session.Values["session_token"].(string); ok {
		return token
	}
	return ""
}
