package common

import (
	"FGW_WEB/pkg/convert"
	"errors"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const pathFile = "test_log.json"

func openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
}

var validEntry = &LogEntry{
	DateTime: convert.GetCurrentDateTime(),
	InfoPC:   &InfoPC{},
	Level:    LogLevelInfo,
	Message: MessageEntry{
		Code:    "I1000",
		Message: "New entry",
	},
	ResponseMessage: nil,
	Detail: &DetailEntry{
		FunctionName: "main.TestWriteEntry",
		FileName:     "logger_test.go",
		LineNumber:   100,
		PathToFile:   "/path/to/logger_test.go",
	},
}

func TestLogger_writeEntry(t *testing.T) {
	cleanupTestFiles(t)

	file, err := openFile(pathFile)
	if err != nil {
		log.Printf("Ошибка при открытии файла: %v", err)
		return
	}

	type fields struct {
		file     *os.File
		infoPC   *InfoPC
		filePath string
	}
	type args struct {
		entry *LogEntry
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Успешно записан лог",
			fields: fields{
				file:     file,
				infoPC:   nil,
				filePath: pathFile,
			},
			args:    args{entry: validEntry},
			wantErr: false,
		},
		{
			name: "Не успешно записан лог",
			fields: fields{
				file:     nil,
				infoPC:   nil,
				filePath: pathFile,
			},
			args:    args{entry: validEntry},
			wantErr: true,
		},
		{
			name: "Не успешное маршалирование",
			fields: fields{
				file:     nil,
				infoPC:   nil,
				filePath: "json.json",
			},
			args:    args{entry: nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				file:     tt.fields.file,
				infoPC:   tt.fields.infoPC,
				filePath: tt.fields.filePath,
			}
			if err = l.writeEntry(tt.args.entry); (err != nil) != tt.wantErr {
				t.Errorf("writeEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
			l.file.Close()
		})
	}
}

func TestLogger_Close(t *testing.T) {
	cleanupTestFiles(t)

	type fields struct {
		file *os.File
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "Успешно закрыт файл",
			fields: fields{file: &os.File{}},
		},
		{
			name:   "Файл не закрыт",
			fields: fields{file: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				file: tt.fields.file,
			}
			l.file.Close()
		})
	}
}

func Test_fileWithFuncAndLineNum(t *testing.T) {
	tests := []struct {
		name  string
		want  string
		want1 string
		want2 int
		want3 string
	}{
		{
			name:  "Успешно получена информация о файле",
			want:  "неизвестно",
			want1: "неизвестно",
			want2: 0,
			want3: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3 := FileWithFuncAndLineNum(SkipNumOfStackFrame)

			if got != tt.want {
				t.Errorf("FileWithFuncAndLineNum() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FileWithFuncAndLineNum() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("FileWithFuncAndLineNum() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("FileWithFuncAndLineNum() got3 = %v, want %v", got3, tt.want3)
			}
		})

	}
}

func Test_splitCodeMessage(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name:  "Успешно разделение сообщения",
			args:  args{msg: "I1000 New entry"},
			want:  "I1000",
			want1: "New entry",
		},
		{
			name:  "Не успешное разделение сообщения",
			args:  args{msg: ""},
			want:  "",
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := splitCodeMessage(tt.args.msg)
			if got != tt.want {
				t.Errorf("splitCodeMessage() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("splitCodeMessage() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLogger_createEntry(t *testing.T) {
	cleanupTestFiles(t)
	t.Run("Успешно создан объект LogEntry", func(t *testing.T) {
		funcName, fileName, lineNum, filePath := FileWithFuncAndLineNum(SkipNumOfStackFrame)
		funcName = "main.TestWrite"
		fileName = "logger_test.go"
		lineNum = 16
		filePath = pathFile

		assert.Equal(t, funcName, "main.TestWrite")
		assert.Equal(t, fileName, "logger_test.go")
		assert.Equal(t, lineNum, 16)
		assert.Equal(t, filePath, pathFile)
	})
}

func TestLogger_createDetails(t *testing.T) {
	cleanupTestFiles(t)

	type fields struct {
		file     *os.File
		infoPC   *InfoPC
		filePath string
	}
	tests := []struct {
		name   string
		fields fields
		want   *DetailEntry
	}{
		{
			name: "Успешно создан объект DetailEntry",
			fields: fields{
				file:     nil,
				infoPC:   nil,
				filePath: "",
			},
			want: &DetailEntry{
				FunctionName: "неизвестно",
				FileName:     "неизвестно",
				LineNumber:   0,
				PathToFile:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				file:     tt.fields.file,
				infoPC:   tt.fields.infoPC,
				filePath: tt.fields.filePath,
			}
			if got := l.createDetails(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createDetails() = %v, want %v", got, tt.want)
			}
			l.file.Close()
		})
	}
}

func TestLogger_createMessage(t *testing.T) {
	cleanupTestFiles(t)

	type fields struct {
		file     *os.File
		infoPC   *InfoPC
		filePath string
	}
	type args struct {
		msg    string
		errStr *string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   MessageEntry
	}{
		{
			name: "Успешно создан объект MessageEntry",
			fields: fields{
				file:     nil,
				infoPC:   nil,
				filePath: pathFile,
			},
			want: MessageEntry{
				Code:    "",
				Message: "",
				Error:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				file:     tt.fields.file,
				infoPC:   tt.fields.infoPC,
				filePath: tt.fields.filePath,
			}
			assert.Equalf(t, tt.want, l.createMessage(tt.args.msg, tt.args.errStr), "createMessage(%v, %v)", tt.args.msg, tt.args.errStr)
			l.file.Close()
		})
	}
}

func TestLogger_logCustom(t *testing.T) {
	cleanupTestFiles(t)

	file, err := openFile(pathFile)
	if err != nil {
		log.Printf("Ошибка при открытии файла: %v", err)
		return
	}

	defer file.Close()
	type fields struct {
		file     *os.File
		infoPC   *InfoPC
		filePath string
	}
	type args struct {
		level    LogLevel
		message  string
		errStr   *string
		response *ResponseEntry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Успешно записан лог",
			fields: fields{
				file:     file,
				infoPC:   &InfoPC{},
				filePath: pathFile,
			},
			args: args{
				level:    LogLevelInfo,
				message:  "I2000 New entry",
				errStr:   nil,
				response: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				file:     tt.fields.file,
				infoPC:   tt.fields.infoPC,
				filePath: tt.fields.filePath,
			}
			l.logCustom(tt.args.level, tt.args.message, tt.args.errStr, tt.args.response)
			l.file.Close()
		})
	}
}

func TestLogger_LogWithResponseE(t *testing.T) {
	cleanupTestFiles(t)

	type fields struct {
		file     *os.File
		infoPC   *InfoPC
		filePath string
	}
	type args struct {
		msg        string
		statusCode int
		method     string
		url        string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Успешно записан лог",
			fields: fields{
				file:     nil,
				infoPC:   nil,
				filePath: pathFile,
			},
			args: args{
				msg:        "E2000 New entry",
				statusCode: 500,
				method:     "GET",
				url:        "/usrs/1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				file:     tt.fields.file,
				infoPC:   tt.fields.infoPC,
				filePath: tt.fields.filePath,
			}
			l.LogHttpErr(tt.args.msg, tt.args.statusCode, tt.args.method, tt.args.url)
			l.file.Close()
		})
	}
}

func TestLogger_LogWithResponseI(t *testing.T) {
	cleanupTestFiles(t)

	type fields struct {
		file     *os.File
		infoPC   *InfoPC
		filePath string
	}
	type args struct {
		msg        string
		statusCode int
		method     string
		url        string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Успешно записан лог",
			fields: fields{
				file:     nil,
				infoPC:   nil,
				filePath: pathFile,
			},
			args: args{
				msg:        "I2000 New entry",
				statusCode: 200,
				method:     "GET",
				url:        "/usrs/1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				file:     tt.fields.file,
				infoPC:   tt.fields.infoPC,
				filePath: tt.fields.filePath,
			}
			l.LogHttpI(tt.args.msg, tt.args.statusCode, tt.args.method, tt.args.url)
			l.file.Close()
		})
	}
}

func TestLogger_LogE(t *testing.T) {
	cleanupTestFiles(t)

	type fields struct {
		file     *os.File
		infoPC   *InfoPC
		filePath string
	}
	type args struct {
		msg string
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			fields: fields{
				file:     nil,
				infoPC:   nil,
				filePath: pathFile,
			},
			args: args{
				msg: "E3000 New entry",
				err: errors.New("error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				file:     tt.fields.file,
				infoPC:   tt.fields.infoPC,
				filePath: tt.fields.filePath,
			}
			l.LogE(tt.args.msg, tt.args.err)
			l.file.Close()
		})
	}
}

func TestLogger_LogW(t *testing.T) {
	cleanupTestFiles(t)

	type fields struct {
		file     *os.File
		infoPC   *InfoPC
		filePath string
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			fields: fields{
				file:     nil,
				infoPC:   nil,
				filePath: "",
			},
			args: args{msg: "W2000 New entry"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				file:     tt.fields.file,
				infoPC:   tt.fields.infoPC,
				filePath: tt.fields.filePath,
			}
			l.LogW(tt.args.msg)
			l.file.Close()
		})
	}
}

func TestLogger_LogI(t *testing.T) {
	cleanupTestFiles(t)

	type fields struct {
		file     *os.File
		infoPC   *InfoPC
		filePath string
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Успешно записан лог",
			fields: fields{
				file:     nil,
				infoPC:   nil,
				filePath: "",
			},
			args: args{msg: "I2000 New entry"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				file:     tt.fields.file,
				infoPC:   tt.fields.infoPC,
				filePath: tt.fields.filePath,
			}
			l.LogI(tt.args.msg)
			l.file.Close()
		})
	}
}

func TestNewLogger(t *testing.T) {
	cleanupTestFiles(t)

	t.Run("Успешно создан объект Logger", func(t *testing.T) {
		logger, err := NewLogger(pathFile)

		assert.NoError(t, err)
		assert.NotNil(t, logger)
		assert.NotNil(t, logger.file)
		assert.Equal(t, pathFile, logger.filePath)
		assert.NotNil(t, logger.infoPC)
		logger.Close()
	})

	t.Run("Не успешное создание объекта Logger", func(t *testing.T) {
		invalidPath := "???.//invalid_path.json"
		logger, err := NewLogger(invalidPath)
		assert.Error(t, err)
		assert.Nil(t, logger)
	})

	t.Run("Не удалось найти файл", func(t *testing.T) {
		cleanupTestFiles(t)
		logger, err := NewLogger("")
		if logger.filePath == "" {
			logger.filePath = pathFile
		}
		assert.NoError(t, err)
		logger.Close()
		cleanupTestFiles(t)
	})
}

func cleanupTestFiles(t *testing.T) {
	t.Cleanup(func() {
		files := []string{"test_log.json", "logCustom.json"}
		for _, file := range files {
			if _, err := os.Stat(file); err == nil {
				os.Remove(file)
			}
		}
	})
}
