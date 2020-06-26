package main

import (
	"fmt"
	"os"
	"strings"
	"sync/atomic"

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
var numOfReposCloned int64

// initialize variable to store any repositories that might not have been cloned
var dirsNotFound []string

// initialize variable to store the numOfDirs created.
var numOfDirs int

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

	// creates 10 workers to start concurrently cloning repos
	for i := 0; i < 10; i++ {
		go cloneWorker(repoDirChan, results)
	}

	for _, row := range sheet.Rows {
		for _, cell := range row {
			if cell.Row > uint(skip) {
				cellPos := cell.Pos()
				if string(cellPos[0]) == column && len(cell.Value) > 0 {
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
	defer progressBar(int(numOfDirs), results)
}

// worker to pass the information from the struct into the clone helper function
func cloneWorker(pairs chan repoDirPair, results chan int) {
	for pair := range pairs {
		results <- clone(pair.Token, pair.RepoURL, pair.Dir)
	}
}

// cloning herlper function to work with the channels to clone the repos
func clone(token, repoURL, directory string) int {
	_, err := g.PlainClone(directory, false, &g.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "makeclones", // This can be anything except an empty string
			Password: token,
		},
		URL: repoURL,
		// do not want any progress info sent to the terminal.
		Progress: nil,
	})
	checkIfError(err)

	// if there is an error, add the repo that was not cloned and remove the directory that was created
	if err != nil {
		dirsNotFound = append(dirsNotFound, repoURL)
		os.Remove(directory)
		return 0
	}

	// if no error, increment number of repos cloned, using atomic to account for race conditions
	atomic.AddInt64(&numOfReposCloned, 1)
	return 1

}
