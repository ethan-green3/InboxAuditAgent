package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func mustEnv(key string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return v
}

func envOrDefault(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}

func nonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func printPrioritySection(title string, messages []message) {
	fmt.Printf("\n%s (%d)\n", title, len(messages))
	fmt.Println("----------------------------------------")

	for i, m := range messages {
		fmt.Printf(
			"%d. %s | %s | %s\n",
			i+1,
			m.ReceivedDateTime.Format("2006-01-02 15:04"),
			nonEmpty(m.From.EmailAddress.Name, m.From.EmailAddress.Address),
			m.Subject,
		)
	}
}
