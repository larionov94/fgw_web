package handler

import (
	"FGW_WEB/internal/config"
	"FGW_WEB/internal/handler/http_err"
	"FGW_WEB/pkg/common"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
)

const (
	pathToDefault        = "/"
	prefixTmplPerformers = "web/html/"
	tmplForceLogoutHTML  = "force_logout.html"
)

type AuthMiddleware struct {
	store        *sessions.CookieStore
	sessName     string
	performerKey string
	roleKey      string
	logg         *common.Logger
}

func NewAuthMiddleware(store *sessions.CookieStore, logg *common.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		store:        store,
		sessName:     config.GetSessionName(),
		performerKey: config.SessionPerformerKey,
		roleKey:      config.SessionRoleKey,
		logg:         logg,
	}
}

// RequireAuth - основной middleware для проверки аутентификации.
func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем заголовки безопасности для всех защищенных запросов.
		m.setSecurityHeaders(w)

		// Получаем сессию с помощью безопасного метода.
		session, err := m.getSecureSession(r)
		if err != nil {
			m.forceLogoutAndRedirect(w, r, "Ошибка получения сессии")

			return
		}

		if session == nil {
			m.forceLogoutAndRedirect(w, r, "Сессия не найдена")

			return
		}

		// Проверяем аутентификацию.
		if auth, ok := session.Values[config.SessionAuthPerformer].(bool); !ok || !auth {
			m.forceLogoutAndRedirect(w, r, "Требуется аутентификация")

			return
		}

		// Проверяем время жизни сессии.
		if m.isSessionExpired(session) {
			m.forceLogoutAndRedirect(w, r, "Сессия истекла")

			return
		}

		// Обновляем активность сессии.
		m.updateSessionActivity(session, w, r)

		// Для HTML-ответов добавляем скрипт управления историей.
		if r.Header.Get("Accept") == "text/html" {
			m.addHistoryManagementScript(w)
		}

		next.ServeHTTP(w, r)
	}
}

// RequireRole - middleware для проверки ролей.
func (m *AuthMiddleware) RequireRole(requireRoles []int, next http.HandlerFunc) http.HandlerFunc {
	allowedRoles := make(map[int]bool)
	for _, role := range requireRoles {
		allowedRoles[role] = true
	}

	return m.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		session, err := m.store.Get(r, m.sessName)
		if err != nil {
			http.Redirect(w, r, pathToDefault, http.StatusFound)
			return
		}

		performerRole, ok := session.Values[m.roleKey].(int)
		if !ok {
			m.forceLogoutAndRedirect(w, r, "Роль не определена")
			return
		}

		if !allowedRoles[performerRole] {
			http_err.SendErrorHTTP(w, http.StatusForbidden, "Доступ запрещен: недостаточно прав.", m.logg, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetPerformerId - получение ID пользователя.
func (m *AuthMiddleware) GetPerformerId(r *http.Request) (int, bool) {
	session, err := m.store.Get(r, m.sessName)
	if err != nil {
		return 0, false
	}

	performerId, ok := session.Values[m.performerKey].(int)
	return performerId, ok
}

// GetRoleId - получение ID роли.
func (m *AuthMiddleware) GetRoleId(r *http.Request) (int, bool) {
	session, err := m.store.Get(r, m.sessName)
	if err != nil {
		return 0, false
	}

	performerRole, ok := session.Values[m.roleKey].(int)
	return performerRole, ok
}

// getSecureSession - безопасное получение сессии с валидацией.
func (m *AuthMiddleware) getSecureSession(r *http.Request) (*sessions.Session, error) {
	session, err := m.store.Get(r, m.sessName)
	if err != nil {
		return nil, err
	}

	// Дополнительная валидация куки.
	if session.IsNew {
		return nil, nil
	}

	return session, nil
}

// isSessionExpired - проверка истечения срока действия сессии.
func (m *AuthMiddleware) isSessionExpired(session *sessions.Session) bool {
	if createdAt, ok := session.Values["created_at"].(int64); ok {
		createTime := time.Unix(createdAt, 0)

		// Максимальное время жизни сессии - 4 часа.
		maxAge := 4 * time.Hour
		if customMaxAge, ok := session.Values["max_age"].(int); ok {
			maxAge = time.Duration(customMaxAge) * time.Second
		}

		return time.Since(createTime) > maxAge
	}

	return true
}

// updateSessionActivity - обновление времени активности.
func (m *AuthMiddleware) updateSessionActivity(session *sessions.Session, w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	session.Values["last_activity"] = now.Unix()

	// Устанавливаем куку с коротким временем жизни для браузера.
	if cookie, err := r.Cookie("activity_check"); err != nil || cookie.Value != "active" {
		http.SetCookie(w, &http.Cookie{
			Name:     "activity_check",
			Value:    "active",
			Path:     pathToDefault,
			MaxAge:   -1,
			HttpOnly: false,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	}

	if err := session.Save(r, w); err != nil {
		return
	}
}

// setSecurityHeaders - установка заголовков безопасности.
func (m *AuthMiddleware) setSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
}

// addHistoryManagementScript - добавление скрипта управления историей.
func (m *AuthMiddleware) addHistoryManagementScript(w http.ResponseWriter) {
	w.Header().Add("X-History-Control", "no-cache")
}

// forceLogoutAndRedirect - принудительный выход и редирект с очисткой истории.
func (m *AuthMiddleware) forceLogoutAndRedirect(w http.ResponseWriter, r *http.Request, reason string) {
	m.logg.LogW(fmt.Sprintf("Принудительный выход: %s", reason))

	// Уничтожаем сессию.
	if session, err := m.store.Get(r, m.sessName); err == nil {

		// Очищаем сессию.
		session.Options.MaxAge = -1
		for key := range session.Values {
			delete(session.Values, key)
		}
		if err = session.Save(r, w); err != nil {
			return
		}
	}

	// Устанавливаем заголовки no-cache.
	m.setSecurityHeaders(w)

	// Загружаем HTML шаблон
	tmpl, err := template.ParseFiles(prefixTmplPerformers + tmplForceLogoutHTML)
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок ответа.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Выполняем шаблон.
	if err = tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Ошибка рендеринга шаблона", http.StatusInternalServerError)
	}
}
