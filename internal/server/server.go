package server

import (
	"fmt"
	_ "pg-sh-scripts/docs"
	"pg-sh-scripts/internal/config"

	"github.com/gin-gonic/gin"
)

type (
	IServer interface {
		Run() error
		Shutdown() error
	}

	Server struct{}
)

func getServer() *gin.Engine {
	return gin.New()
}

func setServerMode(cfg *config.Config) {
	mode := gin.ReleaseMode

	switch cfg.Project.Mode {
	case "local":
		mode = gin.DebugMode
	case "dev":
		mode = gin.TestMode
	}

	gin.SetMode(mode)
}

func setServerProxies(r *gin.Engine, cfg *config.Config) error {
	r.ForwardedByClientIP = true
	if err := r.SetTrustedProxies(cfg.Api.TrustedProxies); err != nil {
		return err
	}
	return nil
}

func runServer(r *gin.Engine, cfg *config.Config) error {
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		return fmt.Errorf("cannot run main router: %w", err)
	}
	return nil
}

func (s *Server) Run() error {
	if err := setDotEnv(); err != nil {
		return err
	}

	cfg := config.GetConfig()

	pgClient, err := setPgConn()
	if err != nil {
		return err
	}
	if err := setMigration(pgClient.GetDB()); err != nil {
		return err
	}

	setServerMode(cfg)

	r := getServer()

	setServeMiddleware(r)
	if err := setServerProxies(r, cfg); err != nil {
		return err
	}

	setSwagger(r)
	setV1Handlers(r, cfg)

	if err := runServer(r, cfg); err != nil {
		return err
	}

	return nil
}

func (s *Server) Shutdown() error {
	if err := closePgConn(); err != nil {
		return err
	}
	return nil
}

func GetServer() IServer {
	return &Server{}
}
