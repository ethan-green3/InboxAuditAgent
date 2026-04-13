package main

import (
	"sort"
	"strings"
)

func isLikelyNoise(m message) bool {
	subject := strings.ToLower(strings.TrimSpace(m.Subject))
	fromName := strings.ToLower(strings.TrimSpace(m.From.EmailAddress.Name))
	fromAddr := strings.ToLower(strings.TrimSpace(m.From.EmailAddress.Address))

	if strings.Contains(fromName, "in teams") || strings.Contains(subject, "sent a message") {
		return true
	}

	if strings.Contains(subject, "[action required] new application") {
		return true
	}

	if strings.Contains(subject, "women leaders in law") {
		return true
	}

	for _, prefix := range noiseSubjectPrefixes {
		if strings.HasPrefix(subject, prefix) {
			return true
		}
	}

	for _, s := range noiseSubjects {
		if strings.Contains(subject, s) {
			return true
		}
	}

	for _, s := range noiseSenders {
		if strings.Contains(fromName, s) || strings.Contains(fromAddr, s) {
			return true
		}
	}

	for _, s := range noiseAddrParts {
		if strings.Contains(fromAddr, s) {
			return true
		}
	}

	return false
}

func filterInboxCandidates(messages []message, myEmail string) []message {
	result := make([]message, 0, len(messages))

	for _, m := range messages {
		if m.ConversationID == "" {
			continue
		}

		if isLikelyNoise(m) {
			continue
		}
		if isOnlyCC(m, myEmail) {
			continue
		}

		// Optional: only consider messages that were actually read,
		// since the COO's pain point is "I read it and forgot it."
		if !m.IsRead {
			continue
		}

		result = append(result, m)
	}

	return result
}

func latestMessagePerConversation(messages []message) map[string]message {
	latest := make(map[string]message)

	for _, m := range messages {
		if m.ConversationID == "" {
			continue
		}

		existing, ok := latest[m.ConversationID]
		if !ok || m.ReceivedDateTime.After(existing.ReceivedDateTime) {
			latest[m.ConversationID] = m
		}
	}

	return latest
}

func needsAttention(inboxByConv map[string]message, sentByConv map[string]message) []message {
	result := make([]message, 0)

	for convID, inboxMsg := range inboxByConv {
		sentMsg, ok := sentByConv[convID]

		// No sent reply in this conversation at all
		if !ok {
			result = append(result, inboxMsg)
			continue
		}

		// Latest sent happened before latest inbox message
		if sentMsg.ReceivedDateTime.Before(inboxMsg.ReceivedDateTime) {
			result = append(result, inboxMsg)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ReceivedDateTime.Before(result[j].ReceivedDateTime)
	})

	return result
}

func isOnlyCC(m message, myEmail string) bool {
	myEmail = strings.ToLower(strings.TrimSpace(myEmail))

	inTo := false
	inCC := false

	for _, r := range m.ToRecipients {
		if strings.ToLower(r.EmailAddress.Address) == myEmail {
			inTo = true
			break
		}
	}

	for _, r := range m.CcRecipients {
		if strings.ToLower(r.EmailAddress.Address) == myEmail {
			inCC = true
			break
		}
	}

	// If not in TO, but in CC → exclude
	if !inTo && inCC {
		return true
	}

	return false
}
