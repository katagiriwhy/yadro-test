package tests

import (
	"os"
	"testing"
)

func TestParseConfig(t *testing.T) {
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
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write([]byte(data)); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

}
