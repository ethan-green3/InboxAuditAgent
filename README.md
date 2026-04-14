## Overview

Inbox Audit Agent is a lightweight automation that performs a weekly review through a GitHub Action of Outlook email activity to identify conversations that may still require a response.

At a high level, the system runs on a scheduled basis (once per week) and analyzes recent email activity by pulling data from both the Inbox and Sent Items using Microsoft Graph. It filters out low-value or irrelevant messages such as newsletters, automated notifications, internal system messages, and emails where the user is only copied.

After filtering, the remaining emails are grouped into conversation threads. For each thread, the system compares the most recent incoming message against the most recent outgoing reply. If no reply has been sent after the latest incoming message, the thread is flagged as potentially needing attention.

The resulting list of threads is then categorized into priority levels based on factors like subject content and sender type, helping surface the most important or time-sensitive items first. Finally, a clean summary report is generated and sent via email, providing a quick and actionable view of any “loose ends” from the past week.

The goal of the workflow is not to make decisions or generate responses, but to act as a reliable safety net—helping ensure that important emails are not overlooked or forgotten.
