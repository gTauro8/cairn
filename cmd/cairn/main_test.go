package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestAddWithoutTagsOmitsTagsKey(t *testing.T) {
	inTempDir(t)

	if err := cmdAdd([]string{"testo"}); err != nil {
		t.Fatalf("cmdAdd() error = %v", err)
	}

	line := readSingleLogLine(t, logFile)
	var got map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if _, ok := got["tags"]; ok {
		t.Fatalf("entry contains tags key: %s", line)
	}
}

func TestAddNormalizesAndDeduplicatesTags(t *testing.T) {
	inTempDir(t)

	if err := cmdAdd([]string{"--tags", "Decision, go, decision", "testo"}); err != nil {
		t.Fatalf("cmdAdd() error = %v", err)
	}

	line := readSingleLogLine(t, logFile)
	var got entry
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	want := []string{"decision", "go"}
	if !reflect.DeepEqual(got.Tags, want) {
		t.Fatalf("tags = %v, want %v", got.Tags, want)
	}
}

func TestLogFiltersByTagAndExcludesEntriesWithoutTags(t *testing.T) {
	inTempDir(t)
	writeLog(t, `{"id":"old","ts":"2026-07-10T10:00:00Z","text":"senza tag"}`+"\n"+
		`{"id":"match","ts":"2026-07-10T11:00:00Z","text":"inclusa","tags":["x"]}`+"\n"+
		`{"id":"other","ts":"2026-07-10T12:00:00Z","text":"esclusa","tags":["y"]}`+"\n")

	out := captureStdout(t, func() {
		if err := cmdLog([]string{"--tag", "x"}); err != nil {
			t.Fatalf("cmdLog() error = %v", err)
		}
	})

	if !strings.Contains(out, "inclusa") {
		t.Fatalf("output does not contain matching entry: %q", out)
	}
	if strings.Contains(out, "senza tag") || strings.Contains(out, "esclusa") {
		t.Fatalf("output contains non-matching entries: %q", out)
	}
}

func TestAddRejectsFlagAfterTextWithoutWriting(t *testing.T) {
	inTempDir(t)

	if err := cmdAdd([]string{"testo", "--tags", "a,b"}); err == nil {
		t.Fatal("cmdAdd() error = nil, want error")
	}
	if _, err := os.Stat(logFile); !os.IsNotExist(err) {
		t.Fatalf("log file exists or stat returned unexpected error: %v", err)
	}
}

func TestAddRequiresText(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
	}{
		{name: "no arguments"},
		{name: "tags only", args: []string{"--tags", "x"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			inTempDir(t)

			if err := cmdAdd(tt.args); err == nil {
				t.Fatal("cmdAdd() error = nil, want error")
			}
			if _, err := os.Stat(logFile); !os.IsNotExist(err) {
				t.Fatalf("log file exists or stat returned unexpected error: %v", err)
			}
		})
	}
}

func TestStateDirectoryIsLimitedToExactWorkingDirectory(t *testing.T) {
	parent := t.TempDir()
	parentState := filepath.Join(parent, stateDir)
	if err := os.Mkdir(parentState, 0o755); err != nil {
		t.Fatalf("os.Mkdir() error = %v", err)
	}
	parentLog := filepath.Join(parent, logFile)
	if err := os.WriteFile(parentLog, []byte(`{"id":"parent","ts":"2026-07-10T10:00:00Z","text":"dal padre"}`+"\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	child := filepath.Join(parent, "child")
	if err := os.Mkdir(child, 0o755); err != nil {
		t.Fatalf("os.Mkdir() error = %v", err)
	}
	chdir(t, child)

	out := captureStdout(t, func() {
		if err := cmdLog(nil); err != nil {
			t.Fatalf("cmdLog() error = %v", err)
		}
	})
	if out != "" {
		t.Fatalf("cmdLog() walked up to parent, output = %q", out)
	}

	if err := cmdAdd([]string{"nel figlio"}); err != nil {
		t.Fatalf("cmdAdd() error = %v", err)
	}
	childLine := readSingleLogLine(t, filepath.Join(child, logFile))
	if !strings.Contains(childLine, `"text":"nel figlio"`) {
		t.Fatalf("child log = %q", childLine)
	}
	if got := readSingleLogLine(t, parentLog); strings.Contains(got, "nel figlio") {
		t.Fatalf("cmdAdd() wrote to parent log: %q", got)
	}
}

func inTempDir(t *testing.T) {
	t.Helper()
	chdir(t, t.TempDir())
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("os.Chdir(%q) error = %v", dir, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(old); err != nil {
			t.Errorf("restore cwd: %v", err)
		}
	})
}

func readSingleLogLine(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile(%q) error = %v", path, err)
	}
	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	if len(lines) != 1 {
		t.Fatalf("log has %d lines, want 1: %q", len(lines), b)
	}
	return lines[0]
}

func writeLog(t *testing.T, contents string) {
	t.Helper()
	if err := os.Mkdir(stateDir, 0o755); err != nil {
		t.Fatalf("os.Mkdir() error = %v", err)
	}
	if err := os.WriteFile(logFile, []byte(contents), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = old
		r.Close()
	}()

	fn()
	if err := w.Close(); err != nil {
		t.Fatalf("stdout writer close: %v", err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("io.ReadAll() error = %v", err)
	}
	return string(b)
}
