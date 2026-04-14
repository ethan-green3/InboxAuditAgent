package main

import (
	"fmt"
	"strings"
)

func buildEmailBody(buckets map[Priority][]message) string {
	var sb strings.Builder

	fmt.Fprintln(&sb, "Friday Inbox Review - Loose Ends")
	fmt.Fprintln(&sb, "========================================")
	fmt.Fprintln(&sb)

	appendSection(&sb, "High Priority", buckets[PriorityHigh])
	appendSection(&sb, "Medium Priority", buckets[PriorityMedium])
	appendSection(&sb, "Low Priority", buckets[PriorityLow])

	return sb.String()
}

func appendSection(sb *strings.Builder, title string, messages []message) {
	fmt.Fprintf(sb, "%s (%d)\n", title, len(messages))
	fmt.Fprintln(sb, "----------------------------------------")

	for i, m := range messages {
		fmt.Fprintf(
			sb,
			"%d. %s | %s | %s\n",
			i+1,
			m.ReceivedDateTime.Format("2006-01-02 15:04"),
			nonEmpty(m.From.EmailAddress.Name, m.From.EmailAddress.Address),
			m.Subject,
		)
	}

	fmt.Fprintln(sb)
	fmt.Fprintln(sb)
}
