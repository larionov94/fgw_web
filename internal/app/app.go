package app

import (
	"FGW_WEB/internal/config"
	"FGW_WEB/internal/config/db"
	"FGW_WEB/internal/handler"
	"FGW_WEB/internal/handler/http_web"
	"FGW_WEB/internal/handler/json_api"
	"FGW_WEB/internal/repository"
	"FGW_WEB/internal/service"
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/sessions"
)

const addr = ":7777"
const fileEnv = ".env"

var store *sessions.CookieStore

func StartApp() {
	config.InitSessionStore()

	logger, err := common.NewLogger("")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	authMiddleware := handler.NewAuthMiddleware(config.Store, logger)

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

	repoRole := repository.NewRoleRepo(mssqlDB, logger)
	serviceRole := service.NewRoleService(repoRole, logger)
	handlerRoleJSON := json_api.NewRoleHandlerJSON(serviceRole, logger)
	handlerRoleHTML := http_web.NewRoleHandlerHTML(serviceRole, logger)

	repoPerformer := repository.NewPerformerRepo(mssqlDB, logger)
	servicePerformer := service.NewPerformerService(repoPerformer, logger)
	handlerPerformerJSON := json_api.NewPerformerHandlerJSON(servicePerformer, logger)
	handlerPerformerHTML := http_web.NewPerformerHandlerHTML(servicePerformer, serviceRole, logger, authMiddleware)

	mux := http.NewServeMux()

	handlerRoleJSON.ServerHTTPJSONRouter(mux)
	handlerRoleHTML.ServerHTTPHTMLRouter(mux)

	handlerPerformerJSON.ServeHTTPJSONRouter(mux)
	handlerPerformerHTML.ServeHTTPHTMLRouter(mux)

	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web/"))))

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
