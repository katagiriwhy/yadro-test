package config

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
	"yadro-test/utils"
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
	ID               int
	Registered       bool
	ScheduledStart   time.Time
	ActualStart      time.Time
	Finished         bool
	Disqualified     bool
	NotFinished      bool
	Comment          string
	LapTimes         []time.Duration
	PenaltyTime      time.Duration
	PenaltyStartTime time.Time
	shotsCounter     []int
	Hits             int
	Shots            int
	CurrentLap       int
	OnFiringRange    bool
	OnPenaltyLap     bool
	LastEventTime    time.Time
}

type Competition struct {
	Config      Config
	Competitors map[int]*Competitor
	Events      []Event
	StartTime   time.Time
	Delta       time.Duration
}

func CreateCompetition(config *Config) (*Competition, error) {
	startTime, err := utils.ParseTime(config.Start)
	if err != nil {
		return nil, err
	}
	delta, err := utils.ParseDurationCfg(config.StartDelta)
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
		output = fmt.Sprintf("[%s] The competitor(%d) registered", utils.FormatTime(event.Time), competitor.ID)

	case 2:
		starTime, err := utils.ParseTime(event.ExtraParams)
		if err != nil {
			return fmt.Sprintf("Error occurred: %s", err.Error())
		}
		competitor.ScheduledStart = starTime
		output = fmt.Sprintf("[%s] The start time for the competitor(%d) was set by a draw to %s", utils.FormatTime(event.Time), competitorID, event.ExtraParams)
	case 3:
		if event.Time.After(competitor.ScheduledStart) {
			competitor.Disqualified = true
			return fmt.Sprintf("The competitor(%d) is disqualified, because of the delay", competitor.ID)
		}
		output = fmt.Sprintf("[%s] The competitor(%d) is on the start line", utils.FormatTime(event.Time), competitor.ID)
	case 4:
		competitor.ActualStart = event.Time
		if competitor.ActualStart.After(competitor.ScheduledStart.Add(c.Delta)) {
			competitor.Disqualified = true
			output = fmt.Sprintf("[%s] The competitor(%d) is disqualified", utils.FormatTime(event.Time), competitor.ID)
			c.Events = append(c.Events, Event{
				EventID:      32,
				CompetitorID: competitor.ID,
				Time:         event.Time,
			})
		} else {
			output = fmt.Sprintf("[%s] The competitor(%d) has started", utils.FormatTime(event.Time), competitor.ID)
		}
	case 5:
		competitor.OnFiringRange = true
		output = fmt.Sprintf("[%s] The competitor(%d) is on the firing range(%s)", utils.FormatTime(event.Time), competitorID, event.ExtraParams)
	case 6:
		shot, err := strconv.Atoi(event.ExtraParams)
		if err != nil {
			return fmt.Sprintf("Error occurred: %s", err.Error())
		}
		competitor.shotsCounter = append(competitor.shotsCounter, shot)
		competitor.Hits++
		output = fmt.Sprintf("[%s] The target(%s) has been hit by competitor(%d)", utils.FormatTime(event.Time), event.ExtraParams, competitor.ID)
	case 7:
		competitor.OnFiringRange = false
		if len(competitor.shotsCounter) > 0 {
			competitor.Shots += slices.Max(competitor.shotsCounter)
		}
		output = fmt.Sprintf("[%s] The competitor(%d) left the firing range", utils.FormatTime(event.Time), competitorID)

	case 8:
		competitor.OnPenaltyLap = true
		competitor.PenaltyStartTime = event.Time
		output = fmt.Sprintf("[%s] The competitor(%d) entered the penalty laps", utils.FormatTime(event.Time), competitorID)

	case 9:
		competitor.OnPenaltyLap = false
		penaltyEnd := event.Time
		competitor.PenaltyTime += penaltyEnd.Sub(competitor.PenaltyStartTime)
		output = fmt.Sprintf("[%s] The competitor(%d) left the penalty laps", utils.FormatTime(event.Time), competitorID)
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
		output = fmt.Sprintf("[%s] The competitor(%d) can`t continue: %s", utils.FormatTime(event.Time), competitor.ID, event.ExtraParams)
	}

	return output
}

func (c *Competition) MakeReport() string {
	var b strings.Builder

	competitors := make([]*Competitor, 0, len(c.Competitors))
	for _, competitor := range c.Competitors {
		competitors = append(competitors, competitor)
	}

	sort.Slice(competitors, func(i, j int) bool {
		if competitors[i].Disqualified {
			return false
		}
		if competitors[j].Disqualified {
			return true
		}
		if competitors[i].NotFinished && !competitors[j].NotFinished {
			return false
		}
		if !competitors[i].Finished && competitors[j].Finished {
			return true
		}
		compTotalTime1 := competitors[i].LastEventTime.Sub(competitors[i].ActualStart)
		compTotalTime2 := competitors[j].LastEventTime.Sub(competitors[j].ActualStart)
		return compTotalTime1 < compTotalTime2
	})

	for _, competitor := range competitors {
		if competitor.Disqualified {
			b.WriteString(fmt.Sprintf("[NotStarted] %d ", competitor.ID))
		} else if competitor.NotFinished {
			b.WriteString(fmt.Sprintf("[NotFinished] %d ", competitor.ID))
		} else {
			totalTime := competitor.LastEventTime.Sub(competitor.ActualStart)
			b.WriteString(fmt.Sprintf("[%s] %d ", utils.FormatDuration(totalTime), competitor.ID))
		}

		b.WriteString(fmt.Sprintf("[{"))
		for i, lapTime := range competitor.LapTimes {
			if i > 0 {
				b.WriteString(", {")
			}
			if i < c.Config.Laps {
				b.WriteString(utils.FormatDuration(lapTime))
				b.WriteString(fmt.Sprintf(", %.3f}", calculateSpeed(c.Config.LapLen, lapTime)))
			} else {
				b.WriteString(",")
			}
		}
		b.WriteString("] ")
		if competitor.PenaltyTime > 0 {
			b.WriteString("{" + utils.FormatDuration(competitor.PenaltyTime))
			b.WriteString(fmt.Sprintf(", %.3f} ", calculateSpeed(c.Config.PenaltyLen, competitor.PenaltyTime)))
		}

		b.WriteString(fmt.Sprintf("%d/%d\n", competitor.Hits, competitor.Shots))
	}

	return b.String()
}

func calculateSpeed(distance int, duration time.Duration) float64 {
	if distance == 0 || duration == 0 {
		return 0
	}

	return float64(distance) / (duration.Seconds())
}
