package tests

import (
	"testing"
	"time"
	"yadro-test/utils"
)

func TestFormatDuration(t *testing.T) {
	t.Parallel()
	tests := []struct {
		data     time.Duration
		expected string
	}{
		{time.Millisecond * 1, "00:00:00:001"},
		{time.Hour + time.Minute*30 + time.Second*15 + time.Millisecond*123, "01:30:15:123"},
	}
	for _, test := range tests {
		res := utils.FormatDuration(test.data)
		if res != test.expected {
			t.Errorf("FormatDuration(%d): expected %s, got %s", test.data, test.expected, res)
		}
	}
}

func TestParseTime(t *testing.T) {
	t.Parallel()
	tests := []struct {
		data     string
		expected time.Time
	}{
		{"09:30:00.000", time.Date(0, 1, 1, 9, 30, 0, 0, time.UTC)},
		{"23:59:59.999", time.Date(0, 1, 1, 23, 59, 59, 999000000, time.UTC)},
		{"09:21:33:000", time.Time{}},
	}
	for _, test := range tests {
		res, err := utils.ParseTime(test.data)
		if err != nil {
			t.Error(err)
		}
		if res != test.expected {
			t.Errorf("ParseTime(%s): expected %s, got %s", test.data, test.expected, res)
		}
	}
}

func TestFormatTime(t *testing.T) {
	t.Parallel()
	tests := []struct {
		data     time.Time
		expected string
	}{
		{time.Date(0, 1, 1, 9, 30, 0, 0, time.UTC), "09:30:00.000"},
		{time.Date(0, 1, 1, 23, 59, 59, 999000000, time.UTC), "23:59:59.999"},
	}
	for _, test := range tests {
		res := utils.FormatTime(test.data)
		if res != test.expected {
			t.Errorf("FormatTime(%s): expected %s, got %s", test.data, test.expected, res)
		}
	}
}

func TestParseDurationCfg(t *testing.T) {
	t.Parallel()
	tests := []struct {
		data     string
		expected time.Duration
	}{
		{"09:15:15", time.Duration(9)*time.Hour + time.Duration(15)*time.Minute + time.Duration(15)*time.Second},
		{"12:33:00", time.Duration(12)*time.Hour + time.Duration(33)*time.Minute},
		{"DAS123.324-00k", time.Duration(0)},
		{"dasd:12:32", time.Duration(0)},
		{"03:dasd:12", time.Duration(0)},
		{"03:12:dasd", time.Duration(0)},
	}
	for _, test := range tests {
		res, err := utils.ParseDurationCfg(test.data)
		if err != nil {
			t.Error(err)
		}
		if res != test.expected {
			t.Errorf("ParseDurationCfg(%s): expected %s, got %s", test.data, test.expected, res)
		}
	}
}
