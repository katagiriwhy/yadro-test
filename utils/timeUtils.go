package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes())
	seconds := int(d.Seconds())
	ms := int(d.Milliseconds())
	return fmt.Sprintf("%02d:%02d:%02d:%03d", hours, minutes, seconds, ms)
}

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04:05.000", timeStr)
}

func FormatTime(t time.Time) string {
	return t.Format("15:04:05.000")
}

func ParseDurationCfg(durStr string) (time.Duration, error) {
	parts := strings.Split(durStr, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid duration format: %s", durStr)
	}
	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, err
	}
	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Second*time.Duration(seconds), nil
}

func ParseDuration(durStr string) (time.Duration, error) {
	parts := strings.Split(durStr, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid duration format")
	}
	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	secondsAndMs := strings.Split(parts[2], ".")

	seconds, err := strconv.Atoi(secondsAndMs[0])
	if err != nil {
		return 0, err
	}

	ms, err := strconv.Atoi(secondsAndMs[1])
	if err != nil {
		return 0, err
	}

	dur := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second + time.Millisecond*time.Duration(ms)

	return dur, nil

}
