package json_api

import (
	"FGW_WEB/internal/handler/json_err"
	"FGW_WEB/internal/model"
	"FGW_WEB/internal/service"
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"FGW_WEB/pkg/convert"
	"encoding/json"
	"net/http"
)

type RoleHandlerJSON struct {
	roleService service.RoleUseCase
	logg        *common.Logger
}

func NewRoleHandlerJSON(roleService service.RoleUseCase, logger *common.Logger) *RoleHandlerJSON {
	return &RoleHandlerJSON{roleService: roleService, logg: logger}
}

func (r *RoleHandlerJSON) ServerHTTPJSONRouter(mux *http.ServeMux) {
	mux.HandleFunc("/api/fgw/roles", r.AllRoleJSON)
	mux.HandleFunc("/api/fgw/roles/add", r.AddRoleJSON)
	mux.HandleFunc("/api/fgw/roles/upd", r.UpdRoleJSON)
}

func (r *RoleHandlerJSON) AllRoleJSON(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if req.Method != http.MethodGet {
		json_err.SendErrorResponse(w, http.StatusMethodNotAllowed, msg.H7000, "", req)

		return
	}

	roles, err := r.roleService.GetAllRole(req.Context())
	if err != nil {
		json_err.SendErrorResponse(w, http.StatusInternalServerError, msg.H7001, err.Error(), req)

		return
	}

	if len(roles) == 0 {
		w.WriteHeader(http.StatusNoContent)
		if err = json.NewEncoder(w).Encode(&model.RoleList{Roles: []*model.Role{}}); err != nil {
			json_err.SendErrorResponse(w, http.StatusNoContent, msg.H7009, err.Error(), req)

			return
		}
	}

	data := model.RoleList{Roles: roles}

	WriteJSON(w, &data, req)
}

func (r *RoleHandlerJSON) AddRoleJSON(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if req.Method != http.MethodPost {
		json_err.SendErrorResponse(w, http.StatusMethodNotAllowed, msg.H7000, "", req)

		return
	}

	var role model.Role
	if err := json.NewDecoder(req.Body).Decode(&role); err != nil {
		json_err.SendErrorResponse(w, http.StatusBadRequest, msg.H7004, err.Error(), req)

		return
	}

	if err := r.roleService.AddRole(req.Context(), &role); err != nil {
		json_err.SendErrorResponse(w, http.StatusInternalServerError, msg.H7001, err.Error(), req)

		return
	}

	WriteJSON(w, role, req)
}

func (r *RoleHandlerJSON) UpdRoleJSON(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if req.Method != http.MethodPut {
		json_err.SendErrorResponse(w, http.StatusMethodNotAllowed, msg.H7000, "", req)

		return
	}

	roleIdStr := req.URL.Query().Get("roleId")
	roleId := convert.ConvStrToInt(roleIdStr)

	var role model.Role
	if err := json.NewDecoder(req.Body).Decode(&role); err != nil {
		json_err.SendErrorResponse(w, http.StatusBadRequest, msg.H7004, err.Error(), req)

		return
	}

	exists, err := r.roleService.ExistRole(req.Context(), roleId)
	if err != nil {
		json_err.SendErrorResponse(w, http.StatusInternalServerError, msg.H7001, err.Error(), req)

		return
	}

	if !exists {
		json_err.SendErrorResponse(w, http.StatusNotFound, msg.H7008, "", req)

		return
	}

	if err = r.roleService.UpdRole(req.Context(), roleId, &role); err != nil {
		json_err.SendErrorResponse(w, http.StatusInternalServerError, msg.H7001, err.Error(), req)

		return
	}

	response := model.RoleUpdate{
		Success: true,
		Message: "Роль успешно обновлена",
	}

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}
