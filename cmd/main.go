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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// config logs
	// aqui no lugar do Stdout poderia ser um db ou elasticsearch por exemplo
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// enviroment configs
	configs, err := config.LoadConfig([]string{"."})
	if err != nil {
		panic(err)
	}

	server := webserver.NewWebServer(configs.WebServerPort)

	err = startRoutes(server, ctx, *configs)
	if err != nil {
		panic(err)
	}

	srvErr := make(chan error, 1)
	go func() {
		fmt.Println("Starting web server on port", configs.WebServerPort)
		srvErr <- server.Start()
	}()

	// Wait for interruption.
	select {
	case <-sigCh:
		{
			slog.Warn("Shutting down gracefully, CTRL+C pressed...")
			return
		}
	case <-srvErr:
		{
			slog.Error("Shutting down gracefully, Server error pressed...")
			return
		}
	case <-ctx.Done():
		slog.Warn("Shutting down due to other reason...")
	}
}

func startRoutes(server *webserver.WebServer, ctx context.Context, configs config.Conf) error {

	var err error
	// repositories and gateways
	var sessionRepository application.SessionRepositoryInterface
	sessionRepository, err = infra.NewMongoSessionRepository(
		ctx,
		configs.DBUri,
		configs.DBPatientCollection,
		configs.DBProfessionalCollection,
		configs.DBSessionCollection,
		configs.DBDatabase)
	if err != nil {
		return err
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
	server.AddRoute("POST /session/", sessionRouteView.HandlerPost)
	server.AddRoute("POST /session/{id}", sessionRouteView.HandlerPost)

	createProfessionalRouteView := routes_view.NewCreateProfessionalRouteView(*createProfessionalUseCase)

	server.AddRoute("POST /professional", createProfessionalRouteView.HandlerPost)
	server.AddRoute("GET /professional", createProfessionalRouteView.HandlerGet)

	listSessionRouteView := routes_view.NewListSessionRouteView(*listSessionsUsecase, *deleteSessionUseCase)
	server.AddRoute("GET /sessions", listSessionRouteView.HandlerGet)
	server.AddRoute("POST /sessions/{id}", listSessionRouteView.HandlerPost)

	server.AddRoute("GET /", routes_view.NewHomeRouteView().HandlerGet)
	return nil
}
