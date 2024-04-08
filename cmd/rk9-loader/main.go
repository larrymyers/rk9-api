package main

import (
	"context"

	"larrymyers.com/rk9api/api"
	"larrymyers.com/rk9api/rk9"
)

func main() {
	builder := api.NewConnectionBuilder()
	conn, err := builder.WithEnvVars().GetConnection()
	if err != nil {
		panic(err)
	}

	client := api.NewClient(conn)

	events, err := rk9.GetEvents()
	for _, event := range events {
		err := client.UpsertEvent(context.Background(), event)
		if err != nil {
			panic(err)
		}
	}
}
