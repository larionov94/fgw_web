package http_err

import (
	"FGW_WEB/pkg/common"
	"fmt"
	"net/http"
)

const (
	// SkipNumOfStackFrame количество кадров стека, которые необходимо пропустить перед записью на ПК, где 0 идентифицирует
	// кадр для самих вызывающих абонентов, а 1 идентифицирует вызывающего абонента. Возвращает количество записей,
	// записанных на компьютер.
	skipNumOfStackFrame = 3
)

func SendErrorHTTP(w http.ResponseWriter, statusCode int, msgErr string, logg *common.Logger, r *http.Request) {
	_, fileName, lineCode, _ := common.FileWithFuncAndLineNum(skipNumOfStackFrame)

	result := fmt.Sprintf("H7777 %s --- %s:%d", msgErr, fileName, lineCode)

	logg.LogHttpErr(result, statusCode, r.Method, r.URL.Path)
	http.Error(w, fmt.Sprintf(" %s", msgErr), statusCode)

}
