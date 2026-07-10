package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

const stateDir = ".cairn"
const logFile = stateDir + "/log.jsonl"

type entry struct {
	ID   string `json:"id"`
	TS   string `json:"ts"`
	Text string `json:"text"`
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
		err = cmdLog()
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
	fmt.Fprintln(os.Stderr, "usage: cairn add <text> | cairn log")
}

func cmdAdd(args []string) error {
	text := strings.TrimSpace(strings.Join(args, " "))
	if text == "" {
		return fmt.Errorf("add richiede un testo non vuoto")
	}

	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		return fmt.Errorf("impossibile creare %s: %w", stateDir, err)
	}

	e := entry{
		ID:   genID(),
		TS:   time.Now().UTC().Format(time.RFC3339),
		Text: text,
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

func cmdLog() error {
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
		fmt.Printf("[%s] %s\n", e.TS, e.Text)
	}
	return scanner.Err()
}

func genID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString([]byte(time.Now().UTC().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(b)
}
