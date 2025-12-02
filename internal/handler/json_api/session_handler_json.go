package json_api

import (
	"FGW_WEB/internal/config"
	"FGW_WEB/internal/service"
	"FGW_WEB/pkg/common"
	"encoding/json"
	"net/http"
	"time"
)

type AuthHandlerJSON struct {
	performerService service.PerformerUseCase
	logg             *common.Logger
}

func NewAuthHandlerJSON(performerService service.PerformerUseCase, logg *common.Logger) *AuthHandlerJSON {
	return &AuthHandlerJSON{performerService: performerService, logg: logg}
}

func (a *AuthHandlerJSON) ServeHTTPJSONRouter(mux *http.ServeMux) {
	mux.HandleFunc("/api/session-check", a.SessionCheckHandler)
}

func (a *AuthHandlerJSON) SessionCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем только HEAD и GET.
	if r.Method != http.MethodHead && r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем сессию.
	session, err := config.Store.Get(r, config.GetSessionName())
	if err != nil {
		// Нет сессии или ошибка.
		w.Header().Set("Session-Status", "no-session")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Проверяем аутентификацию.
	if auth, ok := session.Values[config.SessionAuthPerformer].(bool); !ok || !auth {
		w.Header().Set("Session-Status", "not-authenticated")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Проверяем время жизни сессии
	if createdAt, ok := session.Values["created_at"].(int64); ok {
		createTime := time.Unix(createdAt, 0)

		// 4 часа максимальное время (как в middleware).
		maxAge := 4 * time.Hour

		if time.Since(createTime) > maxAge {
			w.Header().Set("Session-Status", "expired")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	} else {
		// Нет времени создания
		w.Header().Set("Session-Status", "no-creation-time")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Обновляем время последней активности.
	session.Values["last_activity"] = time.Now().Unix()
	if err = session.Save(r, w); err != nil {
		return
	}

	// Сессия валидна.
	w.Header().Set("Session-Status", "active")
	w.WriteHeader(http.StatusOK)

	// Для GET запросов можно вернуть JSON.
	if r.Method == http.MethodGet {
		performerId, _ := session.Values[config.SessionPerformerKey].(int)
		roleId, _ := session.Values[config.SessionRoleKey].(int)
		createdAt, _ := session.Values["created_at"].(int64)

		response := map[string]interface{}{
			"status":      "active",
			"performerId": performerId,
			"roleId":      roleId,
			"createdAt":   time.Unix(createdAt, 0).Format("02.01.2006 15:04:05"),
			"sessionAge":  time.Since(time.Unix(createdAt, 0)).String(),
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			return
		}
	}
}
