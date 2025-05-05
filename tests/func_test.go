package tests

import (
	"testing"
	"time"
	"yadro-test/utils"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		data     time.Duration
		expected string
	}{
		{time.Millisecond * 1, "00:00:00:001"},
		{time.Hour + time.Minute*30 + time.Second*15 + time.Millisecond*123, "01:30:15:123"},
		{time.Hour*2 + time.Minute*32, "03:32:00:000"},
	}
	for _, test := range tests {
		res := utils.FormatDuration(test.data)
		if res != test.expected {
			t.Errorf("FormatDuration(%d): expected %s, got %s", test.data, test.expected, res)
		}
	}
}
