package events

import (
	"fmt"
	"os"
	"strings"
	"time"
	"yadro-test/utils"
)

type Event struct {
	Time         time.Time `json:"time"`
	EventID      int       `json:"eventID"`
	CompetitorID int       `json:"competitorID"`
	ExtraParams  string    `json:"extraParams"`
}

func LoadEvents(filePath string) ([]Event, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var evnts []Event
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		event, err := ParseEvent(line)
		if err != nil {
			return nil, err
		}
		evnts = append(evnts, event)
	}
	return evnts, nil
}

func ParseEvent(line string) (Event, error) {
	var e Event

	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return e, fmt.Errorf("invalid events format")
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
