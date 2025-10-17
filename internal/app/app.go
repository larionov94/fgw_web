package app

import (
	"FGW_WEB/internal/config"
	"FGW_WEB/internal/config/db"
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const addr = ":7777"
const fileEnv = ".env"

func StartApp() {
	logger, err := common.NewLogger("")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	configDB, err := config.NewMSSQLCfg(logger, fileEnv)
	if err != nil {
		logger.LogE(msg.E3000, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mssqlDB, err := db.NewConnMSSQL(ctx, configDB, logger)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(mssqlDB)

	mux := http.NewServeMux()

	// TODO: Реализовать роутинг

	server := config.NewServer(addr, mux, logger)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = server.StartServer(ctx); err != nil {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()
	<-quit
	log.Println("Получен сигнал остановки сервера...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		logger.LogE(msg.E3102, err)
	}
}
