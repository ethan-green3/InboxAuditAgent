package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func listRootMailFolders(client *http.Client, token, mailboxUser string, includeHidden bool) ([]mailFolder, error) {
	escapedUser := url.PathEscape(mailboxUser)
	endpoint := fmt.Sprintf(
		"https://graph.microsoft.com/v1.0/users/%s/mailFolders?$top=100&includeHiddenFolders=%t&$select=id,displayName,childFolderCount,totalItemCount,parentFolderId,isHidden",
		escapedUser,
		includeHidden,
	)

	var all []mailFolder
	for endpoint != "" {
		var resp graphListResponse[mailFolder]
		if err := graphGET(client, token, endpoint, &resp); err != nil {
			return nil, err
		}
		all = append(all, resp.Value...)
		endpoint = resp.NextLink
	}

	return all, nil
}

func listChildFolders(client *http.Client, token, mailboxUser, parentFolderID string, includeHidden bool) ([]mailFolder, error) {
	escapedUser := url.PathEscape(mailboxUser)
	escapedParent := url.PathEscape(parentFolderID)

	endpoint := fmt.Sprintf(
		"https://graph.microsoft.com/v1.0/users/%s/mailFolders/%s/childFolders?$top=100&includeHiddenFolders=%t&$select=id,displayName,childFolderCount,totalItemCount,parentFolderId,isHidden",
		escapedUser,
		escapedParent,
		includeHidden,
	)

	var all []mailFolder
	for endpoint != "" {
		var resp graphListResponse[mailFolder]
		if err := graphGET(client, token, endpoint, &resp); err != nil {
			return nil, err
		}
		all = append(all, resp.Value...)
		endpoint = resp.NextLink
	}

	return all, nil
}

func walkAllFolders(client *http.Client, token, mailboxUser string, roots []mailFolder) ([]mailFolder, error) {
	var all []mailFolder
	queue := append([]mailFolder{}, roots...)
	seen := map[string]bool{}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if seen[current.ID] {
			continue
		}
		seen[current.ID] = true
		all = append(all, current)

		if current.ChildCount <= 0 {
			continue
		}

		children, err := listChildFolders(client, token, mailboxUser, current.ID, true)
		if err != nil {
			return nil, fmt.Errorf("listing children of %q failed: %w", current.DisplayName, err)
		}
		queue = append(queue, children...)
	}

	return all, nil
}

func listFolderMessages(client *http.Client, token, mailboxUser, folderID string, daysBack int, top int) ([]message, error) {
	escapedUser := url.PathEscape(mailboxUser)
	escapedFolderID := url.PathEscape(folderID)

	since := time.Now().UTC().AddDate(0, 0, -daysBack).Format("2006-01-02T15:04:05Z")

	baseURL := fmt.Sprintf(
		"https://graph.microsoft.com/v1.0/users/%s/mailFolders/%s/messages",
		escapedUser,
		escapedFolderID,
	)

	params := url.Values{}
	params.Set("$top", fmt.Sprintf("%d", top))
	params.Set("$orderby", "receivedDateTime desc")
	params.Set("$filter", fmt.Sprintf("receivedDateTime ge %s", since))
	params.Set("$select", "id,subject,receivedDateTime,from,conversationId,isRead,toRecipients,ccRecipients")

	endpoint := baseURL + "?" + params.Encode()

	var all []message
	for endpoint != "" {
		var resp graphListResponse[message]
		if err := graphGET(client, token, endpoint, &resp); err != nil {
			return nil, err
		}
		all = append(all, resp.Value...)
		endpoint = resp.NextLink
	}

	return all, nil
}

func findFolderByName(folders []mailFolder, name string) (*mailFolder, error) {
	for _, f := range folders {
		if strings.EqualFold(strings.TrimSpace(f.DisplayName), strings.TrimSpace(name)) {
			folderCopy := f
			return &folderCopy, nil
		}
	}
	return nil, fmt.Errorf("folder not found")
}
