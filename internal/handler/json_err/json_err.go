package json_err

import (
	"FGW_WEB/pkg/common"
	"encoding/json"
	"net/http"
)

const (
	// SkipNumOfStackFrame количество кадров стека, которые необходимо пропустить перед записью на ПК, где 0 идентифицирует
	// кадр для самих вызывающих абонентов, а 1 идентифицирует вызывающего абонента. Возвращает количество записей,
	// записанных на компьютер.
	skipNumOfStackFrame = 3
)

func SendErrorResponse(w http.ResponseWriter, statusCode int, msgCode string, msg string, r *http.Request) {
	funcName, fileName, lineNumber, filePath := common.FileWithFuncAndLineNum(skipNumOfStackFrame)

	errorResponse := struct {
		Error       string               `json:"error"`
		Code        int                  `json:"code"`
		Description common.ResponseEntry `json:"description"`
		Detail      common.DetailEntry   `json:"detail"`
		Message     string               `json:"message"`
	}{
		msgCode,
		statusCode,
		common.ResponseEntry{
			StatusCode: statusCode,
			MethodHTTP: r.Method,
			URL:        r.URL.Path,
		},
		common.DetailEntry{
			FunctionName: funcName,
			FileName:     fileName,
			LineNumber:   lineNumber,
			PathToFile:   filePath,
		},
		msg,
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
