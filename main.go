package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/user"
)

const logFilePath = "/log/logit.txt"

func getCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		return err.Error()
	}
	return cwd
}

func writeLogEntry(resource string) error {
	currentUser, err := user.Current()
	username := "unknown"
	if err != nil {
		log.Printf("Warning: could not get current user: %v", err)
	} else {
		username = currentUser.Username
	}
	log.Printf("%s accessed %s", username, resource)
	return nil
}

func setupLogFile() error {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open or create log file: %v", err)
	}

	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime)
	cwd := getCwd()

	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slogger.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"logit",
		slog.String("path", cwd),
		slog.Int("status", 200),
	)

	return nil
}

func main() {

	if err := setupLogFile(); err != nil {
		log.Fatalf("Error setting up log file: %v", err)
	}

	resource := "unknown"
	if len(os.Args) > 1 {
		resource = os.Args[1]
	} else {
		fmt.Printf("no arguments provided")
	}

	if err := writeLogEntry(resource); err != nil {
		log.Fatalf("Error writing log entry: %v", err)
	}
}
