package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	tenantID := mustEnv("TENANT_ID")
	clientID := mustEnv("CLIENT_ID")
	clientSecret := mustEnv("CLIENT_SECRET")
	mailboxUser := mustEnv("MAILBOX_USER")

	token, err := getGraphToken(tenantID, clientID, clientSecret)
	if err != nil {
		log.Fatalf("failed to get Graph token: %v", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}

	fmt.Println("Authenticated to Microsoft Graph.")
	fmt.Println("Mailbox:", mailboxUser)
	fmt.Println()

	rootFolders, err := listRootMailFolders(client, token, mailboxUser, true)
	if err != nil {
		log.Fatalf("failed to list root folders: %v", err)
	}

	allFolders, err := walkAllFolders(client, token, mailboxUser, rootFolders)
	if err != nil {
		log.Fatalf("failed walking folder tree: %v", err)
	}

	inboxFolder, err := findFolderByName(allFolders, "Inbox")
	if err != nil {
		log.Fatalf("could not find Inbox: %v", err)
	}

	sentFolder, err := findFolderByName(allFolders, "Sent Items")
	if err != nil {
		log.Fatalf("could not find Sent Items: %v", err)
	}

	inboxMessages, err := listFolderMessages(client, token, mailboxUser, inboxFolder.ID, 7, 100)
	if err != nil {
		log.Fatalf("failed to pull Inbox messages: %v", err)
	}

	sentMessages, err := listFolderMessages(client, token, mailboxUser, sentFolder.ID, 7, 100)
	if err != nil {
		log.Fatalf("failed to pull Sent Items messages: %v", err)
	}

	fmt.Printf("Pulled %d Inbox messages from last 7 days\n", len(inboxMessages))
	fmt.Printf("Pulled %d Sent Items messages from last 7 days\n", len(sentMessages))

	filteredInbox := filterInboxCandidates(inboxMessages, mailboxUser)

	inboxByConversation := latestMessagePerConversation(filteredInbox)
	sentByConversation := latestMessagePerConversation(sentMessages)

	attentionThreads := needsAttention(inboxByConversation, sentByConversation)
	priorityBuckets := bucketByPriority(attentionThreads)

	fmt.Printf("\nThreads needing attention this week: %d\n", len(attentionThreads))

	printPrioritySection("High Priority", priorityBuckets[PriorityHigh])
	printPrioritySection("Medium Priority", priorityBuckets[PriorityMedium])
	printPrioritySection("Low Priority", priorityBuckets[PriorityLow])

}
