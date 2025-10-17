package db

import (
	"FGW_WEB/internal/config"
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

func NewConnMSSQL(ctx context.Context, configDB *config.MSSQLCfg, logger *common.Logger) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("%s://%s:%s@%s?database=%s&charset=%s",
		configDB.MSSQL.Driver,
		configDB.MSSQL.User,
		configDB.MSSQL.Passwd,
		configDB.MSSQL.Server,
		configDB.MSSQL.Name,
		configDB.MSSQL.Charset)
	db, err := sql.Open("mssql", dataSourceName)
	if err != nil {
		logger.LogE(msg.E3200, err)

		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err = db.PingContext(pingCtx); err != nil {
		Close(db)
		log.Printf("%s: %v", msg.E3201, err)

		return nil, err
	}

	return db, nil
}

func Close(db *sql.DB) {
	if db == nil {
		return
	}

	if err := db.Close(); err != nil {
		log.Printf("%s: %v", msg.E3201, err)

		return
	}
	log.Printf(msg.I2200)
}
