package main

import (
	"fmt"
	"yadro-test/config"
	"yadro-test/utils"
)

func main() {
	cfg, err := config.ParseConfig("files/cfg.json")
	if err != nil {
		fmt.Printf("Error parsing config: %s\n", err)
	}
	fmt.Println(*cfg)
	//events, err := internal.LoadEvents("files/events")
	//if err != nil {
	//	fmt.Printf("Error loading events: %s\n", err)
	//}

	duration, err := utils.ParseDuration("15:04:05.012")
	fmt.Println(duration)
}
