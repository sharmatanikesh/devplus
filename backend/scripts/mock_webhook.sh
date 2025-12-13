#!/bin/bash

# Configuration
API_URL="http://localhost:8081/api/v1/webhook/github"
REPO_ID="$1" # Pass the database ID of the repository (or we need github_repo_id which is in payload)
# Actually, the payload uses github_repo_id. 
# We need to ensure the Repo exists in DB first with that ID.

if [ -z "$1" ]; then
  echo "Usage: ./mock_webhook.sh <GITHUB_REPO_ID>"
  echo "Example: ./mock_webhook.sh 123456"
  exit 1
fi

GITHUB_REPO_ID=$1
PR_ID=$((RANDOM % 1000 + 1000))
PR_NUMBER=$((RANDOM % 100 + 1))

echo "Sending Mock Pull Request 'Opened' Event..."
echo "GitHub Repo ID: $GITHUB_REPO_ID"
echo "PR Number: $PR_NUMBER"

curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
  "action": "opened",
  "number": '$PR_NUMBER',
  "pull_request": {
    "id": '$PR_ID',
    "number": '$PR_NUMBER',
    "title": "Feature: Add AI Analysis",
    "state": "open",
    "user": {
      "login": "octocat",
      "id": 1
    },
    "head": {
        "ref": "feature/ai-analysis",
        "sha": "deadbeef"
    }
  },
  "repository": {
    "id": '$GITHUB_REPO_ID',
    "name": "devplus",
    "full_name": "octocat/devplus",
    "owner": {
      "login": "octocat"
    }
  }
}'

echo -e "\n\nDone! Check your backend logs."
