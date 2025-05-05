package main

import (
	"fmt"
	"log"
	"os"
	"yadro-test/config"
	"yadro-test/internal"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg, err := config.ParseConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatalf("Error parsing config.json: %s", err)
	}
	events, err := internal.LoadEvents(os.Getenv("EVENTS_PATH"))
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
