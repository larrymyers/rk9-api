package main

import (
	"context"
	"flag"

	"larrymyers.com/rk9api/api"
)

func main() {
	var initDB bool
	var startServer bool

	flag.BoolVar(&initDB, "init-db", false, "initialize new database")
	flag.BoolVar(&startServer, "start", false, "start server")
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
		return
	}
}
