package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"larrymyers.com/rk9api/api"
	"larrymyers.com/rk9api/server"
)

func main() {
	var initDB bool
	var startServer bool

	flag.BoolVar(&initDB, "init-db", false, "initialize new database")
	flag.BoolVar(&startServer, "server", false, "start server")
	flag.Parse()

	if initDB {
		builder := api.NewConnectionBuilder()
		conn, err := builder.WithEnvVars().GetConnection()
		if err != nil {
			panic(err)
		}

		err = api.ApplySchema(context.Background(), conn)
		if err != nil {
			panic(err)
		}
		return
	}

	if startServer {
		builder := api.NewConnectionBuilder()
		conn, err := builder.WithEnvVars().GetConnection()
		if err != nil {
			panic(err)
		}

		client := api.NewClient(conn)

		router := server.NewRouter(client)

		host := "0.0.0.0"
		hostEnv := os.Getenv("SERVER_HOST_IP")
		if len(hostEnv) > 0 {
			host = hostEnv
		}

		port := "3000"
		portEnv := os.Getenv("SERVER_PORT")
		if len(portEnv) > 0 {
			port = portEnv
		}

		host = host + ":" + port

		log.Println("Server started: " + host)

		err = http.ListenAndServe(host, router)
		if err != nil {
			log.Printf("error shutting down server: %v", err)
		}

		return
	}
}
