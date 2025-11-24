package handler

import (
	"FGW_WEB/internal/config"
	"FGW_WEB/internal/handler/http_err"
	"FGW_WEB/pkg/common"
	"net/http"

	"github.com/gorilla/sessions"
)

type AuthMiddleware struct {
	store        *sessions.CookieStore
	sessName     string
	performerKey string
	roleKey      string
	logg         *common.Logger
}

func NewAuthMiddleware(store *sessions.CookieStore, logg *common.Logger) *AuthMiddleware {
	return &AuthMiddleware{store, config.GetSessionName(), config.SessionPerformerKey, config.SessionRoleKey, logg}
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := m.store.Get(r, m.sessName)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)

			return
		}

		if auth, ok := session.Values[config.SessionAuthPerformer].(bool); !ok || !auth {
			http.Redirect(w, r, "/", http.StatusFound)

			return
		}

		next.ServeHTTP(w, r)
	}
}

func (m *AuthMiddleware) RequireRole(requireRoles []int, next http.HandlerFunc) http.HandlerFunc {
	allowedRoles := make(map[int]bool)
	for _, role := range requireRoles {
		allowedRoles[role] = true
	}

	return func(w http.ResponseWriter, r *http.Request) {
		session, err := m.store.Get(r, m.sessName)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)

			return
		}

		performerRole, ok := session.Values[m.roleKey].(int)
		if !ok {
			http.Redirect(w, r, "/", http.StatusFound)

			return
		}

		if !allowedRoles[performerRole] {
			http_err.SendErrorHTTP(w, http.StatusForbidden, "Доступ запрещен: недостаточно прав.", m.logg, r)

			return
		}

		next.ServeHTTP(w, r)
	}
}

func (m *AuthMiddleware) GetPerformerId(r *http.Request) (int, bool) {
	session, err := m.store.Get(r, m.sessName)
	if err != nil {
		return 0, false
	}

	performerId, ok := session.Values[m.performerKey].(int)

	return performerId, ok
}

func (m *AuthMiddleware) GetRoleId(r *http.Request) (int, bool) {
	session, err := m.store.Get(r, m.sessName)
	if err != nil {
		return 0, false
	}

	performerRole, ok := session.Values[m.roleKey].(int)

	return performerRole, ok
}
