package config

import (
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const fileEnv = ".env"

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

func NewMSSQLCfg(logger *common.Logger) (*MSSQLCfg, error) {
	if err := loadEnvFile(); err != nil {
		return nil, err
	}

	fmt.Println(os.Getenv("MSSQL_DRIVER"))
	fmt.Println(os.Getenv("MSSQL_SERVER"))
	fmt.Println(os.Getenv("MSSQL_NAME"))
	fmt.Println(os.Getenv("MSSQL_USER"))
	fmt.Println(os.Getenv("MSSQL_PASSWD"))
	fmt.Println(os.Getenv("MSSQL_CHARSET"))

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

func loadEnvFile() error {
	//envPath := filepath.Join(fileEnv)
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("%s: %w", msg.E3003, err)
	}

	return nil
}
