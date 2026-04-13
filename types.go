package main

import "time"

type graphListResponse[T any] struct {
	Value    []T    `json:"value"`
	NextLink string `json:"@odata.nextLink"`
}

type mailFolder struct {
	ID             string `json:"id"`
	DisplayName    string `json:"displayName"`
	ChildCount     int    `json:"childFolderCount"`
	TotalItems     int    `json:"totalItemCount"`
	ParentFolderID string `json:"parentFolderId"`
	IsHidden       bool   `json:"isHidden"`
	ODataType      string `json:"@odata.type"`
}

type emailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type recipient struct {
	EmailAddress emailAddress `json:"emailAddress"`
}

type message struct {
	ID               string      `json:"id"`
	Subject          string      `json:"subject"`
	ReceivedDateTime time.Time   `json:"receivedDateTime"`
	ConversationID   string      `json:"conversationId"`
	From             recipient   `json:"from"`
	IsRead           bool        `json:"isRead"`
	ToRecipients     []recipient `json:"toRecipients"`
	CcRecipients     []recipient `json:"ccRecipients"`
}
