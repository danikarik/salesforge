package app

import (
	"net/http"

	"github.com/danikarik/salesforge/internal/model"
	"github.com/gin-gonic/gin"
)

type Service struct {
	mux   *gin.Engine
	store model.SequenceStore
}

type Config struct {
	Store model.SequenceStore
	// Additional configuration options can be added here in the future.
}

// NewService creates a new Service instance with the provided options.
func NewService(cfg Config) *Service {
	srv := &Service{store: cfg.Store}

	r := gin.Default()
	r.POST("/sequences", srv.createSequence)
	r.GET("/sequences/:id", srv.fetchSequence)
	r.PUT("/sequences/:id", srv.updateSequence)
	r.PUT("/sequences/:id/steps/:step_id", srv.updateStep)
	r.DELETE("/sequences/:id/steps/:step_id", srv.deleteStep)

	srv.mux = r
	return srv
}

func (s *Service) Handler() http.Handler {
	return s.mux.Handler()
}
