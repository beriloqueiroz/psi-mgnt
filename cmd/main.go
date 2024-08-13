package main

import (
	"context"
	"fmt"
	routes_view "github.com/beriloqueiroz/psi-mgnt/internal/infra/web/view_routes"

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
		configs.DBProfessionalCollection,
		configs.DBSessionCollection,
		configs.DBDatabase)
	if err != nil {
		panic(err)
	}
	createSessionUseCase := application.NewCreateSessionUseCase(sessionRepository)

	createProfessionalUseCase := application.NewCreateProfessionalUseCase(sessionRepository)

	createSessionRoute := routes.NewCreateSessionRoute(*createSessionUseCase)

	searchPatientsUseCase := application.NewSearchPatientsUseCase(sessionRepository)
	searchPatientsRoute := routes.NewSearchPatientsRoute(*searchPatientsUseCase)

	deleteSessionUseCase := application.NewDeleteSessionUseCase(sessionRepository)
	deleteSessionRoute := routes.NewDeleteSessionRoute(*deleteSessionUseCase)

	listSessionsUsecase := application.NewListSessionsUseCase(sessionRepository)
	listSessionRoute := routes.NewListSessionsRoute(*listSessionsUsecase)

	server.AddRoute("POST /api", createSessionRoute.Handler)
	server.AddRoute("GET /api", listSessionRoute.Handler)
	server.AddRoute("DELETE /api/{id}", deleteSessionRoute.Handler)
	server.AddRoute("GET /api/patient", searchPatientsRoute.Handler)

	// views
	createSessionRouteView := routes_view.NewCreateSessionRouteView(*createSessionUseCase)
	server.AddRoute("GET /session", createSessionRouteView.HandlerGet)
	server.AddRoute("POST /session", createSessionRouteView.HandlerPost)

	createProfessionalRouteView := routes_view.NewCreateProfessionalRouteView(*createProfessionalUseCase)

	server.AddRoute("POST /professional", createProfessionalRouteView.HandlerPost)
	server.AddRoute("GET /professional", createProfessionalRouteView.HandlerGet)

	listSessionRouteView := routes_view.NewListSessionRouteView(*listSessionsUsecase, *deleteSessionUseCase)
	server.AddRoute("GET /sessions", listSessionRouteView.HandlerGet)
	server.AddRoute("POST /sessions/{id}", listSessionRouteView.HandlerPost)

	fmt.Println("Starting web server on port", configs.WebServerPort)
	server.Start()
}
