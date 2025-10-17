package config

import (
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"context"
	"errors"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	logger     *common.Logger
}

func NewServer(addr string, handler http.Handler, logger *common.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		logger: logger,
	}
}

// StartServer запускает HTTP сервер и блокирует, пока он в работе.
func (s *Server) StartServer(ctx context.Context) error {
	s.logger.LogI(msg.I2100 + s.httpServer.Addr)
	errCh := make(chan error, 1)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.LogE(msg.E3102, err)
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

// Shutdown выключает сервер корректно, ожидая завершения текущих запросов.
func (s *Server) Shutdown(ctx context.Context) error {
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		s.logger.LogE(msg.E3102, err)

		return err
	}
	return nil
}
