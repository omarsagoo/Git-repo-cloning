package main

import (
	"fmt"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	s "gopkg.in/Iwark/spreadsheet.v2"
	g "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type repoDirPair struct {
	RepoURL string
	Dir     string
	Token   string
}

// initialize variable to store the number of repos cloned off of the initial sheet.
// for benchmarking purposes
var numOfReposCloned int

// MakeClones from a spreadsheet column
// first optimized version: 21 repos 61 seconds
func MakeClones(sheetID string, tabIndex int, column string, token string, skip int) {
	service, err := s.NewService()
	checkIfError(err)

	sheets, err := service.FetchSpreadsheet(sheetID)
	checkIfError(err)

	sheet, err := sheets.SheetByIndex(uint(tabIndex))
	checkIfError(err)

	repoDirChan := make(chan repoDirPair, 500)
	results := make(chan int, 500)

	name := sheet.Properties.Title
	url := "https://docs.google.com/spreadsheets/d/" + sheetID
	fmt.Printf("Google Sheet URL: %s\nSheet Name: %s\nColumn: %s (%d rows skipped)\n\n", url, name, column, skip)

	var numOfDirs int

	// creates 10 workers to start concurrently cloning repos
	for i := 0; i < 10; i++ {
		go cloneWorker(repoDirChan, results)
	}

	for _, row := range sheet.Rows {
		for _, cell := range row {
			if cell.Row > uint(skip) {
				cellPos := cell.Pos()
				if string(cellPos[0]) == column && len(cell.Value) > 0 {
					checkIfError(err)

					prefix := "github.com/"
					directory := "github.com/" + cell.Value
					repoURL := cell.Value

					if !strings.HasPrefix(repoURL, "https://"+prefix) {
						repoURL = "https://" + prefix + cell.Value
					}

					// warning("creating directory %s...", directory)
					err = os.MkdirAll(directory, os.ModePerm)
					checkIfError(err)

					numOfDirs++
					// add a struct of the directory RepoUrl and github token into the channel
					repoDirChan <- repoDirPair{Dir: directory, RepoURL: repoURL, Token: token}
				}
			}
		}
	}

	// close the channel that was storing the repo struct
	close(repoDirChan)
	progressBar(numOfDirs, results)
}

// worker to pass the information from the struct into the clone helper function
func cloneWorker(pairs chan repoDirPair, results chan int) {
	for pair := range pairs {
		results <- clone(pair.Token, pair.RepoURL, pair.Dir)
	}
}

// cloning herlper function to work with the channels to clone the repos
func clone(token, repoURL, directory string) int {
	// info("cloning %s into %s...", repoURL, directory)
	_, err := g.PlainClone(directory, false, &g.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "makeclones", // This can be anything except an empty string
			Password: token,
		},
		URL: repoURL,
		// Progress: os.Stdout,
		Progress: nil,
	})
	checkIfError(err)

	if err == nil {
		numOfReposCloned++
		return 1
	}
	return 0
}
