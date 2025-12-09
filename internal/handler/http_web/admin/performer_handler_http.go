package admin

import (
	"FGW_WEB/internal/config"
	"FGW_WEB/internal/handler"
	"FGW_WEB/internal/handler/http_err"
	"FGW_WEB/internal/handler/json_api"
	"FGW_WEB/internal/handler/json_err"
	"FGW_WEB/internal/model"
	"FGW_WEB/internal/service"
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"FGW_WEB/pkg/convert"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
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
	mux.HandleFunc("/admin/performers/upd", p.authMiddleware.RequireAuth(p.authMiddleware.RequireRole([]int{3}, p.handleJSONUpdate)))
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

	if performerIdStr := r.URL.Query().Get("performerId"); performerIdStr != "" {
		p.markEditingPerformer(performerIdStr, performers)
	}

	data := struct {
		Title         string
		CurrentPage   string
		Performers    []*model.Performer
		Roles         []*model.Role
		PerformerFIO  string
		PerformerId   int
		PerformerRole string
	}{
		Title:         "Список сотрудников",
		CurrentPage:   "performers",
		Performers:    performers,
		Roles:         roles,
		PerformerFIO:  performer.FIO,
		PerformerId:   performerId,
		PerformerRole: role.Name,
	}

	p.renderPages(w, "admin.html", data, r, tmplAdminPerformersHTML)
}

//func (p *PerformerHandlerHTML) UpdPerformerHTML(w http.ResponseWriter, r *http.Request) {
//	switch r.Method {
//	case http.MethodPost:
//		p.processUpdFormPerformer(w, r)
//	default:
//		http_err.SendErrorHTTP(w, http.StatusMethodNotAllowed, "", p.logg, r)
//	}
//}

func (p *PerformerHandlerHTML) UpdatePerformer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodPut {
		json_err.SendErrorResponse(w, http.StatusMethodNotAllowed, msg.H7000, "", r)

		return
	}

	performerIdStr := r.URL.Query().Get("performerId")
	performerId := convert.ConvStrToInt(performerIdStr)

	var performer model.Performer
	if err := json.NewDecoder(r.Body).Decode(&performer); err != nil {
		json_err.SendErrorResponse(w, http.StatusBadRequest, msg.H7004, err.Error(), r)

		return
	}

	exists, err := p.performerService.ExistPerformer(r.Context(), performerId)
	if err != nil {
		json_err.SendErrorResponse(w, http.StatusInternalServerError, msg.H7001, err.Error(), r)

		return
	}

	if !exists {
		json_err.SendErrorResponse(w, http.StatusNotFound, msg.H7008, "", r)

		return
	}

	if err = p.performerService.UpdPerformer(r.Context(), performerId, &performer); err != nil {
		json_err.SendErrorResponse(w, http.StatusInternalServerError, msg.H7001, err.Error(), r)

		return
	}

	response := model.PerformerUpdate{
		Success: true,
		Message: "Сотрудник успешно обновлен",
	}

	w.WriteHeader(http.StatusOK)
	json_api.WriteJSON(w, response, r)
}

// Обработчик для JSON запросов от Fetch API
func (p *PerformerHandlerHTML) handleJSONUpdate(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем правильный Content-Type для JSON ответа
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Декодируем JSON запрос
	var req struct {
		PerformerId  int `json:"performerId"`
		IdRoleAForms int `json:"idRoleAForms"`
		IdRoleAFGW   int `json:"idRoleAFGW"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		p.logg.LogE("Ошибка декодирования JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Неверный формат JSON: " + err.Error(),
		})
		return
	}

	// Проверяем существование сотрудника
	exists, err := p.performerService.ExistPerformer(r.Context(), req.PerformerId)
	if err != nil {
		p.logg.LogE("Ошибка проверки сотрудника:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Ошибка проверки сотрудника: " + err.Error(),
		})
		return
	}

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Сотрудник не найден",
		})
		return
	}

	// Получаем ID текущего пользователя из сессии
	authPerformerId := 0
	if session, err := config.Store.Get(r, config.GetSessionName()); err == nil {
		if id, ok := session.Values[config.SessionPerformerKey].(int); ok {
			authPerformerId = id
		}
	}

	// Подготавливаем данные для обновления
	performer := model.Performer{
		Id:           req.PerformerId,
		IdRoleAForms: req.IdRoleAForms,
		IdRoleAFGW:   req.IdRoleAFGW,
		AuditRec: model.Audit{
			UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
			UpdatedBy: authPerformerId,
		},
	}

	// Выполняем обновление
	if err = p.performerService.UpdPerformer(r.Context(), req.PerformerId, &performer); err != nil {
		p.logg.LogE("Ошибка обновления сотрудника:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Ошибка обновления: " + err.Error(),
		})
		return
	}

	// Успешный ответ
	response := map[string]interface{}{
		"success":     true,
		"message":     "Роли успешно обновлены",
		"performerId": req.PerformerId,
		"updatedAt":   time.Now().Format("02.01.2006 15:04:05"),
		"updatedBy":   authPerformerId,
	}

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		p.logg.LogE("Ошибка кодирования JSON ответа:", err)
	}
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
