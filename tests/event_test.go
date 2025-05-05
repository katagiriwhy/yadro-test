package tests

import (
	"os"
	"testing"
	"yadro-test/events"
)

func TestLoadEvents(t *testing.T) {
	data := "[09:31:49.285] 1 3\n[09:32:17.531] 1 2\n[09:37:47.892] 1 5\n[09:38:28.673] 1 1\n[09:39:25.079] 1 4\n[09:55:00.000] 2 1 10:00:00.000\n[09:56:30.000] 2 2 10:01:30.000\n[09:58:00.000] 2 3 10:03:00.000\n[09:59:30.000] 2 4 10:04:30.000\n[09:59:45.000] 3 1"
	tmpFile, err := os.CreateTemp("", "events")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write([]byte(data)); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}
	evnts, err := events.LoadEvents(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if len(evnts) == 0 {
		t.Fatal("no events")
	}
}
