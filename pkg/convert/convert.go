package convert

import (
	"FGW_WEB/pkg/common"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const SkipNumOfStackFrame = 3

func pathToStrCode() string {
	var funcName, fileName, lineNumber, filePath = common.FileWithFuncAndLineNum(SkipNumOfStackFrame)

	return fmt.Sprintf("%s -> %s -> %d -> %s", funcName, fileName, lineNumber, filePath)
}

// ConvStrToInt конвертировать строку в число.
func ConvStrToInt(str string) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("Ошибка: [%s] --- ссылка на код: [ %s ] --- значение: [%v]", err.Error(), pathToStrCode(), value)

		return 0
	}

	return value
}

// ParseStrToID пытается разобрать строку в UUID и пишет ошибку в HTTP-ответ при неудаче.
func ParseStrToID(fieldName string) (int, error) {
	value, err := strconv.Atoi(fieldName)
	if err != nil {
		log.Printf("Ошибка: [%s] --- ссылка на код: [ %s ] --- значение: [%v]", err.Error(), pathToStrCode(), value)

		return 0, err
	}

	return value, nil
}

// ParseFormFieldInt преобразует поле в целое число, полученное из HTTP запроса.
func ParseFormFieldInt(r *http.Request, fieldName string) int {

	formValue := r.FormValue(fieldName)
	if formValue == "" {
		formValue = "0"

		return 0
	}
	value, err := strconv.Atoi(formValue)
	if err != nil {
		log.Printf("Ошибка: [%s] --- ссылка на код: [ %s ] --- поле: [%s] --- значение: [%v]", err.Error(), pathToStrCode(), fieldName, value)

		return 0
	}

	return value
}
