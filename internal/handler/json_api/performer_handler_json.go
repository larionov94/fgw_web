package json_api

import (
	"FGW_WEB/internal/service"
	"FGW_WEB/pkg/common"
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
	mux.HandleFunc("/api/fgw/performers", p.AllPerformers)
}

func (p *PerformerHandlerJSON) AllPerformers(w http.ResponseWriter, r *http.Request) {
	performers, err := p.performerService.GetAllPerformers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(performers)
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
