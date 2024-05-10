package main

import (
	"context"
	"fmt"

	"github.com/beriloqueiroz/psi-mgnt/config"
	"github.com/beriloqueiroz/psi-mgnt/internal/application"
	infra "github.com/beriloqueiroz/psi-mgnt/internal/infra/database"
	"github.com/beriloqueiroz/psi-mgnt/internal/infra/web/routes"
	webserver "github.com/beriloqueiroz/psi-mgnt/internal/infra/web/server"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := config.LoadConfig([]string{"."})
	if err != nil {
		panic(err)
	}

	server := webserver.NewWebServer(configs.WebServerPort)

	initCtx := context.Background()
	var sessionRepository application.SessionRepositoryInterface
	sessionRepository, err = infra.NewMongoSessionRepository(
		initCtx,
		configs.DBUri,
		configs.DBPatientCollection,
		configs.DBSessionCollection,
		configs.DBDatabase)
	if err != nil {
		panic(err)
	}
	createSessionUseCase := application.NewCreateSessionUseCase(sessionRepository)
	createSessionRoute := routes.NewCreateSessionRoute(*createSessionUseCase)

	server.AddRoute("POST /", createSessionRoute.Handler)

	fmt.Println("Starting web server on port", configs.WebServerPort)
	server.Start()
}
