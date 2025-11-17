package http_err

import (
	"FGW_WEB/pkg/common"
	"net/http"
)

func WriteServerError(w http.ResponseWriter, r *http.Request, logg *common.Logger, message, err string) {
	http.Error(w, message, http.StatusInternalServerError)
	logg.LogWithResponseE(message+err, http.StatusInternalServerError, r.Method, r.URL.Path)
}

func WriteMethodNotAllowed(w http.ResponseWriter, r *http.Request, logg *common.Logger, message, err string) {
	http.Error(w, message, http.StatusMethodNotAllowed)
	logg.LogWithResponseE(message+err, http.StatusMethodNotAllowed, r.Method, r.URL.Path)
}
