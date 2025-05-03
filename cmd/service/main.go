package main

import (
	"fmt"
	"yadro-test/config"
	"yadro-test/internal"
)

func main() {
	cfg, err := config.ParseConfig("files/cfg.json")
	if err != nil {
		fmt.Printf("Error parsing config: %s\n", err)
	}
	fmt.Println(*cfg)
	events, err := internal.LoadEvents("files/events")
	if err != nil {
		fmt.Printf("Error loading events: %s\n", err)
	}
	fmt.Println(events)
}
