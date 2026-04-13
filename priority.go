package main

import (
	"sort"
	"strings"
)

type Priority string

const (
	PriorityHigh   Priority = "High"
	PriorityMedium Priority = "Medium"
	PriorityLow    Priority = "Low"
)

var highPrioritySubjects = []string{
	"deadline",
	"expert report",
	"please respond",
	"2nd request",
	"second request",
	"action required",
	"request for information",
	"notification of service",
	"notice of representation",
	"subpoena",
	"settlement",
	"settlement statement",
	"demand",
	"demand letter",
	"charge of discrimination",
	"twccrd",
	"eeoc",
	"twc crd",
	"mediation",
	"final decision",
	"time-sensitive",
	"rush question",
	"followup",
	"follow-up",
}

var highPrioritySenders = []string{
	"court",
	"district clerk",
	"county clerk",
	"eeo intake",
	"eeoc",
	"twccrd",
	"twc",
	"mediat",
	"opposing counsel",
}

var mediumPrioritySubjects = []string{
	"client",
	"case update",
	"contract questions",
	"filevine",
	"integration",
	"webhook",
	"automations",
	"lead docket",
	"scheduling a call",
	"interview",
	"application",
	"referral",
	"mailing address for settlement check",
}

var mediumPrioritySenders = []string{
	"carter law group",
	"tracy pace",
	"tania robinson",
	"heather davis",
	"amy carter",
	"javier perez",
	"david sanchez",
	"rachel ropp",
	"kyndal hetmer",
}

func classifyPriority(m message) Priority {
	subject := normalizePriorityText(m.Subject)
	fromName := normalizePriorityText(m.From.EmailAddress.Name)
	fromAddr := normalizePriorityText(m.From.EmailAddress.Address)
	combinedSender := fromName + " " + fromAddr

	for _, s := range highPrioritySubjects {
		if strings.Contains(subject, s) {
			return PriorityHigh
		}
	}

	for _, s := range highPrioritySenders {
		if strings.Contains(combinedSender, s) {
			return PriorityHigh
		}
	}

	for _, s := range mediumPrioritySubjects {
		if strings.Contains(subject, s) {
			return PriorityMedium
		}
	}

	for _, s := range mediumPrioritySenders {
		if strings.Contains(combinedSender, s) {
			return PriorityMedium
		}
	}

	return PriorityLow
}

func normalizePriorityText(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func bucketByPriority(messages []message) map[Priority][]message {
	result := map[Priority][]message{
		PriorityHigh:   {},
		PriorityMedium: {},
		PriorityLow:    {},
	}

	for _, m := range messages {
		p := classifyPriority(m)
		result[p] = append(result[p], m)
	}

	sortMessagesOldestFirst(result[PriorityHigh])
	sortMessagesOldestFirst(result[PriorityMedium])
	sortMessagesOldestFirst(result[PriorityLow])

	return result
}

func sortMessagesOldestFirst(messages []message) {
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].ReceivedDateTime.Before(messages[j].ReceivedDateTime)
	})
}
