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
	"fmt"
	"html/template"
	"net/http"
	"time"
)

const (
	tmplAdminRolesHTML = "roles.html"
)

type RoleHandlerHTML struct {
	roleService      service.RoleUseCase
	performerService service.PerformerUseCase
	logg             *common.Logger
	authMiddleware   *handler.AuthMiddleware
}

func NewRoleHandlerHTML(roleService service.RoleUseCase, logger *common.Logger, authMiddleware *handler.AuthMiddleware, performerService service.PerformerUseCase) *RoleHandlerHTML {
	return &RoleHandlerHTML{roleService: roleService, logg: logger, authMiddleware: authMiddleware, performerService: performerService}
}

func (r *RoleHandlerHTML) ServerHTTPHTMLRouter(mux *http.ServeMux) {
	mux.HandleFunc("/admin/roles", r.authMiddleware.RequireAuth(r.authMiddleware.RequireRole([]int{3}, r.AllRoleHTML)))
	//mux.HandleFunc("/admin/roles/add", r.AddRoleHTML)
	mux.HandleFunc("/admin/roles/upd", r.authMiddleware.RequireAuth(r.authMiddleware.RequireRole([]int{3}, r.HandleJSONUpdate)))
}

func (r *RoleHandlerHTML) AllRoleHTML(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	performerId, performerRoleId, err := r.getSessionPerformerData(w, req)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusUnauthorized, err.Error(), r.logg, req)

		return
	}

	if req.Method != http.MethodGet {
		http_err.SendErrorHTTP(w, http.StatusMethodNotAllowed, "", r.logg, req)

		return
	}

	roles, err := r.roleService.GetAllRole(req.Context())
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusInternalServerError, err.Error(), r.logg, req)

		return
	}

	role, err := r.roleService.FindRoleById(req.Context(), performerRoleId)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusNotFound, err.Error(), r.logg, req)

		return
	}

	performer, err := r.performerService.FindByIdPerformer(req.Context(), performerId)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusNotFound, err.Error(), r.logg, req)

		return
	}

	data := struct {
		Title         string
		CurrentPage   string
		Roles         []*model.Role
		PerformerId   int
		PerformerRole string
		PerformerFIO  string
	}{
		Title:         "Список ролей",
		CurrentPage:   "roles",
		Roles:         roles,
		PerformerId:   performerId,
		PerformerRole: role.Name,
		PerformerFIO:  performer.FIO,
	}

	r.renderPages(w, tmplAdminHTML, data, req, tmplAdminRolesHTML, tmplAdminPerformersHTML)
}

func (r *RoleHandlerHTML) AddRoleHTML(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if req.Method != http.MethodPost {
		http_err.SendErrorHTTP(w, http.StatusMethodNotAllowed, "", r.logg, req)

		return
	}

	roleIdStr := req.FormValue("roleId")
	nameStr := req.FormValue("name")
	descStr := req.FormValue("description")
	createdByStr := req.FormValue("createdBy")

	if roleIdStr == "" || nameStr == "" || descStr == "" || createdByStr == "" {
		r.renderErrorPage(w, http.StatusBadRequest, msg.E3214, req)

		return
	}

	roleId := convert.ConvStrToInt(roleIdStr)
	createdBy := convert.ConvStrToInt(createdByStr)

	role := &model.Role{
		Id:   roleId,
		Name: nameStr,
		Desc: descStr,
		AuditRec: model.Audit{
			CreatedBy: createdBy,
		},
	}

	if err := r.roleService.AddRole(req.Context(), role); err != nil {
		http_err.SendErrorHTTP(w, http.StatusInternalServerError, err.Error(), r.logg, req)

		return
	}

	http.Redirect(w, req, "/admin/roles", http.StatusSeeOther)
}

func (r *RoleHandlerHTML) HandleJSONUpdate(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var reqs struct {
		RoleId      int    `json:"roleId"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(req.Body).Decode(&reqs); err != nil {
		json_err.SendErrorResponse(w, http.StatusBadRequest, msg.H7004, err.Error(), req)

		return
	}

	exists, err := r.roleService.ExistRole(req.Context(), reqs.RoleId)
	if err != nil {
		json_err.SendErrorResponse(w, http.StatusInternalServerError, msg.H7001, err.Error(), req)

		return
	}

	if !exists {
		json_err.SendErrorResponse(w, http.StatusNotFound, msg.H7008, "", req)

		return
	}

	if session, err := config.Store.Get(req, config.GetSessionName()); err == nil {
		if id, ok := session.Values[config.SessionPerformerKey].(int); ok {
			authPerformerId = id
		}
	}

	role := model.Role{
		Id:   reqs.RoleId,
		Name: reqs.Name,
		Desc: reqs.Description,
		AuditRec: model.Audit{
			UpdatedBy: authPerformerId,
		},
	}

	if err = r.roleService.UpdRole(req.Context(), reqs.RoleId, &role); err != nil {
		json_err.SendErrorResponse(w, http.StatusInternalServerError, msg.H7001, err.Error(), req)

		return
	}

	response := map[string]interface{}{
		"success":   true,
		"message":   "Роль успешна обновлена",
		"roleId":    reqs.RoleId,
		"updatedAt": time.Now().Format("2006-01-02 15:04:05"),
		"updatedBy": authPerformerId,
	}

	w.WriteHeader(http.StatusOK)
	json_api.WriteJSON(w, response, req)
}

func (r *RoleHandlerHTML) renderErrorPage(w http.ResponseWriter, statusCode int, msgCode string, req *http.Request) {
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
		Method:     req.Method,
		Path:       req.URL.Path,
	}

	w.WriteHeader(statusCode)
	r.logg.LogHttpErr(msgCode, statusCode, req.Method, req.URL.Path)
	r.renderPage(w, tmplErrorHTML, data, req)
}

func (r *RoleHandlerHTML) renderPage(w http.ResponseWriter, tmpl string, data interface{}, req *http.Request) {
	parseTmpl, err := template.New(tmpl).Funcs(
		template.FuncMap{
			"formatDateTime": convert.FormatDateTime,
		}).ParseFiles(prefixTmplAdmin + tmpl)
	if err != nil {
		r.renderErrorPage(w, http.StatusInternalServerError, msg.H7002+err.Error(), req)

		return
	}

	if err = parseTmpl.ExecuteTemplate(w, tmpl, data); err != nil {
		r.renderErrorPage(w, http.StatusInternalServerError, msg.H7003+err.Error(), req)

		return
	}
}

func (r *RoleHandlerHTML) renderPages(
	w http.ResponseWriter, tmpl string, data interface{}, req *http.Request, addTemplates ...string) {

	templatePaths := []string{prefixDefaultTmpl + tmpl}

	for _, addTmpl := range addTemplates {
		templatePaths = append(templatePaths, prefixAdminTmpl+addTmpl)
	}

	parseTmpl, err := template.New(tmpl).Funcs(template.FuncMap{
		"formatDateTime": convert.FormatDateTime,
		"add":            func(a, b int) int { return a + b },
		"sub":            func(a, b int) int { return a - b },
	}).ParseFiles(templatePaths...)

	if err != nil {
		r.renderErrorPage(w, http.StatusInternalServerError, msg.H7002+err.Error(), req)

		return
	}

	if err = parseTmpl.ExecuteTemplate(w, tmpl, data); err != nil {
		r.renderErrorPage(w, http.StatusInternalServerError, msg.H7003+err.Error(), req)

		return
	}
}

// getSessionPerformerData получить данные о сеансе сотрудника.
func (r *RoleHandlerHTML) getSessionPerformerData(w http.ResponseWriter, req *http.Request) (int, int, error) {
	performerId, ok := r.authMiddleware.GetPerformerId(req)
	if !ok {
		http_err.SendErrorHTTP(w, http.StatusUnauthorized, msg.H7005, r.logg, req)

		return 0, 0, fmt.Errorf("%s", msg.H7005)
	}
	authPerformerId = performerId

	performerRole, ok := r.authMiddleware.GetRoleId(req)
	if !ok {
		http_err.SendErrorHTTP(w, http.StatusUnauthorized, msg.H7005, r.logg, req)

		return 0, 0, fmt.Errorf("%s", msg.H7005)
	}

	return performerId, performerRole, nil
}
