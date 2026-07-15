package main

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
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
	if _, ok := got["files"]; ok {
		t.Fatalf("entry contains files key: %s", line)
	}
	if string(got["source"]) != `"manual"` {
		t.Fatalf("source = %s, want manual", got["source"])
	}
}

func TestAddRecordsExplicitProvenance(t *testing.T) {
	inTempDir(t)

	if err := cmdAdd([]string{"--source", "git", "--commit", "abc123", "testo"}); err != nil {
		t.Fatalf("cmdAdd() error = %v", err)
	}

	line := readSingleLogLine(t, logFile)
	var got entry
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if got.Source != "git" || got.Commit != "abc123" {
		t.Fatalf("provenance = (%q, %q), want (git, abc123)", got.Source, got.Commit)
	}
}

func TestAddParsesAndDeduplicatesFilesPreservingCase(t *testing.T) {
	inTempDir(t)

	if err := cmdAdd([]string{"--files", " a.go, B.go ,a.go", "testo"}); err != nil {
		t.Fatalf("cmdAdd() error = %v", err)
	}

	line := readSingleLogLine(t, logFile)
	var got entry
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	want := []string{"a.go", "B.go"}
	if !reflect.DeepEqual(got.Files, want) {
		t.Fatalf("files = %v, want %v", got.Files, want)
	}
}

func TestLogFiltersByFileAndExcludesEntriesWithoutFiles(t *testing.T) {
	inTempDir(t)
	writeLog(t, `{"id":"old","ts":"2026-07-10T10:00:00Z","text":"senza file"}`+"\n"+
		`{"id":"match","ts":"2026-07-10T11:00:00Z","text":"inclusa","files":["a.go"]}`+"\n"+
		`{"id":"other","ts":"2026-07-10T12:00:00Z","text":"esclusa","files":["b.go"]}`+"\n")

	out := captureStdout(t, func() {
		if err := cmdLog([]string{"--file", "a.go"}); err != nil {
			t.Fatalf("cmdLog() error = %v", err)
		}
	})

	if !strings.Contains(out, "inclusa") {
		t.Fatalf("output does not contain matching entry: %q", out)
	}
	if strings.Contains(out, "senza file") || strings.Contains(out, "esclusa") {
		t.Fatalf("output contains non-matching entries: %q", out)
	}
}

func TestLogJSONPreservesCompleteMatchingEntry(t *testing.T) {
	inTempDir(t)
	want := `{"id":"match","ts":"2026-07-10T11:00:00Z","text":"inclusa","tags":["x"],"files":["a.go"],"source":"git"}`
	writeLog(t, want+"\n"+
		`{"id":"other","ts":"2026-07-10T12:00:00Z","text":"esclusa","tags":["y"]}`+"\n")

	out := captureStdout(t, func() {
		if err := cmdLog([]string{"--tag", "x", "--json"}); err != nil {
			t.Fatalf("cmdLog() error = %v", err)
		}
	})

	if out != want+"\n" {
		t.Fatalf("output = %q, want %q", out, want+"\n")
	}
}

func TestCheckAcceptsValidLogAndExistingFiles(t *testing.T) {
	inTempDir(t)
	if err := os.WriteFile("a.go", []byte("package a\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	writeLog(t, `{"id":"ok","ts":"2026-07-10T11:00:00Z","text":"valida","files":["a.go"]}`+"\n")

	out := captureStdout(t, func() {
		if err := cmdCheck(nil); err != nil {
			t.Fatalf("cmdCheck() error = %v", err)
		}
	})
	if out != "" {
		t.Fatalf("output = %q, want empty", out)
	}
}

func TestCheckReportsInvalidEntriesAndMissingFiles(t *testing.T) {
	inTempDir(t)
	writeLog(t, "not-json\n"+
		`{"id":"","ts":"2026-07-10T11:00:00Z","text":"","files":["missing.go"]}`+"\n")

	var checkErr error
	out := captureStdout(t, func() {
		checkErr = cmdCheck(nil)
	})
	if checkErr == nil {
		t.Fatal("cmdCheck() error = nil, want error")
	}
	for _, want := range []string{
		"linea 1: JSON non valido",
		`linea 2: campo obbligatorio "id" mancante o vuoto`,
		`linea 2: campo obbligatorio "text" mancante o vuoto`,
		"linea 2: file referenziato non trovato: missing.go",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output = %q, want substring %q", out, want)
		}
	}
	if !strings.Contains(checkErr.Error(), "4 problemi trovati") {
		t.Fatalf("error = %q, want issue count", checkErr)
	}
}

func TestCheckRequiresCommitForGitSource(t *testing.T) {
	inTempDir(t)
	writeLog(t, `{"id":"git","ts":"2026-07-10T11:00:00Z","text":"senza hash","source":"Git"}`+"\n")

	var checkErr error
	out := captureStdout(t, func() {
		checkErr = cmdCheck(nil)
	})
	if checkErr == nil {
		t.Fatal("cmdCheck() error = nil, want error")
	}
	if !strings.Contains(out, `source "Git" richiede il campo commit`) {
		t.Fatalf("output = %q, want missing commit diagnostic", out)
	}
}

func TestHookInstallIsSafeAndIdempotent(t *testing.T) {
	repo := initGitRepo(t)
	chdir(t, repo)

	for i := 0; i < 2; i++ {
		out := captureStdout(t, func() {
			if err := cmdHook([]string{"install"}); err != nil {
				t.Fatalf("cmdHook(install) error = %v", err)
			}
		})
		if !strings.Contains(out, "hook post-commit installato") {
			t.Fatalf("output = %q, want install confirmation", out)
		}
	}

	hookPath := filepath.Join(repo, ".githooks", "post-commit")
	got, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	if string(got) != installedHook {
		t.Fatalf("installed hook differs from canonical hook")
	}
	info, err := os.Stat(hookPath)
	if err != nil {
		t.Fatalf("os.Stat() error = %v", err)
	}
	if info.Mode().Perm()&0o111 == 0 {
		t.Fatalf("hook mode = %v, want executable", info.Mode().Perm())
	}
	if got := strings.TrimSpace(runGitTest(t, repo, "config", "--local", "--get", "core.hooksPath")); got != ".githooks" {
		t.Fatalf("core.hooksPath = %q, want .githooks", got)
	}
}

func TestHookInstallRefusesExistingHooksPath(t *testing.T) {
	repo := initGitRepo(t)
	runGitTest(t, repo, "config", "--local", "core.hooksPath", "custom-hooks")
	chdir(t, repo)

	if err := cmdHook([]string{"install"}); err == nil {
		t.Fatal("cmdHook(install) error = nil, want conflict")
	}
	if _, err := os.Stat(filepath.Join(repo, ".githooks")); !os.IsNotExist(err) {
		t.Fatalf(".githooks created despite conflict: %v", err)
	}
}

func TestHookInstallRefusesDifferentPostCommit(t *testing.T) {
	repo := initGitRepo(t)
	hookDir := filepath.Join(repo, ".githooks")
	if err := os.Mkdir(hookDir, 0o755); err != nil {
		t.Fatalf("os.Mkdir() error = %v", err)
	}
	hookPath := filepath.Join(hookDir, "post-commit")
	if err := os.WriteFile(hookPath, []byte("#!/bin/sh\necho custom\n"), 0o755); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	chdir(t, repo)

	if err := cmdHook([]string{"install"}); err == nil {
		t.Fatal("cmdHook(install) error = nil, want conflict")
	}
	got, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	if string(got) != "#!/bin/sh\necho custom\n" {
		t.Fatalf("custom hook was overwritten: %q", got)
	}
}

func TestHookRunCapturesMarkedCommit(t *testing.T) {
	repo := initGitRepo(t)
	file := filepath.Join(repo, "a.go")
	if err := os.WriteFile(file, []byte("package a\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	runGitTest(t, repo, "add", "a.go")
	runGitTest(t, repo, "commit", "-q", "-m", "initial")
	if err := os.WriteFile(file, []byte("package a\n\nconst A = 1\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	runGitTest(t, repo, "add", "a.go")
	runGitTest(t, repo, "commit", "-q", "-m", "record decision\n\nCairn-Note: true\nCairn-Tags: Decision\n\nCo-Authored-By: Test <test@example.invalid>")
	wantCommit := strings.TrimSpace(runGitTest(t, repo, "rev-parse", "HEAD"))
	chdir(t, repo)

	if err := cmdHook([]string{"run"}); err != nil {
		t.Fatalf("cmdHook(run) error = %v", err)
	}
	line := readSingleLogLine(t, logFile)
	var got entry
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if got.Source != "git" || got.Commit != wantCommit {
		t.Fatalf("provenance = (%q, %q), want (git, %q)", got.Source, got.Commit, wantCommit)
	}
	if !reflect.DeepEqual(got.Tags, []string{"decision", "git"}) {
		t.Fatalf("tags = %v, want [decision git]", got.Tags)
	}
	if !reflect.DeepEqual(got.Files, []string{"a.go"}) {
		t.Fatalf("files = %v, want [a.go]", got.Files)
	}
}

func TestVersionDefaultsToDev(t *testing.T) {
	out := captureStdout(t, func() {
		if err := cmdVersion(nil); err != nil {
			t.Fatalf("cmdVersion() error = %v", err)
		}
	})
	if out != "cairn dev\n" {
		t.Fatalf("output = %q, want cairn dev", out)
	}
}

func TestVersionedHookMatchesCanonicalHook(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller() failed")
	}
	hookPath := filepath.Join(filepath.Dir(currentFile), "..", "..", ".githooks", "post-commit")
	got, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	if string(got) != installedHook {
		t.Fatal("versioned hook differs from installedHook")
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

func initGitRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	runGitTest(t, repo, "init", "-q")
	runGitTest(t, repo, "config", "user.name", "Cairn Test")
	runGitTest(t, repo, "config", "user.email", "cairn-test@example.invalid")
	return repo
}

func runGitTest(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}
