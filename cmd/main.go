package main

import (
	"context"
	"fmt"
	"github.com/beriloqueiroz/psi-mgnt/config"
	"github.com/beriloqueiroz/psi-mgnt/internal/application"
	infra "github.com/beriloqueiroz/psi-mgnt/internal/infra/database"
	"github.com/beriloqueiroz/psi-mgnt/internal/infra/web/api_routes"
	webserver "github.com/beriloqueiroz/psi-mgnt/internal/infra/web/server"
	routes_view "github.com/beriloqueiroz/psi-mgnt/internal/infra/web/view_routes"
	"log/slog"
	"os"
	"os/signal"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// graceful exit
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	initCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// config logs
	// aqui no lugar do Stdout poderia ser um db ou elasticsearch por exemplo
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// enviroment configs
	configs, err := config.LoadConfig([]string{"."})
	if err != nil {
		panic(err)
	}

	// telemetry
	//otel := otel_b.OtelB{}
	//shutdown, err := otel.InitTraceProvider("web server psi-mgmt", configs.OtelExporterEndpoint)
	//if err != nil {
	//	slog.Error(err.Error(), err)
	//}
	//defer func() {
	//	fmt.Println("oiaaaa")
	//	if err := shutdown(initCtx); err != nil {
	//		slog.Error("failed shutdown TraceProvider: %w", err)
	//	}
	//}()

	//server := webserver.NewWebServer(configs.WebServerPort, otel.WithRouteTag)
	server := webserver.NewWebServer(configs.WebServerPort, nil)

	// repositories and gateways
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

	// usecases
	createSessionUseCase := application.NewCreateSessionUseCase(sessionRepository)
	updateSessionUseCase := application.NewUpdateSessionUseCase(sessionRepository)
	findSessionUseCase := application.NewFindSessionUseCase(sessionRepository)
	deleteSessionUseCase := application.NewDeleteSessionUseCase(sessionRepository)
	listSessionsUsecase := application.NewListSessionsUseCase(sessionRepository)
	searchPatientsUseCase := application.NewSearchPatientsUseCase(sessionRepository)
	searchProfessionalsUseCase := application.NewSearchProfessionalsUseCase(sessionRepository)
	createProfessionalUseCase := application.NewCreateProfessionalUseCase(sessionRepository)

	// api api_routes
	createSessionRoute := api_routes.NewCreateSessionRoute(*createSessionUseCase)
	createProfessionalRoute := api_routes.NewCreateProfessionalRoute(*createProfessionalUseCase)
	searchPatientsRoute := api_routes.NewSearchPatientsRoute(*searchPatientsUseCase)
	searchProfesionalsRoute := api_routes.NewSearchProfessionalsRoute(*searchProfessionalsUseCase)
	deleteSessionRoute := api_routes.NewDeleteSessionRoute(*deleteSessionUseCase)
	listSessionRoute := api_routes.NewListSessionsRoute(*listSessionsUsecase)

	server.AddRoute("POST /api", createSessionRoute.Handler)
	server.AddRoute("POST /api/professional", createProfessionalRoute.Handler)
	server.AddRoute("GET /api", listSessionRoute.Handler)
	server.AddRoute("DELETE /api/{id}", deleteSessionRoute.Handler)
	server.AddRoute("GET /api/patient", searchPatientsRoute.Handler)
	server.AddRoute("GET /api/professional", searchProfesionalsRoute.Handler)

	// views
	sessionRouteView := routes_view.NewSessionRouteView(*createSessionUseCase, *updateSessionUseCase, *findSessionUseCase)
	server.AddRoute("GET /session", sessionRouteView.HandlerGet)
	server.AddRoute("GET /session/{id}", sessionRouteView.HandlerGet)
	server.AddRoute("POST /session", sessionRouteView.HandlerPost)
	server.AddRoute("POST /session/{id}", sessionRouteView.HandlerPost)

	createProfessionalRouteView := routes_view.NewCreateProfessionalRouteView(*createProfessionalUseCase)

	server.AddRoute("POST /professional", createProfessionalRouteView.HandlerPost)
	server.AddRoute("GET /professional", createProfessionalRouteView.HandlerGet)

	listSessionRouteView := routes_view.NewListSessionRouteView(*listSessionsUsecase, *deleteSessionUseCase)
	server.AddRoute("GET /sessions", listSessionRouteView.HandlerGet)
	server.AddRoute("POST /sessions/{id}", listSessionRouteView.HandlerPost)

	server.AddRoute("GET /", routes_view.NewHomeRouteView().HandlerGet)

	fmt.Println("Starting web server on port", configs.WebServerPort)
	server.Start()

	// Wait for interruption.
	select {
	case <-sigCh:
		slog.Warn("Shutting down gracefully, CTRL+C pressed...")
	case <-initCtx.Done():
		slog.Warn("Shutting down due to other reason...")
	}
}
