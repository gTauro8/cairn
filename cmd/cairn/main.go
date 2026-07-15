package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const stateDir = ".cairn"
const logFile = stateDir + "/log.jsonl"

var version = "dev"

const installedHook = `#!/bin/sh
# Installato da "cairn hook install". Non deve mai far fallire il commit.
CAIRN_BIN=""
if command -v cairn >/dev/null 2>&1; then
    CAIRN_BIN="cairn"
else
    REPO_ROOT=$(git rev-parse --show-toplevel 2>/dev/null)
    if [ -n "$REPO_ROOT" ] && [ -x "$REPO_ROOT/cairn" ]; then
        CAIRN_BIN="$REPO_ROOT/cairn"
    fi
fi

if [ -z "$CAIRN_BIN" ]; then
    echo "post-commit: cairn non trovato (ne' su PATH ne' in <repo>/cairn) - nota non catturata" >&2
    exit 0
fi

"$CAIRN_BIN" hook run || echo "post-commit: cairn hook run fallito - nota non catturata" >&2
exit 0
`

type entry struct {
	ID     string   `json:"id"`
	TS     string   `json:"ts"`
	Text   string   `json:"text"`
	Tags   []string `json:"tags,omitempty"`
	Files  []string `json:"files,omitempty"`
	Source string   `json:"source,omitempty"`
	Commit string   `json:"commit,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "add":
		err = cmdAdd(os.Args[2:])
	case "log":
		err = cmdLog(os.Args[2:])
	case "check":
		err = cmdCheck(os.Args[2:])
	case "hook":
		err = cmdHook(os.Args[2:])
	case "version":
		err = cmdVersion(os.Args[2:])
	default:
		usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "cairn:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: cairn add [--tags a,b,c] [--files a,b] [--source s] [--commit hash] <text> | cairn log [--tag x] [--file p] [--json] | cairn check | cairn hook install | cairn version")
}

func cmdAdd(args []string) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	tagsFlag := fs.String("tags", "", "tag separati da virgola")
	filesFlag := fs.String("files", "", "file del repo a cui la nota si riferisce, separati da virgola")
	sourceFlag := fs.String("source", "manual", "origine della nota")
	commitFlag := fs.String("commit", "", "hash del commit di origine")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if tok, ok := misplacedFlag(fs); ok {
		return fmt.Errorf("%q dopo il testo non ha effetto: mettilo prima, es. cairn add --tags a,b \"testo\"", tok)
	}

	text := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if text == "" {
		return fmt.Errorf("add richiede un testo non vuoto")
	}

	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		return fmt.Errorf("impossibile creare %s: %w", stateDir, err)
	}

	e := entry{
		ID:     genID(),
		TS:     time.Now().UTC().Format(time.RFC3339),
		Text:   text,
		Tags:   parseTags(*tagsFlag),
		Files:  parseFiles(*filesFlag),
		Source: strings.TrimSpace(*sourceFlag),
		Commit: strings.TrimSpace(*commitFlag),
	}
	line, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("impossibile serializzare la nota: %w", err)
	}

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("impossibile aprire %s: %w", logFile, err)
	}
	defer f.Close()

	if _, err := f.Write(append(line, '\n')); err != nil {
		return fmt.Errorf("impossibile scrivere su %s: %w", logFile, err)
	}

	return nil
}

func cmdLog(args []string) error {
	fs := flag.NewFlagSet("log", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	tagFilter := fs.String("tag", "", "mostra solo le note con questo tag")
	fileFilter := fs.String("file", "", "mostra solo le note che riferiscono questo file")
	jsonOutput := fs.Bool("json", false, "emette le note come JSON Lines")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("log non accetta argomenti posizionali")
	}

	f, err := os.Open(logFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("impossibile leggere %s: %w", logFile, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var e entry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			return fmt.Errorf("riga corrotta in %s: %w", logFile, err)
		}
		if *tagFilter != "" && !hasTag(e.Tags, *tagFilter) {
			continue
		}
		if *fileFilter != "" && !hasFile(e.Files, *fileFilter) {
			continue
		}
		if *jsonOutput {
			fmt.Println(scanner.Text())
			continue
		}
		fmt.Printf("[%s] %s\n", e.TS, e.Text)
	}
	return scanner.Err()
}

func cmdCheck(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("check non accetta argomenti")
	}

	f, err := os.Open(logFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("impossibile leggere %s: %w", logFile, err)
	}
	defer f.Close()

	issues := 0
	lineNumber := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineNumber++
		var e entry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			fmt.Printf("linea %d: JSON non valido: %v\n", lineNumber, err)
			issues++
			continue
		}

		for _, required := range []struct {
			name  string
			value string
		}{
			{name: "id", value: e.ID},
			{name: "ts", value: e.TS},
			{name: "text", value: e.Text},
		} {
			if strings.TrimSpace(required.value) == "" {
				fmt.Printf("linea %d: campo obbligatorio %q mancante o vuoto\n", lineNumber, required.name)
				issues++
			}
		}

		for _, path := range e.Files {
			if _, err := os.Stat(path); err != nil {
				if os.IsNotExist(err) {
					fmt.Printf("linea %d: file referenziato non trovato: %s\n", lineNumber, path)
					issues++
					continue
				}
				return fmt.Errorf("impossibile verificare il file %q alla linea %d: %w", path, lineNumber, err)
			}
		}

		if strings.EqualFold(strings.TrimSpace(e.Source), "git") && strings.TrimSpace(e.Commit) == "" {
			fmt.Printf("linea %d: source %q richiede il campo commit\n", lineNumber, e.Source)
			issues++
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("impossibile leggere %s: %w", logFile, err)
	}
	if issues != 0 {
		return fmt.Errorf("check fallito: %d problemi trovati", issues)
	}
	return nil
}

func cmdHook(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("hook richiede un sottocomando: install")
	}

	switch args[0] {
	case "install":
		return cmdHookInstall(args[1:])
	case "run":
		return cmdHookRun(args[1:])
	default:
		return fmt.Errorf("sottocomando hook sconosciuto %q", args[0])
	}
}

func cmdHookInstall(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("hook install non accetta argomenti")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("impossibile determinare la directory corrente: %w", err)
	}
	root, err := gitOutput("rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("la directory corrente non e' una root Git: %w", err)
	}
	if canonicalPath(cwd) != canonicalPath(strings.TrimSpace(root)) {
		return fmt.Errorf("hook install va eseguito dalla root Git: %s", strings.TrimSpace(root))
	}

	hooksPath, err := localHooksPath()
	if err != nil {
		return err
	}
	if hooksPath != "" && hooksPath != ".githooks" {
		return fmt.Errorf("core.hooksPath e' gia' impostato a %q: configurazione non modificata", hooksPath)
	}

	hookPath := filepath.Join(".githooks", "post-commit")
	existing, err := os.ReadFile(hookPath)
	switch {
	case err == nil && string(existing) != installedHook:
		return fmt.Errorf("%s esiste gia' con contenuto diverso: file non sovrascritto", hookPath)
	case err != nil && !os.IsNotExist(err):
		return fmt.Errorf("impossibile leggere %s: %w", hookPath, err)
	case os.IsNotExist(err):
		if err := os.MkdirAll(filepath.Dir(hookPath), 0o755); err != nil {
			return fmt.Errorf("impossibile creare %s: %w", filepath.Dir(hookPath), err)
		}
		if err := os.WriteFile(hookPath, []byte(installedHook), 0o755); err != nil {
			return fmt.Errorf("impossibile installare %s: %w", hookPath, err)
		}
	}
	if err := os.Chmod(hookPath, 0o755); err != nil {
		return fmt.Errorf("impossibile rendere eseguibile %s: %w", hookPath, err)
	}

	if hooksPath == "" {
		if _, err := gitOutput("config", "--local", "core.hooksPath", ".githooks"); err != nil {
			return fmt.Errorf("hook scritto ma core.hooksPath non configurato: %w", err)
		}
	}

	fmt.Println("hook post-commit installato in .githooks/post-commit")
	return nil
}

func cmdHookRun(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("hook run non accetta argomenti")
	}

	message, err := gitOutput("log", "-1", "--format=%B")
	if err != nil {
		return err
	}
	noteFlag := strings.ToLower(trailerValue(message, "Cairn-Note"))
	if noteFlag != "true" && noteFlag != "yes" {
		return nil
	}

	subject, err := gitOutput("log", "-1", "--format=%s")
	if err != nil {
		return err
	}
	shortHash, err := gitOutput("log", "-1", "--format=%h")
	if err != nil {
		return err
	}
	fullHash, err := gitOutput("rev-parse", "HEAD")
	if err != nil {
		return err
	}
	subject = strings.TrimSpace(subject)
	shortHash = strings.TrimSpace(shortHash)
	fullHash = strings.TrimSpace(fullHash)

	tags := trailerValue(message, "Cairn-Tags")
	if tags == "" {
		tags = "git"
	} else {
		tags += ",git"
	}

	changed, err := gitOutput("diff-tree", "--no-commit-id", "--name-only", "-r", fullHash)
	if err != nil {
		return err
	}
	files := strings.Join(nonEmptyLines(changed), ",")
	text := fmt.Sprintf("%s (commit %s)", subject, shortHash)

	return cmdAdd([]string{"--tags", tags, "--files", files, "--source", "git", "--commit", fullHash, text})
}

func cmdVersion(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("version non accetta argomenti")
	}
	fmt.Println("cairn", version)
	return nil
}

func trailerValue(message, name string) string {
	for _, line := range strings.Split(message, "\n") {
		colon := strings.IndexByte(line, ':')
		if colon < 0 || !strings.EqualFold(strings.TrimSpace(line[:colon]), name) {
			continue
		}
		return strings.TrimSpace(line[colon+1:])
	}
	return ""
}

func nonEmptyLines(raw string) []string {
	var lines []string
	for _, line := range strings.Split(raw, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func localHooksPath() (string, error) {
	cmd := exec.Command("git", "config", "--local", "--get", "core.hooksPath")
	out, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(out)), nil
	}
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return "", nil
	}
	return "", fmt.Errorf("impossibile leggere core.hooksPath: %w", err)
}

func gitOutput(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		detail := strings.TrimSpace(string(out))
		if detail == "" {
			return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
		}
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), detail)
	}
	return string(out), nil
}

func canonicalPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return filepath.Clean(path)
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err == nil {
		return resolved
	}
	return filepath.Clean(abs)
}

// misplacedFlag rileva un flag noto del FlagSet finito tra gli argomenti
// posizionali: succede quando l'utente lo scrive dopo il testo, e flag.Parse
// si ferma al primo token non-flag lasciando tutto il resto in fs.Args().
func misplacedFlag(fs *flag.FlagSet) (string, bool) {
	names := make(map[string]bool)
	fs.VisitAll(func(f *flag.Flag) {
		names[f.Name] = true
	})

	for _, tok := range fs.Args() {
		name := strings.TrimLeft(tok, "-")
		if name == tok {
			continue // nessun trattino iniziale: non può essere un flag
		}
		if eq := strings.IndexByte(name, '='); eq >= 0 {
			name = name[:eq]
		}
		if names[name] {
			return tok, true
		}
	}
	return "", false
}

func parseTags(raw string) []string {
	if raw == "" {
		return nil
	}
	seen := make(map[string]bool)
	var tags []string
	for _, part := range strings.Split(raw, ",") {
		t := strings.ToLower(strings.TrimSpace(part))
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		tags = append(tags, t)
	}
	return tags
}

func hasTag(tags []string, target string) bool {
	target = strings.ToLower(strings.TrimSpace(target))
	for _, t := range tags {
		if t == target {
			return true
		}
	}
	return false
}

// parseFiles, a differenza di parseTags, non normalizza il case: i path del
// filesystem sono case-sensitive sulla maggior parte dei sistemi.
func parseFiles(raw string) []string {
	if raw == "" {
		return nil
	}
	seen := make(map[string]bool)
	var files []string
	for _, part := range strings.Split(raw, ",") {
		f := strings.TrimSpace(part)
		if f == "" || seen[f] {
			continue
		}
		seen[f] = true
		files = append(files, f)
	}
	return files
}

func hasFile(files []string, target string) bool {
	target = strings.TrimSpace(target)
	for _, f := range files {
		if f == target {
			return true
		}
	}
	return false
}

func genID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString([]byte(time.Now().UTC().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(b)
}
