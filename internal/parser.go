package internal

import (
	"fmt"
	"os"
	"strings"
	"yadro-test/config"
	"yadro-test/utils"
)

func ParseEvent(line string) (config.Event, error) {
	var e config.Event

	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return e, fmt.Errorf("invalid event format")
	}

	timeStr := strings.Trim(parts[0], "[]")
	t, err := utils.ParseTime(timeStr)
	if err != nil {
		return e, err
	}
	e.Time = t

	_, err = fmt.Sscanf(parts[1], "%d", &e.EventID)
	if err != nil {
		return e, err
	}

	_, err = fmt.Sscanf(parts[2], "%d", &e.CompetitorID)
	if err != nil {
		return e, err
	}

	if len(parts) > 3 {
		e.ExtraParams = strings.Join(parts[3:], " ")
	}

	return e, nil
}

func LoadEvents(filePath string) ([]config.Event, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var events []config.Event
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		event, err := ParseEvent(line)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
