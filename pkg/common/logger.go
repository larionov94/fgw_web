package common

import (
	"FGW_WEB/pkg/common/msg"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	// skipNumOfStackFrame количество кадров стека, которые необходимо пропустить перед записью на ПК, где 0 идентифицирует
	// кадр для самих вызывающих абонентов, а 1 идентифицирует вызывающего абонента. Возвращает количество записей,
	// записанных на компьютер.
	skipNumOfStackFrame    = 5
	CodeLength             = 6 // CodeLength извлечение подстроки из поля.
	DefaultMaxStackFrames  = 15
	DefaultFilePermissions = 0644
	DefaultPathToLog       = "logCustom.json"
)

type LogLevel string

const (
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
)

type MessageEntry struct {
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Error   *string `json:"error,omitempty"`
}

type ResponseEntry struct {
	StatusCode int    `json:"statusCode"`
	MethodHTTP string `json:"methodHTTP"`
	URL        string `json:"url"`
}

type DetailEntry struct {
	FunctionName string `json:"functionName"`
	FileName     string `json:"fileName"`
	LineNumber   int    `json:"lineNumber"`
	PathToFile   string `json:"pathToFile"`
}

type LogEntry struct {
	DateTime        string         `json:"dateTime"`
	InfoPC          *InfoPC        `json:"infoPC"`
	Level           LogLevel       `json:"level"`
	Message         MessageEntry   `json:"message"`
	ResponseMessage *ResponseEntry `json:"responseMessage,omitempty"`
	Detail          *DetailEntry   `json:"detail"`
}

type Logger struct {
	file     *os.File
	infoPC   *InfoPC
	filePath string
}

// NewLogger возвращает новый объект лога.
func NewLogger(filePath string) (*Logger, error) {
	if filePath == "" {
		filePath = DefaultPathToLog
	}
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, DefaultFilePermissions)
	if err != nil {
		return nil, err
	}

	infoPC, err := NewInfoPC()
	if err != nil {
		return nil, err
	}

	return &Logger{
		file: file,
		infoPC: &InfoPC{
			Domain: infoPC.HostName(),
			IPAddr: infoPC.AddrIP(),
		},
		filePath: filePath,
	}, nil
}

func (l *Logger) LogI(msg string) {
	l.logCustom(LogLevelInfo, msg, nil, nil)
}

func (l *Logger) LogW(msg string) {
	l.logCustom(LogLevelWarn, msg, nil, nil)
}

func (l *Logger) LogE(msg string, err error) {
	var errStr *string
	if err != nil {
		errMsg := err.Error()
		errStr = &errMsg
	}
	l.logCustom(LogLevelError, msg, errStr, nil)
}

func (l *Logger) LogHttpI(msg string, statusCode int, method, url string) {
	response := &ResponseEntry{
		StatusCode: statusCode,
		MethodHTTP: method,
		URL:        url,
	}
	l.logCustom(LogLevelInfo, msg, nil, response)
}

func (l *Logger) LogHttpErr(msg string, statusCode int, method, url string) {
	response := &ResponseEntry{
		StatusCode: statusCode,
		MethodHTTP: method,
		URL:        url,
	}
	l.logCustom(LogLevelError, msg, nil, response)
}

// logCustom логирование пользовательских сообщений с поддержкой уровня логирования и дополнительной информацией.
func (l *Logger) logCustom(level LogLevel, message string, errStr *string, response *ResponseEntry) {
	entry := &LogEntry{
		DateTime:        time.Now().Format(time.DateTime),
		InfoPC:          l.infoPC,
		Level:           level,
		Message:         l.createMessage(message, errStr),
		ResponseMessage: response,
		Detail:          l.createDetails(skipNumOfStackFrame),
	}

	if err := l.writeEntry(entry); err != nil {
		log.Printf("%s: %v", msg.E3001, err)
	}
}

// createMessage создает и заполняет структуру сообщения.
func (l *Logger) createMessage(msg string, errStr *string) MessageEntry {
	code, message := splitCodeMessage(msg)
	return MessageEntry{
		Code:    code,
		Message: message,
		Error:   errStr,
	}
}

// createDetails создает и заполняет информацию о месте вызова.
func (l *Logger) createDetails(skipNumOfStackFrame int) *DetailEntry {
	funcName, fileName, lineNumber, filePath := FileWithFuncAndLineNum(skipNumOfStackFrame)
	return &DetailEntry{
		FunctionName: funcName,
		FileName:     fileName,
		LineNumber:   lineNumber,
		PathToFile:   filePath,
	}
}

// splitCodeMessage разбивает строку на код и сообщение.
// Например, "E1001: Ошибка при выполнении запроса" -> code = "E1001", msg = "Ошибка при выполнении запроса".
func splitCodeMessage(msg string) (string, string) {
	if len(msg) < CodeLength || msg == "" {
		return "", msg
	}

	return msg[:CodeLength-1], msg[CodeLength:]
}

// FileWithFuncAndLineNum возвращает имя функции, имя файла, номер строки, путь файла.
func FileWithFuncAndLineNum(skipNumOfStack int) (string, string, int, string) {
	pc := make([]uintptr, DefaultMaxStackFrames)
	frameCount := runtime.Callers(skipNumOfStack, pc)
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

// writeEntry запись в файл.
func (l *Logger) writeEntry(entry *LogEntry) error {
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("%s: %w", msg.E3002, err)
	}

	data = append(data, ',', '\n')

	if _, err := l.file.Write(data); err != nil {
		return fmt.Errorf("%s: %w", msg.E3001, err)
	}
	fmt.Println(string(data))

	return nil
}

// Close закрывает файл.
func (l *Logger) Close() {
	if l.file != nil {
		if err := l.file.Close(); err != nil {
			log.Printf("%s: %v", msg.E3000, err)
		}
		l.file = nil
	}

	log.Printf("%s", msg.I2000)
}
