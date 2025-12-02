package admin

import (
	"FGW_WEB/internal/handler/http_err"
	"FGW_WEB/internal/model"
	"FGW_WEB/internal/service"
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"FGW_WEB/pkg/convert"
	"html/template"
	"net/http"
)

type RoleHandlerHTML struct {
	roleService service.RoleUseCase
	logg        *common.Logger
}

func NewRoleHandlerHTML(roleService service.RoleUseCase, logger *common.Logger) *RoleHandlerHTML {
	return &RoleHandlerHTML{roleService: roleService, logg: logger}
}

func (r *RoleHandlerHTML) ServerHTTPHTMLRouter(mux *http.ServeMux) {
	mux.HandleFunc("/fgw/roles", r.AllRoleHTML)
	mux.HandleFunc("/fgw/roles/add", r.AddRoleHTML)
	mux.HandleFunc("/fgw/roles/upd", r.UpdRoleHTML)
}

func (r *RoleHandlerHTML) AllRoleHTML(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if req.Method != http.MethodGet {
		http_err.SendErrorHTTP(w, http.StatusMethodNotAllowed, "", r.logg, req)

		return
	}

	roles, err := r.roleService.GetAllRole(req.Context())
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusInternalServerError, err.Error(), r.logg, req)

		return
	}

	data := struct {
		Title string
		Roles []*model.Role
	}{
		Title: "Список ролей",
		Roles: roles,
	}

	if roleIdStr := req.URL.Query().Get("roleId"); roleIdStr != "" {
		r.markEditingRole(roleIdStr, roles)
	}

	r.renderPage(w, "roles.html", data, req)
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

	http.Redirect(w, req, "/fgw/roles", http.StatusSeeOther)
}

func (r *RoleHandlerHTML) UpdRoleHTML(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		r.processUpdFormRole(w, req)
	case http.MethodGet:
		r.renderUpdFormRole(w, req)
	default:
		http_err.SendErrorHTTP(w, http.StatusMethodNotAllowed, "", r.logg, req)
	}
}

func (r *RoleHandlerHTML) renderUpdFormRole(w http.ResponseWriter, req *http.Request) {
	roleIdStr := req.URL.Query().Get("roleId")
	http.Redirect(w, req, "/fgw/roles?roleId="+roleIdStr, http.StatusSeeOther)
}

func (r *RoleHandlerHTML) processUpdFormRole(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := req.ParseForm(); err != nil {
		http_err.SendErrorHTTP(w, http.StatusInternalServerError, err.Error(), r.logg, req)

		return
	}

	roleIdStr := req.FormValue("roleId")
	nameStr := req.FormValue("name")
	descStr := req.FormValue("description")
	updateByStr := req.FormValue("updatedBy")

	if roleIdStr == "" || nameStr == "" || descStr == "" || updateByStr == "" {
		r.renderErrorPage(w, http.StatusBadRequest, msg.E3214, req)

		return
	}

	roleId := convert.ConvStrToInt(roleIdStr)
	updateBy := convert.ConvStrToInt(updateByStr)

	role := &model.Role{
		Id:   roleId,
		Name: nameStr,
		Desc: descStr,
		AuditRec: model.Audit{
			UpdatedBy: updateBy,
		},
		IsEditing: false,
	}

	if err := r.roleService.UpdRole(req.Context(), roleId, role); err != nil {
		http_err.SendErrorHTTP(w, http.StatusInternalServerError, err.Error(), r.logg, req)

		return
	}

	http.Redirect(w, req, "/fgw/roles", http.StatusSeeOther)
}

func (r *RoleHandlerHTML) markEditingRole(id string, roles []*model.Role) {
	roleId := convert.ConvStrToInt(id)
	for _, role := range roles {
		if role.Id == roleId {
			role.IsEditing = true
		}
	}
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
		}).ParseFiles(prefixTmplPerformers + tmpl)
	if err != nil {
		r.renderErrorPage(w, http.StatusInternalServerError, msg.H7002+err.Error(), req)

		return
	}

	if err = parseTmpl.ExecuteTemplate(w, tmpl, data); err != nil {
		r.renderErrorPage(w, http.StatusInternalServerError, msg.H7003+err.Error(), req)

		return
	}
}
