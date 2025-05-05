package main

import (
	"fmt"
	"log"
	"yadro-test/config"
	"yadro-test/internal"
)

func main() {
	cfg, err := config.ParseConfig("files/cfg.json")
	if err != nil {
		log.Fatalf("Error parsing config.json: %s", err)
	}
	events, err := internal.LoadEvents("files/events")
	if err != nil {
		log.Fatalf("Error loading events: %s", err)
	}

	competition, err := config.CreateCompetition(cfg)
	if err != nil {
		log.Fatalf("Error creating competition: %s", err)
	}

	for _, event := range events {
		out := competition.ProcessEvent(event)
		fmt.Println(out)
	}

	fmt.Println("`Resulting table`")
	fmt.Println(competition.MakeReport())
}
