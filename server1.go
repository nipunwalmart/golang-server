package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type PullRequestEvent struct {
	Action  string `json:"action"`
	PullReq struct {
		Number int `json:"number"`
		Base   struct {
			Sha string `json:"sha"`
		} `json:"base"`
		Head struct {
			Sha string `json:"sha"`
		} `json:"head"`
	} `json:"pull_request"`
	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"repository"`
}

func main() {
	http.HandleFunc("/payload", handleWebhook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Parse the event payload
	eventType := r.Header.Get("X-GitHub-Event")
	switch eventType {
	case "pull_request":
		var event PullRequestEvent
		err := json.Unmarshal(body, &event)
		if err != nil {
			http.Error(w, "Error parsing pull request event", http.StatusBadRequest)
			return
		}

		// Handle pull request event
		handlePullRequestEvent(event)
		// default:
		// 	fmt.Printf("Received event type: %s\n", eventType)
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
}

func handlePullRequestEvent(event PullRequestEvent) {
	switch event.Action {
	case "opened", "synchronize":
		fmt.Printf("Pull request #%d updated\n", event.PullReq.Number)

		// Get the Git diff between the base and head commits
		diffCmd := exec.Command("git", "diff", "--name-only", event.PullReq.Base.Sha, event.PullReq.Head.Sha)
		diffOutput, err := diffCmd.Output()
		if err != nil {
			fmt.Println("Error executing 'git diff' command:", err)
			return
		}

		changedFiles := strings.Split(string(diffOutput), "\n")
		for _, file := range changedFiles {
			if file != "" {
				fmt.Println("Changed file:", file)
			}
		}
	}
}
