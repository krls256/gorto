package http

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"log/slog"
)

type Server struct {
	logger   *slog.Logger
	cfg      Config
	fiberApp *fiber.App
}

func NewServer(ctx context.Context, appName string, logger *slog.Logger, cfg Config, handlers []Handler) *Server {
	s := &Server{
		logger: logger,
		cfg:    cfg,
		fiberApp: fiber.New(fiber.Config{
			StrictRouting: false,
			CaseSensitive: false,
			ReadTimeout:   cfg.ReadTimeout,
			WriteTimeout:  cfg.WriteTimeout,
			IdleTimeout:   cfg.IdleTimeout,
			AppName:       appName,
		})}

	s.fiberApp.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c fiber.Ctx, e any) {
			slog.Error("http exception", e)
		},
	}))

	s.fiberApp.Use(cors.New(cors.ConfigDefault))

	for _, h := range handlers {
		h.Register(s.fiberApp)
	}

	return s
}

func (s *Server) AsyncRun() {
	go func() {
		if err := s.Run(); err != nil {
			s.logger.Error(err.Error())
		}
	}()
}

func (s *Server) Run() error {
	if !s.cfg.Silent {
		s.logger.Info("Running server...")
	}

	return s.fiberApp.Listen(fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port), fiber.ListenConfig{
		DisableStartupMessage: s.cfg.Silent,
	})
}

func (s *Server) Shutdown(ctx context.Context) error {
	if !s.cfg.Silent {
		s.logger.Info("Shutdown server...")
	}

	defer func() {
		if !s.cfg.Silent {
			s.logger.Info("Server successfully stopped.")
		}
	}()

	if err := s.fiberApp.ShutdownWithContext(ctx); err != nil {
		if !s.cfg.Silent {
			s.logger.Error("Server forced to shutdown:", err)
		}

		return err
	}

	return nil
}
