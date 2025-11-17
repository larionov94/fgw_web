package json_api

import (
	"FGW_WEB/internal/handler/json_err"
	"FGW_WEB/internal/model"
	"FGW_WEB/internal/service"
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"encoding/json"
	"net/http"
)

type PerformerHandlerJSON struct {
	performerService service.PerformerUseCase
	logg             *common.Logger
}

func NewPerformerHandlerJSON(performerService service.PerformerUseCase, logg *common.Logger) *PerformerHandlerJSON {
	return &PerformerHandlerJSON{performerService: performerService, logg: logg}
}

func (p *PerformerHandlerJSON) ServeHTTPJSONRouter(mux *http.ServeMux) {
	mux.HandleFunc("/api/fgw/login", p.AuthPerformerJSON)
	mux.HandleFunc("/api/fgw/performers", p.AllPerformersJSON)
}

func (p *PerformerHandlerJSON) AllPerformersJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodGet {
		json_err.SendErrorResponse(w, http.StatusMethodNotAllowed, msg.H7000, r)

		return
	}

	performers, err := p.performerService.GetAllPerformers(r.Context())
	if err != nil {
		json_err.SendErrorResponse(w, http.StatusInternalServerError, msg.H7001, r)

		return
	}

	if len(performers) == 0 {
		w.WriteHeader(http.StatusNoContent)
		if err = json.NewEncoder(w).Encode(&model.PerformerList{Performers: []model.Performer{}}); err != nil {
			return
		}
	}

	data := model.PerformerList{Performers: performers}

	if err = json.NewEncoder(w).Encode(&data); err != nil {
		json_err.SendErrorResponse(w, http.StatusInternalServerError, msg.H7001, r)

		return
	}
}

func (p *PerformerHandlerJSON) AuthPerformerJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodPost {
		json_err.SendErrorResponse(w, http.StatusMethodNotAllowed, msg.H7000, r)

		return
	}

	var req struct {
		Id       int    `json:"id"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json_err.SendErrorResponse(w, http.StatusBadRequest, msg.H7004, r)

		return
	}

	result, err := p.performerService.AuthPerformer(r.Context(), req.Id, req.Password)
	if err != nil {
		json_err.SendErrorResponse(w, http.StatusUnauthorized, msg.H7005, r)

		return
	}

	WriteJSON(w, result, r)
}

func WriteJSON(w http.ResponseWriter, v interface{}, r *http.Request) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		json_err.SendErrorResponse(w, http.StatusInternalServerError, msg.H7006, r)

		return
	}
}
