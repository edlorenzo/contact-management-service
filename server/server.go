package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"contact-management-service/config"
	"contact-management-service/internal/contacts"
	"contact-management-service/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	defaultAddr     = ":8000"
	shutdownTimeout = time.Second * 5
)

type Server struct {
	*http.Server
}

func (s *Server) GracefulShutdown() error {
	timeout := shutdownTimeout
	done := make(chan error, 1)
	go func() {
		ctx := context.Background()
		var cancel context.CancelFunc
		if timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		log.Info().Msg("server: shutting down gracefully...")
		done <- s.Shutdown(ctx)
		log.Info().Msg("server: shutdown")
	}()
	return <-done
}

func New(
	conf *config.Config,
	svc contacts.Service,
	repo contacts.Repo,
) *Server {
	s := &Server{}
	r := gin.New()
	r.Use(RequestLogger())
	router.Route(r, svc, repo, conf)
	if conf.AppConfig.Addr == "" {
		conf.AppConfig.Addr = defaultAddr
	}
	s.Server = &http.Server{
		Addr:         conf.AppConfig.Addr,
		ReadTimeout:  conf.AppConfig.ReadTimeout,
		WriteTimeout: conf.AppConfig.WriteTimeout,
		Handler:      r,
	}
	return s
}

func (s *Server) Run() error {
	log.Info().Msg(fmt.Sprintf("server: running on %s", s.Addr))
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
