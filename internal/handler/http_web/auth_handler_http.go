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
)

const (
	tmplAdminPerformerHTML = "admin.html"
)

type AuthHandlerHTML struct {
	performerService service.PerformerUseCase
	roleService      service.RoleUseCase
	logg             *common.Logger
	authMiddleware   *handler.AuthMiddleware
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
	mux.HandleFunc("/admin", a.authMiddleware.RequireAuth(a.authMiddleware.RequireRole([]int{3}, a.AuthPerformerHTML)))
}

func (a *AuthHandlerHTML) StartPageAdmin(w http.ResponseWriter, r *http.Request) {
	a.renderPage(w, tmplAdminPerformerHTML, nil, r)
}

func (a *AuthHandlerHTML) StartPage(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, config.GetSessionName())
	performerId := session.Values[config.SessionPerformerKey].(int)
	performerRole := session.Values[config.SessionRoleKey].(int)

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
	performerId, ok := a.authMiddleware.GetPerformerId(r)
	if ok && performerId > 0 {
		http.Redirect(w, r, "/fgw", http.StatusFound)

		return
	}
	a.renderPage(w, tmplAuthHTML, nil, r)
}

func (a *AuthHandlerHTML) LoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

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
	session, _ := config.Store.Get(r, config.GetSessionName())

	session.Values[config.SessionAuthPerformer] = false
	session.Values[config.SessionPerformerKey] = nil
	session.Values[config.SessionRoleKey] = nil
	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		a.logg.LogE("Ошибка при выходе: ", err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *AuthHandlerHTML) AuthPerformerHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

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
		session, _ := config.Store.Get(r, config.GetSessionName())
		session.Values[config.SessionAuthPerformer] = true
		session.Values[config.SessionPerformerKey] = performerId
		session.Values[config.SessionRoleKey] = authResult.Performer.IdRoleAForms

		err = session.Save(r, w)
		if err != nil {
			a.renderErrorPage(w, http.StatusInternalServerError, "Ошибка создания сессии", r)
			return
		}
		if authResult.Performer.IdRoleAForms == 3 {
			http.Redirect(w, r, "/admin", http.StatusFound)
		} else {
			http.Redirect(w, r, "/fgw", http.StatusFound)
		}
	} else {
		http.Redirect(w, r, "/login?error="+url.QueryEscape(authResult.Message), http.StatusFound)
	}
}

func (a *AuthHandlerHTML) renderErrorPage(w http.ResponseWriter, statusCode int, msgCode string, r *http.Request) {
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
	parseTmpl, err := template.New(tmpl).Funcs(
		template.FuncMap{
			"formatDateTime": convert.FormatDateTime,
		}).ParseFiles(prefixTmplPerformers + tmpl)
	if err != nil {
		a.renderErrorPage(w, http.StatusInternalServerError, msg.H7002+err.Error(), r)

		return
	}

	if err = parseTmpl.ExecuteTemplate(w, tmpl, data); err != nil {
		a.renderErrorPage(w, http.StatusInternalServerError, msg.H7003+err.Error(), r)

		return
	}
}
