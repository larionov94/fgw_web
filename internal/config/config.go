package config

import (
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

type MSSQLEntryCfg struct {
	Driver  string `env:"driver"`
	Server  string `env:"server"`
	Name    string `env:"name"`
	User    string `env:"user"`
	Passwd  string `env:"passwd"`
	Charset string `env:"charset"`
}

type MSSQLCfg struct {
	MSSQL  MSSQLEntryCfg
	logger *common.Logger
}

func NewMSSQLCfg(logger *common.Logger, pathFile string) (*MSSQLCfg, error) {
	if err := loadEnvFile(pathFile); err != nil {
		return nil, err
	}

	return &MSSQLCfg{
		MSSQL: MSSQLEntryCfg{
			Driver:  os.Getenv("MSSQL_DRIVER"),
			Server:  os.Getenv("MSSQL_SERVER"),
			Name:    os.Getenv("MSSQL_NAME"),
			User:    os.Getenv("MSSQL_USER"),
			Passwd:  os.Getenv("MSSQL_PASSWD"),
			Charset: os.Getenv("MSSQL_CHARSET"),
		},
		logger: logger,
	}, nil
}

func loadEnvFile(pathFile string) error {
	envPath := filepath.Join(pathFile)
	err := godotenv.Load(envPath)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.E3003, err)
	}

	return nil
}

func GenerateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// UpdateSessionStoreOptions - обновление настроек хранилища (добавить в config)
func UpdateSessionStoreOptions(maxAge int) {
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}
