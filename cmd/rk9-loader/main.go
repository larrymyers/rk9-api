package main

import (
	"context"
	"fmt"

	"larrymyers.com/rk9api/api"
	"larrymyers.com/rk9api/rk9"
)

func main() {
	builder := api.NewConnectionBuilder()
	conn, err := builder.WithEnvVars().GetConnection()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := api.NewClient(conn)

	events, err := rk9.GetEvents()
	for _, event := range events {
		err := client.UpsertEvent(context.Background(), event)
		if err != nil {
			panic(err)
		}
	}

	events, err = client.GetEvents(context.Background())
	if err != nil {
		panic(err)
	}

	for _, event := range events {
		if event.HasStarted() {
			fmt.Printf("Event: %s\n", event.Name)

			rounds, err := rk9.GetRounds(event)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Rounds: %d\n", rounds.Masters)

			for round := range rounds.Masters {
				matches, err := rk9.GetRound(event, rk9.Masters, round+1)
				if err != nil {
					panic(err)
				}

				fmt.Printf("Round %d: %d matches\n", round, len(matches))

				for _, match := range matches {
					err := client.UpsertMatch(context.Background(), match)
					if err != nil {
						panic(err)
					}
				}
			}

			break
		}
	}
}
