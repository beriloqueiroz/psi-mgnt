package main

import (
	"database/sql"
	"fmt"

	"github.com/beriloqueiroz/psi-mgnt/config"
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

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	server := webserver.NewWebServer(configs.WebServerPort)
	msgRoute := routes.NewWebMsgRoute("olá dev loco")
	server.AddRoute("GET /", msgRoute.Handler)

	fmt.Println("Starting web server on port", configs.WebServerPort)
	server.Start()
}
