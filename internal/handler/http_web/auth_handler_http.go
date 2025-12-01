package http_web

import (
	"FGW_WEB/internal/config"
	"FGW_WEB/internal/handler"
	"FGW_WEB/internal/handler/http_err"
	"FGW_WEB/internal/service"
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"FGW_WEB/pkg/convert"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/sessions"
)

const (
	tmplAdminPerformerHTML = "admin.html"
	tmplRedirectHTML       = "redirect.html" // Общий шаблон для всех редиректов
)

// Константы для редиректов
const (
	RedirectDelayFast    = 100  // 0.1 секунда
	RedirectDelayNormal  = 300  // 0.3 секунды
	FallbackDelayDefault = 3000 // 3 секунды
)

type AuthHandlerHTML struct {
	performerService service.PerformerUseCase
	roleService      service.RoleUseCase
	logg             *common.Logger
	authMiddleware   *handler.AuthMiddleware
}

type RedirectData struct {
	Title           string
	Message         string
	NoScriptMessage string
	TargetURL       string
	CurrentURL      string
	TempURL         string
	Delay           int
	FallbackDelay   int
	ClearHistory    bool
	AddTempState    bool // Флаг для сложного управления историей
}

func NewAuthHandlerHTML(
	performerService service.PerformerUseCase,
	roleService service.RoleUseCase,
	logg *common.Logger,
	authMiddleware *handler.AuthMiddleware) *AuthHandlerHTML {

	return &AuthHandlerHTML{
		performerService: performerService,
		roleService:      roleService,
		logg:             logg,
		authMiddleware:   authMiddleware}
}

func (a *AuthHandlerHTML) ServerHTTPRouter(mux *http.ServeMux) {
	mux.HandleFunc("/", a.ShowAuthForm)
	mux.HandleFunc("/login", a.LoginPage)
	mux.HandleFunc("/auth", a.AuthPerformerHTML)
	mux.HandleFunc("/logout", a.Logout)
	mux.HandleFunc("/fgw", a.authMiddleware.RequireAuth(a.StartPage))
	mux.HandleFunc("/admin", a.authMiddleware.RequireAuth(a.authMiddleware.RequireRole([]int{3}, a.StartPageAdmin)))
}

func (a *AuthHandlerHTML) StartPageAdmin(w http.ResponseWriter, r *http.Request) {
	performerId, ok1 := a.authMiddleware.GetPerformerId(r)
	performerRole, ok2 := a.authMiddleware.GetRoleId(r)

	if !ok1 || !ok2 {
		a.redirectToLoginWithHistoryClear(w, r)
		return
	}

	a.setSecureHTMLHeaders(w)

	data := struct {
		PerformerId   int
		PerformerRole int
	}{
		PerformerId:   performerId,
		PerformerRole: performerRole,
	}

	a.renderPage(w, tmplAdminPerformerHTML, data, r)
}

func (a *AuthHandlerHTML) StartPage(w http.ResponseWriter, r *http.Request) {
	performerId, ok1 := a.authMiddleware.GetPerformerId(r)
	performerRole, ok2 := a.authMiddleware.GetRoleId(r)

	if !ok1 || !ok2 {
		a.redirectToLoginWithHistoryClear(w, r)
		return
	}

	a.setSecureHTMLHeaders(w)

	data := struct {
		PerformerId   int
		PerformerRole int
	}{
		PerformerId:   performerId,
		PerformerRole: performerRole,
	}

	a.renderPage(w, tmplStartPageHTML, data, r)
}

func (a *AuthHandlerHTML) ShowAuthForm(w http.ResponseWriter, r *http.Request) {
	session, err := config.Store.Get(r, config.GetSessionName())
	if err == nil {
		if auth, ok := session.Values[config.SessionAuthPerformer].(bool); ok && auth {
			a.safeRedirectBasedOnRole(w, r, session)
			return
		}
	}

	a.LoginPage(w, r)
}

func (a *AuthHandlerHTML) LoginPage(w http.ResponseWriter, r *http.Request) {
	a.setSecureHTMLHeaders(w)

	if r.Method != http.MethodGet {
		http_err.SendErrorHTTP(w, http.StatusMethodNotAllowed, "", a.logg, r)
		return
	}

	errorMsg := r.URL.Query().Get("error")

	data := struct {
		ErrorMessage string
	}{
		ErrorMessage: errorMsg,
	}

	a.renderPage(w, tmplAuthHTML, data, r)
}

func (a *AuthHandlerHTML) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := config.Store.Get(r, config.GetSessionName())
	if err != nil {
		a.sendLogoutPageWithHistoryClear(w, r)
		return
	}

	if token, ok := session.Values["session_token"].(string); ok {
		if mw, ok := interface{}(a.authMiddleware).(interface{ RemoveSessionToken(token string) }); ok {
			mw.RemoveSessionToken(token)
		}
	}

	for key := range session.Values {
		delete(session.Values, key)
	}

	session.Options.MaxAge = -1
	session.Options.HttpOnly = true
	session.Options.Secure = true
	session.Options.SameSite = http.SameSiteStrictMode

	session.Save(r, w)

	http.SetCookie(w, &http.Cookie{
		Name:     config.GetSessionName(),
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	a.sendLogoutPageWithHistoryClear(w, r)
}

func (a *AuthHandlerHTML) AuthPerformerHTML(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http_err.SendErrorHTTP(w, http.StatusMethodNotAllowed, "", a.logg, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		a.renderErrorPage(w, http.StatusBadRequest, msg.H7007, r)
		return
	}

	performerIdStr := r.FormValue("performerId")
	performerPass := r.FormValue("performerPassword")

	if performerIdStr == "" || performerPass == "" {
		a.renderErrorPage(w, http.StatusUnauthorized, msg.E3211, r)
		return
	}

	performerId := convert.ConvStrToInt(performerIdStr)

	authResult, err := a.performerService.AuthPerformer(r.Context(), performerId, performerPass)
	if err != nil {
		if authResult != nil && !authResult.Success {
			http.Redirect(w, r, "/login?error="+url.QueryEscape(authResult.Message), http.StatusFound)
		} else {
			http.Redirect(w, r, "/login?error="+url.QueryEscape(msg.H7005), http.StatusFound)
		}
		return
	}

	if authResult.Success {
		err := a.createSecureSession(w, r, performerId, authResult.Performer.IdRoleAForms)
		if err != nil {
			a.renderErrorPage(w, http.StatusInternalServerError, "Ошибка создания сессии", r)
			return
		}

		a.sendLoginSuccessPage(w, r, authResult.Performer.IdRoleAForms)
	} else {
		http.Redirect(w, r, "/login?error="+url.QueryEscape(authResult.Message), http.StatusFound)
	}
}

// НОВЫЙ МЕТОД: safeRedirectBasedOnRole с использованием общего шаблона
func (a *AuthHandlerHTML) safeRedirectBasedOnRole(w http.ResponseWriter, r *http.Request, session *sessions.Session) {
	target := "/fgw"
	if role, ok := session.Values[config.SessionRoleKey].(int); ok && role == 3 {
		target = "/admin"
	}

	data := RedirectData{
		Title:           "Перенаправление",
		Message:         "Вы уже авторизованы. Выполняется безопасное перенаправление...",
		NoScriptMessage: "Включите JavaScript для безопасного перехода.",
		TargetURL:       target,
		CurrentURL:      r.URL.Path,
		TempURL:         "/temp-redirect",
		Delay:           RedirectDelayFast,
		FallbackDelay:   FallbackDelayDefault,
		ClearHistory:    true,
		AddTempState:    false, // Для этого случая не нужно сложное управление историей
	}

	a.renderRedirectPage(w, r, data)
}

// УНИФИЦИРОВАННЫЙ МЕТОД ДЛЯ ВСЕХ РЕДИРЕКТОВ
func (a *AuthHandlerHTML) renderRedirectPage(w http.ResponseWriter, r *http.Request, data RedirectData) {
	// Устанавливаем значения по умолчанию
	if data.Title == "" {
		data.Title = "Перенаправление"
	}
	if data.Message == "" {
		data.Message = "Выполняется безопасное перенаправление..."
	}
	if data.NoScriptMessage == "" {
		data.NoScriptMessage = "Включите JavaScript для безопасного перехода."
	}
	if data.CurrentURL == "" {
		data.CurrentURL = r.URL.Path
	}
	if data.TempURL == "" && data.AddTempState {
		//data.TempURL = "/temp-redirect-" + convert.IntToStr(int(time.Now().UnixNano()%1000000))
	}
	if data.Delay == 0 {
		data.Delay = RedirectDelayNormal
	}
	if data.FallbackDelay == 0 {
		data.FallbackDelay = FallbackDelayDefault
	}

	a.setSecureHTMLHeaders(w)
	a.renderPage(w, tmplRedirectHTML, data, r)
}

// Обновленный sendLoginSuccessPage
func (a *AuthHandlerHTML) sendLoginSuccessPage(w http.ResponseWriter, r *http.Request, roleId int) {
	target := "/fgw"
	if roleId == 3 {
		target = "/admin"
	}

	data := RedirectData{
		Title:           "Успешный вход",
		Message:         "Вход выполнен успешно. Выполняется безопасное перенаправление...",
		NoScriptMessage: "Включите JavaScript для безопасного перехода.",
		TargetURL:       target,
		CurrentURL:      "/auth",
		TempURL:         "/login-temp-redirect",
		Delay:           RedirectDelayNormal,
		FallbackDelay:   2000,
		ClearHistory:    true,
		AddTempState:    true, // Для входа нужна полная очистка истории
	}

	a.renderRedirectPage(w, r, data)
}

// Обновленный sendLogoutPageWithHistoryClear
func (a *AuthHandlerHTML) sendLogoutPageWithHistoryClear(w http.ResponseWriter, r *http.Request) {
	data := RedirectData{
		Title:           "Выход из системы",
		Message:         "Вы успешно вышли из системы. Выполняется безопасное перенаправление на страницу входа...",
		NoScriptMessage: "Включите JavaScript для безопасного выхода.",
		TargetURL:       "/login",
		CurrentURL:      r.URL.Path,
		TempURL:         "/logout-temp-redirect",
		Delay:           RedirectDelayNormal,
		FallbackDelay:   FallbackDelayDefault,
		ClearHistory:    true,
		AddTempState:    true, // Для выхода нужна полная очистка истории
	}

	a.renderRedirectPage(w, r, data)
}

func (a *AuthHandlerHTML) redirectToLoginWithHistoryClear(w http.ResponseWriter, r *http.Request) {
	a.sendLogoutPageWithHistoryClear(w, r)
}

// Вспомогательные методы

func (a *AuthHandlerHTML) setSecureHTMLHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
}

func (a *AuthHandlerHTML) createSecureSession(w http.ResponseWriter, r *http.Request, performerId, roleId int) error {
	session, _ := config.Store.Get(r, config.GetSessionName())

	token := config.GenerateSessionToken()

	session.Values[config.SessionAuthPerformer] = true
	session.Values[config.SessionPerformerKey] = performerId
	session.Values[config.SessionRoleKey] = roleId
	session.Values["session_token"] = token
	session.Values["created_at"] = time.Now().Unix()
	session.Values["last_activity"] = time.Now().Unix()

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   1800,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	a.setSecureHTMLHeaders(w)

	return session.Save(r, w)
}

func (a *AuthHandlerHTML) renderErrorPage(w http.ResponseWriter, statusCode int, msgCode string, r *http.Request) {
	a.setSecureHTMLHeaders(w)

	data := struct {
		Title      string
		MsgCode    string
		StatusCode int
		Method     string
		Path       string
	}{
		Title:      "Ошибка",
		MsgCode:    msgCode,
		StatusCode: statusCode,
		Method:     r.Method,
		Path:       r.URL.Path,
	}

	w.WriteHeader(statusCode)
	a.logg.LogHttpErr(msgCode, statusCode, r.Method, r.URL.Path)
	a.renderPage(w, tmplErrorHTML, data, r)
}

func (a *AuthHandlerHTML) renderPage(w http.ResponseWriter, tmpl string, data interface{}, r *http.Request) {
	templatePath := prefixTmplPerformers + tmpl

	parseTmpl, err := template.New(tmpl).Funcs(
		template.FuncMap{
			"formatDateTime": convert.FormatDateTime,
		}).ParseFiles(templatePath)

	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		a.logg.LogE("Ошибка загрузки шаблона:", err)
		return
	}

	if err = parseTmpl.ExecuteTemplate(w, tmpl, data); err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		a.logg.LogE("Ошибка выполнения шаблона:", err)
		return
	}
}
