package http_web

import (
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
	logg             *common.Logger
}

func NewPerformerHandlerHTML(performerService service.PerformerUseCase, logg *common.Logger) *PerformerHandlerHTML {
	return &PerformerHandlerHTML{performerService: performerService, logg: logg}
}

func (p *PerformerHandlerHTML) ServeHTTPHTMLRouter(mux *http.ServeMux) {
	mux.HandleFunc("/", p.ShowAuthForm)
	mux.HandleFunc("/fgw/performers", p.AllPerformersHTML)
	mux.HandleFunc("/login", p.AuthPerformerHTML)
	mux.HandleFunc("/fgw", p.StartPage)
	mux.HandleFunc("/fgw/performers/upd", p.UpdatePerformerHTML)
}

func (p *PerformerHandlerHTML) ShowAuthForm(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, tmplAuthHTML, nil, r)
}

func (p *PerformerHandlerHTML) AllPerformersHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != http.MethodGet {
		http_err.WriteMethodNotAllowed(w, r, p.logg, msg.H7000, "")

		return
	}

	performers, err := p.performerService.GetAllPerformers(r.Context())
	if err != nil {
		http_err.WriteServerError(w, r, p.logg, msg.H7001, err.Error())

		return
	}

	data := struct {
		Title      string
		Performers []model.Performer
	}{
		Title:      "Список сотрудников",
		Performers: performers,
	}

	p.renderPage(w, tmplPerformersHTML, data, r)
}

func (p *PerformerHandlerHTML) AuthPerformerHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != http.MethodPost {
		http_err.WriteMethodNotAllowed(w, r, p.logg, msg.H7000, "")

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
		http.Redirect(w, r, "/fgw", http.StatusFound)
	} else {
		http.Redirect(w, r, "/login?error"+url.QueryEscape(authResult.Message), http.StatusFound)
	}
}

func (p *PerformerHandlerHTML) UpdatePerformerHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != http.MethodPost {
		http_err.WriteMethodNotAllowed(w, r, p.logg, msg.H7000, "")

		return
	}

	if err := r.ParseForm(); err != nil {
		p.renderErrorPage(w, http.StatusBadRequest, msg.H7007, r)

		return
	}

	performerIdStr := r.FormValue("performerId")
	idRoleAFormsStr := r.FormValue("idRoleAForms")
	idRoleAFGWStr := r.FormValue("idRoleAFGW")
	updatedByStr := r.FormValue("updatedBy")

	if performerIdStr == "" || idRoleAFormsStr == "" || idRoleAFGWStr == "" || updatedByStr == "" {
		p.renderErrorPage(w, http.StatusUnauthorized, msg.E3214, r)

		return
	}

	performerId := convert.ConvStrToInt(performerIdStr)

	exists, err := p.performerService.ExistPerformer(r.Context(), performerId)
	if err != nil {
		http_err.WriteServerError(w, r, p.logg, msg.H7008, err.Error())

		return
	}

	if !exists {
		p.renderErrorPage(w, http.StatusUnauthorized, msg.E3212, r)

		return
	}

	performer := model.Performer{
		Id:           performerId,
		IdRoleAForms: convert.ConvStrToInt(idRoleAFormsStr),
		IdRoleAFGW:   convert.ConvStrToInt(idRoleAFGWStr),
		AuditRec: model.Audit{
			UpdatedAt: time.Now().String(),
			UpdatedBy: 6680, // TODO: заменить на авторизованного сотрудника
		},
	}

	if err = p.performerService.UpdPerformer(r.Context(), performerId, &performer); err != nil {
		http_err.WriteServerError(w, r, p.logg, msg.H7007, err.Error())

		return
	}

	http.Redirect(w, r, "/fgw/performers", http.StatusSeeOther)
}

func (p *PerformerHandlerHTML) StartPage(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, tmplStartPageHTML, nil, r)
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
	p.logg.LogWithResponseE(msgCode, statusCode, r.Method, r.URL.Path)
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
