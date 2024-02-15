package main

import (
	"github.com/juanmabaracat/diagnosis-service/internal/app"
	"github.com/juanmabaracat/diagnosis-service/internal/infrastracture/http"
	"github.com/juanmabaracat/diagnosis-service/internal/infrastracture/storage/memory"
)

// @Title			Patient Diagnoses API
// @Version			1.0.0
// @Description		This is API service to handle patient diagnoses
// @host			localhost:8080
// @BasePath 		/api/v1
func main() {
	repository := memory.NewRepository()
	appServices := app.NewServices(&repository, &repository)
	server := http.NewServer(appServices)
	server.Run(":8080")
}
