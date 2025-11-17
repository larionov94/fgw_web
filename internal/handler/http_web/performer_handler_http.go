package http_web

import (
	"FGW_WEB/internal/handler/http_err"
	"FGW_WEB/internal/service"
	"FGW_WEB/internal/service/dto"
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"FGW_WEB/pkg/convert"
	"html/template"
	"net/http"
)

const (
	tmplPerformersHTML   = "performers.html"
	tmplErrorHTML        = "error.html"
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
	mux.HandleFunc("/fgw/performers", p.AllPerformersHTML)
	mux.HandleFunc("/fgw/login", p.AuthPerformerHTML)
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
		Performers []dto.PerformerDTO
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

	//performers, err := p.performerService.GetAllPerformerAuth(r.Context())
	//if err != nil {
	//	http_err.WriteServerError(w, r, p.logg, msg.H7001, err.Error())
	//
	//	return
	//}

	//performerId := convert.ParseFormFieldInt(r,"id")
	//performerPass := r.FormValue("pass")

	//for _, performer := range performers {
	//	if performer.Id == performerId && performer.Pass == performerPass {
	//
	//	}
	//}
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
