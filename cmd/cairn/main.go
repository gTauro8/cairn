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
	"strings"
	"time"
)

const stateDir = ".cairn"
const logFile = stateDir + "/log.jsonl"

type entry struct {
	ID    string   `json:"id"`
	TS    string   `json:"ts"`
	Text  string   `json:"text"`
	Tags  []string `json:"tags,omitempty"`
	Files []string `json:"files,omitempty"`
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
	fmt.Fprintln(os.Stderr, "usage: cairn add [--tags a,b,c] [--files a,b] <text> | cairn log [--tag x] [--file p]")
}

func cmdAdd(args []string) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	tagsFlag := fs.String("tags", "", "tag separati da virgola")
	filesFlag := fs.String("files", "", "file del repo a cui la nota si riferisce, separati da virgola")
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
		ID:    genID(),
		TS:    time.Now().UTC().Format(time.RFC3339),
		Text:  text,
		Tags:  parseTags(*tagsFlag),
		Files: parseFiles(*filesFlag),
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
	if err := fs.Parse(args); err != nil {
		return err
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
		fmt.Printf("[%s] %s\n", e.TS, e.Text)
	}
	return scanner.Err()
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
