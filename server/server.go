package main

import (
	"encoding/json"
	"fmt"
	"golang-yaml/validation"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	m "github.com/hashicorp/go-multierror"
)

type PushEventPayload struct {
	Commits []struct {
		Url      string   `json:"url"`
		Added    []string `json:"added"`
		Removed  []string `json:"removed"`
		Modified []string `json:"modified"`
	} `json:"commits"`
}

// TODO : we need check with the commits used in PR and must be specific to PR events.
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	event := r.Header.Get("X-GitHub-Event")

	if event == "push" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(" failed to read request : ", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		var pushEvent PushEventPayload
		err = json.Unmarshal(body, &pushEvent)
		if err != nil {
			log.Println("failed to unmarshal : ", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		var errorList error
		for _, commit := range pushEvent.Commits {
			modifiedFiles := append(append(commit.Added, commit.Removed...), commit.Modified...)
			for _, file := range modifiedFiles {
				myurl := commit.Url
				myurl = strings.Replace(myurl, "github.com", "raw.githubusercontent.com", 1)
				myurl = strings.Replace(myurl, "/commit/", "/", 1)
				myurl += "/" + file
				fmt.Println(myurl)

				response, err := http.Get(myurl)
				if err != nil {
					fmt.Println(err)
					return
				}

				defer response.Body.Close()

				content, err1 := ioutil.ReadAll(response.Body)
				if err1 != nil {
					fmt.Println(err1)
					continue
				}

				individualErrors := validation.ValidateYamlFile(file, content)
				errorList = m.Append(errorList, individualErrors)
			}
		}

		listOfErrors := errorList.(*m.Error)
		if len(listOfErrors.Errors) > 0 {
			fmt.Println(listOfErrors)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, listOfErrors)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "modified files have no error")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "event handled successfully")
}

// TODO : we need to think of getting all the yaml files from the repo
func main() {
	port := "8080"
	log.Printf("server is on port %s : ", port)
	http.HandleFunc("/payload", webhookHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
