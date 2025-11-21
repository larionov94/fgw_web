package http_err

import (
	"FGW_WEB/pkg/common"
	"net/http"
)

func WriteServerError(w http.ResponseWriter, r *http.Request, logg *common.Logger, message, err string) {
	http.Error(w, message, http.StatusInternalServerError)
	logg.LogHttpErr(message+err, http.StatusInternalServerError, r.Method, r.URL.Path)
}

func WriteMethodNotAllowed(w http.ResponseWriter, r *http.Request, logg *common.Logger, message, err string) {
	http.Error(w, message, http.StatusMethodNotAllowed)
	logg.LogHttpErr(message+err, http.StatusMethodNotAllowed, r.Method, r.URL.Path)
}

func WriteUnauthorized(w http.ResponseWriter, r *http.Request, logg *common.Logger, message, err string) {
	http.Error(w, message, http.StatusUnauthorized)
	logg.LogHttpErr(message+err, http.StatusUnauthorized, r.Method, r.URL.Path)
}

func WriteBadRequest(w http.ResponseWriter, r *http.Request, logg *common.Logger, message, err string) {
	http.Error(w, message, http.StatusBadRequest)
	logg.LogHttpErr(message+err, http.StatusBadRequest, r.Method, r.URL.Path)
}

func WriteForbidden(w http.ResponseWriter, r *http.Request, logg *common.Logger, message, err string) {
	http.Error(w, message, http.StatusForbidden)
	logg.LogHttpErr(message+err, http.StatusForbidden, r.Method, r.URL.Path)
}
