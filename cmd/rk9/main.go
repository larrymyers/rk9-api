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
		fmt.Printf("%s %s\n", event.Name, event.Location)
	}
}
