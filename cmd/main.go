package main

import (
	"github.com/juanmabaracat/diagnosis-service/internal/app"
	"github.com/juanmabaracat/diagnosis-service/internal/infrastracture/http"
	"github.com/juanmabaracat/diagnosis-service/internal/infrastracture/storage/memory"
)

func main() {
	repository := memory.NewRepository()
	appServices := app.NewServices(&repository, &repository)
	server := http.NewServer(appServices)
	server.Run(":8080")
}
