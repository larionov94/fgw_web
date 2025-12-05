package admin

import (
	"FGW_WEB/internal/handler"
	"FGW_WEB/internal/handler/http_err"
	"FGW_WEB/internal/model"
	"FGW_WEB/internal/service"
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"FGW_WEB/pkg/convert"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

const (
	tmplAdminPerformersHTML = "performers.html"
	tmplErrorHTML           = "error.html"
	prefixTmplAdmin         = "web/html/admin/"
	urlAdminPerformers      = "/admin/performers"

	prefixDefaultTmpl = "web/html/"
	prefixAdminTmpl   = "web/html/admin/"
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
	mux.HandleFunc("/admin/performers/upd", p.authMiddleware.RequireAuth(p.authMiddleware.RequireRole([]int{3}, p.UpdPerformerHTML)))
}

func (p *PerformerHandlerHTML) AllPerformersHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != http.MethodGet {
		http_err.SendErrorHTTP(w, http.StatusMethodNotAllowed, "", p.logg, r)

		return
	}

	// Получаем параметры пагинации
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Размер страницы (можно вынести в конфиг)
	pageSize := 10

	// Получаем общее количество
	totalCount, err := p.performerService.GetPerformersCount(r.Context())
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusInternalServerError, err.Error(), p.logg, r)
		return
	}

	// Рассчитываем смещение
	offset := (page - 1) * pageSize

	// Получаем исполнителей с пагинацией
	performers, err := p.performerService.GetPerformersWithPagination(r.Context(), offset, pageSize)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusInternalServerError, err.Error(), p.logg, r)
		return
	}

	roles, err := p.roleService.GetAllRole(r.Context())
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusInternalServerError, err.Error(), p.logg, r)
		return
	}

	performerId, _ := p.authMiddleware.GetPerformerId(r)
	authPerformerId = performerId

	performerRole, _ := p.authMiddleware.GetRoleId(r)

	performer, err := p.performerService.FindByIdPerformer(r.Context(), performerId)
	if err != nil {
		log.Println(err.Error())
	}

	role, err := p.roleService.FindRoleById(r.Context(), performerRole)
	if err != nil {
		log.Println(err.Error())
	}

	// Получаем ID для редактирования
	if performerIdStr := r.URL.Query().Get("performerId"); performerIdStr != "" {
		p.markEditingPerformer(performerIdStr, performers)
	}

	// Рассчитываем пагинацию
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}

	// Ограничиваем номер страницы
	if page > totalPages {
		page = totalPages
	}

	// Генерируем диапазон страниц для отображения
	pages := generatePageRange(page, totalPages, 5)

	// Рассчитываем отображаемый диапазон элементов
	startItem := offset + 1
	endItem := offset + len(performers)
	if endItem > totalCount {
		endItem = totalCount
	}

	data := struct {
		Title         string
		CurrentPage   string
		Performers    []*model.Performer
		Roles         []*model.Role
		PerformerFIO  string
		PerformerId   int
		PerformerRole string

		// Пагинация
		Page           int
		PageSize       int
		TotalCount     int
		TotalPages     int
		Pages          []int
		StartItem      int
		EndItem        int
		PerformerIdStr int // Для сохранения в пагинации
	}{
		Title:         "Список сотрудников",
		CurrentPage:   "performers",
		Performers:    performers,
		Roles:         roles,
		PerformerFIO:  performer.FIO,
		PerformerId:   performerId,
		PerformerRole: role.Name,

		// Пагинация
		Page:           page,
		PageSize:       pageSize,
		TotalCount:     totalCount,
		TotalPages:     totalPages,
		Pages:          pages,
		StartItem:      startItem,
		EndItem:        endItem,
		PerformerIdStr: performerId,
	}

	p.renderPages(w, "admin.html", data, r, tmplAdminPerformersHTML)
}

// Вспомогательная функция для генерации диапазона страниц
func generatePageRange(current, total, maxPages int) []int {
	var pages []int

	if total <= maxPages {
		// Если страниц меньше или равно maxPages, показываем все
		for i := 1; i <= total; i++ {
			pages = append(pages, i)
		}
	} else {
		// Определяем начальную и конечную страницу
		start := current - maxPages/2
		end := current + maxPages/2

		if start < 1 {
			start = 1
			end = maxPages
		}

		if end > total {
			end = total
			start = total - maxPages + 1
		}

		for i := start; i <= end; i++ {
			pages = append(pages, i)
		}
	}

	return pages
}

func (p *PerformerHandlerHTML) UpdPerformerHTML(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p.processUpdFormPerformer(w, r)
	default:
		http_err.SendErrorHTTP(w, http.StatusMethodNotAllowed, "", p.logg, r)
	}
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

	if idRoleAFormsStr == "" || idRoleAFGWStr == "" {
		p.renderErrorPage(w, http.StatusUnauthorized, msg.E3214, r)

		return
	}

	performerId := convert.ConvStrToInt(performerIdStr)
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
			UpdatedBy: authPerformerId,
		},
	}

	if err = p.performerService.UpdPerformer(r.Context(), performerId, &performer); err != nil {
		http_err.SendErrorHTTP(w, http.StatusInternalServerError, err.Error(), p.logg, r)

		return
	}

	http.Redirect(w, r, urlAdminPerformers, http.StatusSeeOther)
}

func (p *PerformerHandlerHTML) markEditingPerformer(id string, performers []*model.Performer) {
	performerId := convert.ConvStrToInt(id)
	for _, performer := range performers {
		if performer.Id == performerId {
			performer.IsEditing = true
		}
	}
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
			"add":            func(a, b int) int { return a + b },
			"sub":            func(a, b int) int { return a - b },
			"formatDateTime": convert.FormatDateTime,
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
			"add":            func(a, b int) int { return a + b },
			"sub":            func(a, b int) int { return a - b },
			"formatDateTime": convert.FormatDateTime,
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
