package json_err

import (
	"FGW_WEB/pkg/common"
	"encoding/json"
	"net/http"
	"runtime"
	"strings"
)

const (
	// SkipNumOfStackFrame количество кадров стека, которые необходимо пропустить перед записью на ПК, где 0 идентифицирует
	// кадр для самих вызывающих абонентов, а 1 идентифицирует вызывающего абонента. Возвращает количество записей,
	// записанных на компьютер.
	SkipNumOfStackFrame   = 3
	DefaultMaxStackFrames = 15
)

func SendErrorResponse(w http.ResponseWriter, statusCode int, msgCode string, r *http.Request) {
	funcName, fileName, lineNumber, filePath := fileWithFuncAndLineNum()

	errorResponse := struct {
		Error       string               `json:"error"`
		Code        int                  `json:"code"`
		Description common.ResponseEntry `json:"description"`
		Detail      common.DetailEntry   `json:"detail"`
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
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

// fileWithFuncAndLineNum возвращает имя функции, имя файла, номер строки, путь файла.
func fileWithFuncAndLineNum() (string, string, int, string) {
	pc := make([]uintptr, DefaultMaxStackFrames)
	frameCount := runtime.Callers(SkipNumOfStackFrame, pc)
	if frameCount == 0 {
		return "неизвестно", "неизвестно", 0, ""
	}

	frames := runtime.CallersFrames(pc[:frameCount])
	frame, ok := frames.Next()
	if !ok {
		return "неизвестно", "неизвестно", 0, ""
	}

	idxFile := strings.LastIndexByte(frame.File, '/')

	return frame.Function, frame.File[idxFile+1:], frame.Line, frame.File
}
