package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	// GitHub Actions passes action inputs as env vars named
	// INPUT_<NAME_UPPERCASED_WITH_UNDERSCORES>. The input "who-to-greet"
	// in action.yml becomes INPUT_WHO-TO-GREET normalized to
	// INPUT_WHO_TO_GREET.
	who := os.Getenv("INPUT_WHO_TO_GREET")
	if who == "" {
		who = "World"
	}

	fmt.Printf("Hello, %s!\n", who)

	now := time.Now().Format(time.RFC3339)

	// Outputs are written as key=value lines to the file at $GITHUB_OUTPUT.
	// GitHub mounts this file into the container automatically for
	// docker-based actions, so it just works the same as in a normal job step.
	outputFile := os.Getenv("GITHUB_OUTPUT")
	if outputFile == "" {
		fmt.Println("GITHUB_OUTPUT not set, skipping output write (not running in Actions)")
		return
	}

	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open GITHUB_OUTPUT: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, "time=%s\n", now); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write output: %v\n", err)
		os.Exit(1)
	}
}
