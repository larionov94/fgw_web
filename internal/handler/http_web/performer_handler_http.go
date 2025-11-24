package http_web

import (
	"FGW_WEB/internal/config"
	"FGW_WEB/internal/handler"
	"FGW_WEB/internal/handler/http_err"
	"FGW_WEB/internal/model"
	"FGW_WEB/internal/service"
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"FGW_WEB/pkg/convert"
	"html/template"
	"net/http"
	"net/url"
	"time"
)

const (
	tmplPerformersHTML   = "performers.html"
	tmplErrorHTML        = "error.html"
	tmplAuthHTML         = "auth.html"
	tmplStartPageHTML    = "index.html" // /fgw
	prefixTmplPerformers = "web/html/"
)

type PerformerHandlerHTML struct {
	performerService service.PerformerUseCase
	roleService      service.RoleUseCase
	logg             *common.Logger
	authMiddleware   *handler.AuthMiddleware
}

func NewPerformerHandlerHTML(performerService service.PerformerUseCase, roleService service.RoleUseCase, logg *common.Logger, authMiddleware *handler.AuthMiddleware) *PerformerHandlerHTML {
	return &PerformerHandlerHTML{performerService: performerService, roleService: roleService, logg: logg, authMiddleware: authMiddleware}
}

func (p *PerformerHandlerHTML) ServeHTTPHTMLRouter(mux *http.ServeMux) {
	mux.HandleFunc("/", p.ShowAuthForm)
	mux.HandleFunc("/login", p.AuthPerformerHTML)
	mux.HandleFunc("/logout", p.Logout)
	mux.HandleFunc("/fgw", p.authMiddleware.RequireAuth(p.StartPage))
	mux.HandleFunc("/fgw/performers", p.authMiddleware.RequireAuth(p.authMiddleware.RequireRole([]int{3}, p.AllPerformersHTML)))
	mux.HandleFunc("/fgw/performers/upd", p.authMiddleware.RequireAuth(p.authMiddleware.RequireRole([]int{3}, p.UpdPerformerHTML)))
}

func (p *PerformerHandlerHTML) ShowAuthForm(w http.ResponseWriter, r *http.Request) {
	performerId, ok := p.authMiddleware.GetPerformerId(r)
	if ok && performerId > 0 {
		http.Redirect(w, r, "/fgw", http.StatusFound)

		return
	}
	p.renderPage(w, tmplAuthHTML, nil, r)
}

func (p *PerformerHandlerHTML) AllPerformersHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != http.MethodGet {
		http_err.SendErrorHTTP(w, http.StatusMethodNotAllowed, "", p.logg, r)

		return
	}

	performers, err := p.performerService.GetAllPerformers(r.Context())
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusInternalServerError, err.Error(), p.logg, r)

		return
	}

	roles, err := p.roleService.GetAllRole(r.Context())
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusInternalServerError, err.Error(), p.logg, r)

		return
	}

	data := struct {
		Title      string
		Performers []*model.Performer
		Roles      []*model.Role
	}{
		Title:      "Список сотрудников",
		Performers: performers,
		Roles:      roles,
	}

	if performerIdStr := r.URL.Query().Get("performerId"); performerIdStr != "" {
		p.markEditingPerformer(performerIdStr, performers)
	}

	p.renderPage(w, tmplPerformersHTML, data, r)
}

func (p *PerformerHandlerHTML) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, config.GetSessionName())

	session.Values[config.SessionAuthPerformer] = false
	session.Values[config.SessionPerformerKey] = nil
	session.Values[config.SessionRoleKey] = nil
	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		p.logg.LogE("Ошибка при выходе: ", err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (p *PerformerHandlerHTML) AuthPerformerHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != http.MethodPost {
		http_err.SendErrorHTTP(w, http.StatusMethodNotAllowed, "", p.logg, r)

		return
	}

	if err := r.ParseForm(); err != nil {
		p.renderErrorPage(w, http.StatusBadRequest, msg.H7007, r)

		return
	}

	performerIdStr := r.FormValue("performerId")
	performerPass := r.FormValue("performerPassword")

	if performerIdStr == "" || performerPass == "" {
		p.renderErrorPage(w, http.StatusUnauthorized, msg.E3211, r)

		return
	}

	performerId := convert.ConvStrToInt(performerIdStr)

	authResult, err := p.performerService.AuthPerformer(r.Context(), performerId, performerPass)
	if err != nil {
		if authResult != nil && !authResult.Success {
			p.renderErrorPage(w, http.StatusUnauthorized, authResult.Message, r)
		} else {
			p.renderErrorPage(w, http.StatusUnauthorized, msg.H7005, r)
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
			p.renderErrorPage(w, http.StatusInternalServerError, "Ошибка создания сессии", r)

			return
		}
		http.Redirect(w, r, "/fgw", http.StatusFound)
	} else {
		http.Redirect(w, r, "/login?error"+url.QueryEscape(authResult.Message), http.StatusFound)
	}
}

func (p *PerformerHandlerHTML) UpdPerformerHTML(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p.processUpdFormPerformer(w, r)
	case http.MethodGet:
		p.renderUpdFormPerformer(w, r)
	default:
		http_err.SendErrorHTTP(w, http.StatusMethodNotAllowed, "", p.logg, r)
	}
}

func (p *PerformerHandlerHTML) renderUpdFormPerformer(w http.ResponseWriter, r *http.Request) {
	performerIdStr := r.URL.Query().Get("performerId")
	http.Redirect(w, r, "/fgw/performerId?performer="+performerIdStr, http.StatusFound)
}

func (p *PerformerHandlerHTML) processUpdFormPerformer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := r.ParseForm(); err != nil {
		p.renderErrorPage(w, http.StatusBadRequest, msg.H7007, r)

		return
	}

	performerIdStr := r.FormValue("performerId")
	idRoleAFormsStr := r.FormValue("idRoleAForms")
	idRoleAFGWStr := r.FormValue("idRoleAFGW")
	updatedByStr := r.FormValue("updatedBy")

	if idRoleAFormsStr == "" || idRoleAFGWStr == "" || updatedByStr == "" {
		p.renderErrorPage(w, http.StatusUnauthorized, msg.E3214, r)

		return
	}

	performerId := convert.ConvStrToInt(performerIdStr)
	updatedBy := convert.ConvStrToInt(updatedByStr)
	idRoleAForms := convert.ConvStrToInt(idRoleAFormsStr)
	idRoleAFGW := convert.ConvStrToInt(idRoleAFGWStr)

	exists, err := p.performerService.ExistPerformer(r.Context(), performerId)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusInternalServerError, err.Error(), p.logg, r)

		return
	}

	if !exists {
		p.renderErrorPage(w, http.StatusUnauthorized, msg.E3212, r)

		return
	}

	performer := model.Performer{
		Id:           performerId,
		IdRoleAForms: idRoleAForms,
		IdRoleAFGW:   idRoleAFGW,
		AuditRec: model.Audit{
			UpdatedAt: time.Now().String(),
			UpdatedBy: updatedBy, // TODO: заменить на авторизованного сотрудника
		},
	}

	if err := p.performerService.UpdPerformer(r.Context(), performerId, &performer); err != nil {
		http_err.SendErrorHTTP(w, http.StatusInternalServerError, err.Error(), p.logg, r)

		return
	}

	http.Redirect(w, r, "/fgw/performers", http.StatusSeeOther)
}

func (p *PerformerHandlerHTML) markEditingPerformer(id string, performers []*model.Performer) {
	performerId := convert.ConvStrToInt(id)
	for _, performer := range performers {
		if performer.Id == performerId {
			performer.IsEditing = true
		}
	}
}

func (p *PerformerHandlerHTML) StartPage(w http.ResponseWriter, r *http.Request) {
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

	p.renderPage(w, tmplStartPageHTML, data, r)
}

func (p *PerformerHandlerHTML) renderErrorPage(w http.ResponseWriter, statusCode int, msgCode string, r *http.Request) {
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
	p.logg.LogHttpErr(msgCode, statusCode, r.Method, r.URL.Path)
	p.renderPage(w, tmplErrorHTML, data, r)
}

func (p *PerformerHandlerHTML) renderPage(w http.ResponseWriter, tmpl string, data interface{}, r *http.Request) {
	parseTmpl, err := template.New(tmpl).Funcs(
		template.FuncMap{
			"formatDateTime": convert.FormatDateTime,
		}).ParseFiles(prefixTmplPerformers + tmpl)
	if err != nil {
		p.renderErrorPage(w, http.StatusInternalServerError, msg.H7002+err.Error(), r)

		return
	}

	if err = parseTmpl.ExecuteTemplate(w, tmpl, data); err != nil {
		p.renderErrorPage(w, http.StatusInternalServerError, msg.H7003+err.Error(), r)

		return
	}
}
