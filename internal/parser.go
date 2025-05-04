package internal

import (
	"fmt"
	"os"
	"strings"
	"time"
	"yadro-test/config"
)

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04:05.000", timeStr)
}

func FormatTime(t time.Time) string {
	return t.Format("15:04:05.000")
}

func ParseDuration(durStr string) (time.Duration, error) {
	parts := strings.Split(durStr, ".")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid duration format")
	}

	hms := parts[0]
	ms := parts[1]

	t, err := time.Parse("15:04:05", hms)
	if err != nil {
		return 0, err
	}

	parsedMs, err := parseMilliseconds(ms)
	if err != nil {
		return 0, err
	}

	dur := time.Hour*time.Duration(t.Hour()) +
		time.Minute*time.Duration(t.Minute()) +
		time.Second*time.Duration(t.Second()) +
		time.Millisecond*time.Duration(parsedMs)

	return dur, nil
}

func parseMilliseconds(ms string) (int, error) {
	var val int
	_, err := fmt.Sscanf(ms, "%d", &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func ParseEvent(line string) (config.Event, error) {
	var e config.Event

	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return e, fmt.Errorf("invalid event format")
	}

	timeStr := strings.Trim(parts[0], "[]")
	t, err := ParseTime(timeStr)
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
