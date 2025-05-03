package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	Laps        int    `json:"laps"`
	LapLen      int    `json:"lapLen"`
	PenaltyLen  int    `json:"penaltyLen"`
	FiringLines int    `json:"firingLines"`
	Start       string `json:"start"`
	StartDelta  string `json:"startDelta"`
}

type Event struct {
	Time         time.Time `json:"time"`
	EventID      uint      `json:"eventID"`
	CompetitorID uint      `json:"competitorID"`
	ExtraParams  string    `json:"extraParams"`
}

type Competitor struct {
	ID             int
	Registered     bool
	ScheduledStart time.Time
	ActualStart    time.Time
	Finished       bool
	Disqualified   bool
	NotFinished    bool
	Comment        string
	LapTimes       []time.Duration
	PenaltyTime    time.Duration
	Hits           int
	Shots          int
	CurrentLap     int
	OnFiringRange  bool
	OnPenaltyLap   bool
	LastEventTime  time.Time
}

func ParseConfig(path string) (*Config, error) {
	var Config Config

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &Config)
	if err != nil {
		return nil, err
	}
	return &Config, nil
}
