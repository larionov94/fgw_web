package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
}

func validEnvFile() (*os.File, error) {
	file, err := openFile("test.env")
	if err != nil {
		fmt.Println("Error opening file:", err)
	}

	_, err = file.Write([]byte(validEnv))
	if err != nil {
		fmt.Println("Error writing file:", err)
	}

	return file, nil
}

var validEnv = `
MSSQL_DRIVER=FreeTDS
MSSQL_SERVER=111.111.111.111,1111
MSSQL_NAME=11_11_TEST
MSSQL_USER=11q11
MSSQL_PASSWD=11q11
MSSQL_CHARSET=WINDOWS-1251
`

const pathToEnv = "test.env"

func TestNewMSSQLCfg(t *testing.T) {

	t.Run("Успешно: данные в файле", func(t *testing.T) {
		file, err := validEnvFile()
		if err != nil {
			t.Error(err)
		}
		defer file.Close()
		defer os.Remove("test.env")

		got, err := NewMSSQLCfg(nil, pathToEnv)

		assert.NoError(t, err, "NewMSSQLCfg() error = %v, wantErr %v", err, false)
		assert.Equal(t, got.MSSQL.Driver, "FreeTDS", "NewMSSQLCfg() got = %v, want %v", got.MSSQL.Driver, "FreeTDS")
	})

	t.Run("Не успешно: нет пути до файла", func(t *testing.T) {
		got, err := NewMSSQLCfg(nil, "")

		assert.Error(t, err, "NewMSSQLCfg() error = %v, wantErr %v", err, false)
		assert.Nil(t, got, "NewMSSQLCfg() error = %v, wantErr %v", err, false)
	})
}

func Test_loadEnvFile(t *testing.T) {
	defer os.Remove("test.env")
	t.Run("Успешное создание файла", func(t *testing.T) {
		err := loadEnvFile(pathToEnv)
		assert.NoError(t, err, "loadEnvFile() error = %v, wantErr %v", err, false)
	})

	t.Run("Не успешное создание файла", func(t *testing.T) {
		err := loadEnvFile("")
		assert.Error(t, err, "loadEnvFile() error = %v, wantErr %v", err, true)
	})
}
