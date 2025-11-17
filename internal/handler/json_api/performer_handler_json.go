package json_api

import (
	"FGW_WEB/internal/service"
	"FGW_WEB/internal/service/dto"
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
	mux.HandleFunc("/api/fgw/login", p.AuthPerformer)
	mux.HandleFunc("/api/fgw/performers", p.AllPerformersJSON)
}

func (p *PerformerHandlerJSON) AllPerformersJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodGet {
		p.sendErrorResponse(w, http.StatusMethodNotAllowed, msg.H7000, r)

		return
	}

	performers, err := p.performerService.GetAllPerformers(r.Context())
	if err != nil {
		p.sendErrorResponse(w, http.StatusInternalServerError, msg.H7001, r)

		return
	}

	if len(performers) == 0 {
		w.WriteHeader(http.StatusNoContent)
		if err = json.NewEncoder(w).Encode(&dto.PerformerDTOList{Performers: []dto.PerformerDTO{}}); err != nil {
			return
		}
	}

	data := dto.PerformerDTOList{Performers: performers}

	if err = json.NewEncoder(w).Encode(data); err != nil {
		p.sendErrorResponse(w, http.StatusInternalServerError, msg.H7001, r)

		return
	}
}

func (p *PerformerHandlerJSON) AuthPerformer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID       int    `json:"id"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := p.performerService.AuthPerformer(r.Context(), req.ID, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (p *PerformerHandlerJSON) sendErrorResponse(w http.ResponseWriter, statusCode int, msgCode string, r *http.Request) {
	errorResponse := struct {
		Error       string               `json:"error"`
		Code        int                  `json:"code"`
		Description common.ResponseEntry `json:"description"`
	}{
		msgCode,
		statusCode,
		common.ResponseEntry{
			StatusCode: statusCode,
			MethodHTTP: r.Method,
			URL:        r.URL.Path,
		},
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.logg.LogWithResponseE(msgCode, statusCode, r.Method, r.URL.Path)
}
