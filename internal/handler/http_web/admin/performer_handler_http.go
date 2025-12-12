package admin

import (
	"FGW_WEB/internal/config"
	"FGW_WEB/internal/handler"
	"FGW_WEB/internal/handler/http_err"
	"FGW_WEB/internal/handler/http_web"
	"FGW_WEB/internal/handler/json_api"
	"FGW_WEB/internal/handler/json_err"
	"FGW_WEB/internal/model"
	"FGW_WEB/internal/service"
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"FGW_WEB/pkg/convert"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

const (
	tmplAdminPerformersHTML = "performers.html"
	tmplErrorHTML           = "error.html"
	tmplAdminHTML           = "admin.html"
	prefixTmplAdmin         = "web/html/admin/"

	prefixDefaultTmpl = "web/html/"
	prefixAdminTmpl   = "web/html/admin/"

	pageSize          = 55
	maxPage           = 5
	numberPageDefault = 1
)

var authPerformerId int

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
	mux.HandleFunc("/admin/performers", p.authMiddleware.RequireAuth(p.authMiddleware.RequireRole([]int{3}, p.AllPerformersHTML)))
	mux.HandleFunc("/admin/performers/upd", p.authMiddleware.RequireAuth(p.authMiddleware.RequireRole([]int{3}, p.HandleJSONUpdate)))
}

func (p *PerformerHandlerHTML) AllPerformersHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != http.MethodGet {
		http_err.SendErrorHTTP(w, http.StatusMethodNotAllowed, "", p.logg, r)

		return
	}

	performerId, performerRoleId, err := p.getSessionPerformerData(w, r)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusUnauthorized, err.Error(), p.logg, r)

		return
	}

	pageStr := r.URL.Query().Get("page")
	searchPattern := r.URL.Query().Get("search")

	page, err := http_web.GetParametersPagination(pageStr, numberPageDefault)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusNotFound, "", p.logg, r)

		return
	}

	totalCount, performers, err, done := p.searchPerformerWithPagination(w, r, page, searchPattern, err)
	if done {
		return
	}

	roles, err := p.roleService.GetAllRole(r.Context())
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusNotFound, err.Error(), p.logg, r)

		return
	}

	performer, err := p.performerService.FindByIdPerformer(r.Context(), performerId)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusNotFound, err.Error(), p.logg, r)

		return
	}
	totalPages, err := http_web.CalculatePage(totalCount, pageSize, page)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusNotFound, err.Error(), p.logg, r)

		return
	}
	role, err := p.roleService.FindRoleById(r.Context(), performerRoleId)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusNotFound, err.Error(), p.logg, r)

		return
	}
	pages := http_web.GeneratePageRange(page, totalPages, maxPage)

	countPerformers := len(performers)
	startItem, endItem, err := http_web.CalculateRangeOfElements((page-1)*pageSize, totalCount, countPerformers)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusNotFound, err.Error(), p.logg, r)

		return

	}
	data := struct {
		Title         string
		CurrentPage   string
		Performers    []*model.Performer
		Roles         []*model.Role
		PerformerFIO  string
		PerformerId   int
		PerformerRole string
		Pagination    model.Pagination
		SearchQuery   string
		IsSearch      bool
	}{
		Title:         "Список сотрудников",
		CurrentPage:   "performers",
		Performers:    performers,
		Roles:         roles,
		PerformerFIO:  performer.FIO,
		PerformerId:   performerId,
		PerformerRole: role.Name,
		Pagination: model.Pagination{
			Page:           page,
			PageSize:       pageSize,
			TotalCount:     totalCount,
			TotalPages:     totalPages,
			Pages:          pages,
			StartItem:      startItem,
			EndItem:        endItem,
			PerformerIdStr: performerId,
		},
		SearchQuery: searchPattern,
		IsSearch:    searchPattern != "",
	}

	p.renderPages(w, tmplAdminHTML, data, r, tmplAdminPerformersHTML)
}

// searchPerformerWithPagination поиск сотрудника с пагинацией.
func (p *PerformerHandlerHTML) searchPerformerWithPagination(w http.ResponseWriter, r *http.Request, page int, searchPattern string, err error) (int, []*model.Performer, error, bool) {
	var totalCount int
	var performers []*model.Performer
	var offset = (page - 1) * pageSize

	if searchPattern != "" {
		performers, err = p.performerService.SearchPerformerById(r.Context(), searchPattern)
		if err != nil {
			http_err.SendErrorHTTP(w, http.StatusNotFound, "", p.logg, r)

			return 0, nil, nil, true
		}
		totalCount = len(performers)

		start := offset
		end := start + pageSize
		if start > totalCount {
			start = 0
		}
		if end > totalCount {
			end = totalCount
		}
		performers = performers[start:end]
	} else {
		totalCount, err = p.performerService.GetPerformersCount(r.Context())
		if err != nil {
			http_err.SendErrorHTTP(w, http.StatusNotFound, err.Error(), p.logg, r)

			return 0, nil, nil, true
		}

		performers, err = p.performerService.GetPerformersWithPagination(r.Context(), offset, pageSize)
		if err != nil {
			http_err.SendErrorHTTP(w, http.StatusNotFound, err.Error(), p.logg, r)

			return 0, nil, nil, true
		}
	}

	return totalCount, performers, nil, false
}

// HandleJSONUpdate обработчик для JSON запросов от Fetch API.
func (p *PerformerHandlerHTML) HandleJSONUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var req struct {
		PerformerId  int `json:"performerId"`
		IdRoleAForms int `json:"idRoleAForms"`
		IdRoleAFGW   int `json:"idRoleAFGW"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json_err.SendErrorResponse(w, http.StatusBadRequest, msg.H7004, err.Error(), r)

		return
	}

	exists, err := p.performerService.ExistPerformer(r.Context(), req.PerformerId)
	if err != nil {
		json_err.SendErrorResponse(w, http.StatusInternalServerError, msg.H7001, err.Error(), r)

		return
	}

	if !exists {
		json_err.SendErrorResponse(w, http.StatusNotFound, msg.H7008, "", r)

		return
	}

	if session, err := config.Store.Get(r, config.GetSessionName()); err == nil {
		if id, ok := session.Values[config.SessionPerformerKey].(int); ok {
			authPerformerId = id
		}
	}

	performer := model.Performer{
		Id:           req.PerformerId,
		IdRoleAForms: req.IdRoleAForms,
		IdRoleAFGW:   req.IdRoleAFGW,
		AuditRec: model.Audit{
			UpdatedBy: authPerformerId,
		},
	}

	if err = p.performerService.UpdPerformer(r.Context(), req.PerformerId, &performer); err != nil {
		json_err.SendErrorResponse(w, http.StatusInternalServerError, msg.H7001, err.Error(), r)

		return
	}

	response := map[string]interface{}{
		"success":     true,
		"message":     "Роли успешно обновлены",
		"performerId": req.PerformerId,
		"updatedAt":   time.Now().Format("2006-01-02 15:04:05"),
		"updatedBy":   authPerformerId,
	}

	w.WriteHeader(http.StatusOK)
	json_api.WriteJSON(w, response, r)
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
			"add":            func(a, b int) int { return a + b },
			"sub":            func(a, b int) int { return a - b },
		}).ParseFiles(prefixDefaultTmpl + tmpl)
	if err != nil {
		p.renderErrorPage(w, http.StatusInternalServerError, msg.H7002+err.Error(), r)

		return
	}

	if err = parseTmpl.ExecuteTemplate(w, tmpl, data); err != nil {
		p.renderErrorPage(w, http.StatusInternalServerError, msg.H7003+err.Error(), r)

		return
	}
}

func (p *PerformerHandlerHTML) renderPages(
	w http.ResponseWriter, tmpl string, data interface{}, r *http.Request, addTemplates ...string) {

	templatePaths := []string{prefixDefaultTmpl + tmpl}

	for _, addTmpl := range addTemplates {
		templatePaths = append(templatePaths, prefixAdminTmpl+addTmpl)
	}

	parseTmpl, err := template.New(tmpl).Funcs(
		template.FuncMap{
			"formatDateTime": convert.FormatDateTime,
			"add":            func(a, b int) int { return a + b },
			"sub":            func(a, b int) int { return a - b },
		}).ParseFiles(templatePaths...)

	if err != nil {
		p.renderErrorPage(w, http.StatusInternalServerError, msg.H7002+err.Error(), r)

		return
	}

	if err = parseTmpl.ExecuteTemplate(w, tmpl, data); err != nil {
		p.renderErrorPage(w, http.StatusInternalServerError, msg.H7003+err.Error(), r)

		return
	}
}

// getSessionPerformerData получить данные о сеансе сотрудника.
func (p *PerformerHandlerHTML) getSessionPerformerData(w http.ResponseWriter, r *http.Request) (int, int, error) {
	performerId, ok := p.authMiddleware.GetPerformerId(r)
	if !ok {
		http_err.SendErrorHTTP(w, http.StatusUnauthorized, msg.H7005, p.logg, r)

		return 0, 0, fmt.Errorf("%s", msg.H7005)
	}
	authPerformerId = performerId

	performerRole, ok := p.authMiddleware.GetRoleId(r)
	if !ok {
		http_err.SendErrorHTTP(w, http.StatusUnauthorized, msg.H7005, p.logg, r)

		return 0, 0, fmt.Errorf("%s", msg.H7005)
	}

	return performerId, performerRole, nil
}
