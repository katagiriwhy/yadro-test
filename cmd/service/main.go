package main

import (
	"fmt"
	"log"
	"os"
	competition2 "yadro-test/competition"
	"yadro-test/config"
	events2 "yadro-test/events"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("/home/danil/yadro-test/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg, err := config.ParseConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatalf("Error parsing config.json: %s", err)
	}
	events, err := events2.LoadEvents(os.Getenv("EVENTS_PATH"))
	if err != nil {
		log.Fatalf("Error loading events: %s", err)
	}

	competition, err := competition2.CreateCompetition(cfg)
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
