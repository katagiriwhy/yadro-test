package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"yadro-test/internal"
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
	EventID      int       `json:"eventID"`
	CompetitorID int       `json:"competitorID"`
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

type Competition struct {
	Config      Config
	Competitors map[int]*Competitor
	Events      []Event
	StartTime   time.Time
	Delta       time.Duration
}

func CreateCompetition(config *Config) (*Competition, error) {
	startTime, err := internal.ParseTime(config.Start)
	if err != nil {
		return nil, err
	}
	delta, err := time.ParseDuration(config.StartDelta)
	if err != nil {
		return nil, err
	}

	return &Competition{
		Config:      *config,
		StartTime:   startTime,
		Delta:       delta,
		Competitors: make(map[int]*Competitor),
	}, nil
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

var (
	penaltyStart time.Time
)

func (c *Competition) ProcessEvent(event Event) string {
	competitorID := event.CompetitorID
	competitor, exist := c.Competitors[competitorID]
	if !exist {
		competitor = &Competitor{
			ID: competitorID,
		}
		c.Competitors[competitorID] = competitor
	}

	competitor.LastEventTime = event.Time

	var output string

	switch event.EventID {
	case 1:
		competitor.Registered = true
		output = fmt.Sprintf("[%s] The competitor(%d) registered", internal.FormatTime(event.Time), competitor.ID)

	case 2:
		starTime, err := internal.ParseTime(event.ExtraParams)
		if err != nil {
			return fmt.Sprintf("Error occurred: %s", err.Error())
		}
		competitor.ScheduledStart = starTime
		output = fmt.Sprintf("[%s] The start time for the competitor(%d) was set by a draw to %s", internal.FormatTime(event.Time), competitorID, event.ExtraParams)
	case 3:
		output = fmt.Sprintf("[%s] The competitor(%d) is on the start line", internal.FormatTime(event.Time), competitor.ID)
	case 4:
		competitor.ActualStart = event.Time
		if competitor.ActualStart.After(competitor.ScheduledStart.Add(c.Delta)) {
			competitor.Disqualified = true
			output = fmt.Sprintf("[%s] The competitor(%d) is disqualified", internal.FormatTime(event.Time), competitor.ID)
			c.Events = append(c.Events, Event{
				EventID:      32,
				CompetitorID: competitor.ID,
				Time:         event.Time,
			})
		} else {
			output = fmt.Sprintf("[%s] The competitor(%d) has started", internal.FormatTime(event.Time), competitor.ID)
		}
	case 5:
		competitor.OnFiringRange = true
		output = fmt.Sprintf("[%s] The competitor(%d) is on the firing range(%s)", internal.FormatTime(event.Time), competitorID, event.ExtraParams)
	case 6:
		competitor.Shots++
		competitor.Hits++
		output = fmt.Sprintf("[%s] The target(%s) has been hit by competitor(%d)", internal.FormatTime(event.Time), event.ExtraParams, competitor.ID)
	case 7:
		competitor.OnFiringRange = false
		output = fmt.Sprintf("[%s] The competitor(%d) left the firing range", internal.FormatTime(event.Time), competitorID)

	case 8:
		competitor.OnPenaltyLap = true
		penaltyStart = event.Time
		output = fmt.Sprintf("[%s] The competitor(%d) entered the penalty laps", internal.FormatTime(event.Time), competitorID)

	case 9:
		competitor.OnPenaltyLap = false
		penaltyEnd := event.Time
		competitor.PenaltyTime = penaltyEnd.Sub(penaltyStart)
		output = fmt.Sprintf("[%s] The competitor(%d) left the penalty laps", internal.FormatTime(event.Time), competitorID)
	case 10:
		competitor.CurrentLap++
		if competitor.CurrentLap == 1 {
			lapTime := event.Time.Sub(competitor.ActualStart)
			competitor.LapTimes = append(competitor.LapTimes, lapTime)
		} else if competitor.CurrentLap <= c.Config.Laps {
			lapTime := event.Time.Sub(competitor.ActualStart)
			competitor.LapTimes = append(competitor.LapTimes, lapTime)
		}

		if competitor.CurrentLap == c.Config.Laps {
			competitor.Finished = true
			c.Events = append(c.Events, Event{
				EventID:      33,
				CompetitorID: competitor.ID,
				Time:         event.Time,
			})
		}

	case 11:
		competitor.NotFinished = true
		competitor.Comment = event.ExtraParams
		output = fmt.Sprintf("[%s] The competitor(%d) can`t continue: %s", internal.FormatTime(event.Time), competitor.ID, event.ExtraParams)
	}

	return output
}

func (c *Competitor) MakeReport() string {

}
