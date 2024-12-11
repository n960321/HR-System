package server

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Server struct {
	app *http.Server
}

type Config struct {
	Port string `mapstructure:"port"`
}

func NewServer(cfg Config, handler http.Handler) *Server {
	return &Server{
		app: &http.Server{
			Addr:    cfg.Port,
			Handler: handler,
		},
	}
}
func (svr *Server) Run() {
	go func() {
		log.Info().Msgf("Server Start Listening on %s", svr.app.Addr)
		if err := svr.app.ListenAndServe(); err != nil {
			log.Fatal().Err(err)
		}
	}()
}

func (svr *Server) Shutdown(ctx context.Context) {
	if err := svr.app.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server show failed")
	}
}
