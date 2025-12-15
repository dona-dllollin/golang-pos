package http

import (
	"encoding/json"
	"net/http"

	"github.com/dona-dllollin/belajar-clean-arch/internal/config"
	productHttp "github.com/dona-dllollin/belajar-clean-arch/internal/delivery/http/producthandler/handler"
	customMiddleware "github.com/dona-dllollin/belajar-clean-arch/internal/middleware"
	"github.com/dona-dllollin/belajar-clean-arch/pkgs/logger"
	"github.com/dona-dllollin/belajar-clean-arch/pkgs/validation"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

type Server struct {
	engine    *chi.Mux
	db        *pgx.Conn
	validator validation.Validation
	cfg       *config.Config
}

func NewServer(
	validator validation.Validation,
	db *pgx.Conn,
) *Server {
	return &Server{
		engine:    chi.NewRouter(),
		db:        db,
		cfg:       config.LoadConfig(),
		validator: validator,
	}
}

func (s Server) Run() error {

	// use middleware
	s.engine.Use(middleware.Recoverer)
	s.engine.Use(customMiddleware.CORSMiddleware)

	// load all route
	s.MapRoute()

	// check http service
	s.engine.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "Welcome to Ecommerce Clean Architecture"}); err != nil {
			// If an error occurs during encoding, log it and send an error response
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Fatalf("Error encoding JSON: %v", err)
		}
	})

	// start http server
	logger.Info("HTTP server is listening on PORT: ", s.cfg.Port)
	if err := http.ListenAndServe(s.cfg.Port, s.engine); err != nil {
		logger.Fatalf("Running HTTP server: %d", err)
	}
	return nil
}

func (s Server) MapRoute() {
	s.engine.Route("/api/v1", func(r chi.Router) {
		r.Route("/products", func(r chi.Router) {
			productHttp.Routes(r, s.db, s.validator, s.cfg.ImagePath)
		})
	})
}
