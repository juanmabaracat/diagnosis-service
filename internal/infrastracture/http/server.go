package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/juanmabaracat/diagnosis-service/docs"
	"github.com/juanmabaracat/diagnosis-service/internal/app"
	"github.com/juanmabaracat/diagnosis-service/internal/infrastracture/http/diagnoses"
	"github.com/swaggo/http-swagger/v2"
	"log"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	appServices app.Services
	router      chi.Router
}

func NewServer(services app.Services) *Server {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(30 * time.Second))
	router.Use(commonMiddleware)
	server := &Server{
		appServices: services,
		router:      router,
	}

	server.addHTTPRoutes()
	return server
}

func (s *Server) addHTTPRoutes() {
	handler := diagnoses.NewHandler(s.appServices.DiagnosisServices)
	s.router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))
	s.router.Route("/api/v1", func(r chi.Router) {
		r.Get("/patient/diagnoses", handler.GetDiagnoses)
		r.Post("/patient/{"+diagnoses.PatientIDURLParam+"}/diagnoses", handler.AddDiagnosis)
	})
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})
}

func (s *Server) Run(port string) {
	slog.Info("Listening on http://localhost" + port)
	err := http.ListenAndServe(port, s.router)
	if err != nil {
		log.Fatal(err)
	}
}
