package main

import (
	"fmt"

	"larrymyers.com/rk9"
)

func main() {
	events, err := rk9.GetEvents()
	if err != nil {
		panic(err)
	}

	for _, event := range events {
		fmt.Printf("%s\n%s - %s\n%s\n%s\n%s\n\n", event.Name, event.StartDate, event.EndDate, event.Location, event.DetailsURL, event.PairingsURL)
	}
}
