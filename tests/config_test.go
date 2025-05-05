package tests

import (
	"os"
	"testing"
	config2 "yadro-test/config"
)

func TestParseConfig(t *testing.T) {
	t.Parallel()

	data := `{
	"laps": 2,
  "lapLen": 3500,
  "penaltyLen": 150,
  "firingLines": 2,
  "start": "10:00:00.000",
  "startDelta": "00:01:30"
}`
	tmpFile, err := os.CreateTemp("", "test_config.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temp file: %v", err)
		}
	}(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(data)); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	config, err := config2.ParseConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}
	if config.Laps != 2 {
		t.Fatalf("Laps should be 2, got %d", config.Laps)
	}
	if config.PenaltyLen != 150 {
		t.Fatalf("Penalty length should be 150, got %d", config.PenaltyLen)
	}
	if config.FiringLines != 2 {
		t.Fatalf("Firing lines should be 2, got %d", config.FiringLines)
	}

}
