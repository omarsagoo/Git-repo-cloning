package main

import(
	"fmt"
	"flag"
	"os"
	"time"
)

func start() {

	start := time.Now()

	fmt.Println("start main")
	var sheet string
	var column string
	var authToken string
	var skipRows int
	var tabIndex int

	flag.StringVar(&sheet, "sheet", "", "Google Sheets spreadsheet ID (Required)")
	flag.StringVar(&column, "column", "", "Column to scrape. Make sure data is in the format username/reponame (Required)")
	flag.StringVar(&authToken, "token", "", "GitHub Personal Access Token (Create one at https://github.com/settings/tokens/new) with full control of private repositories (Required)")
	flag.IntVar(&skipRows, "skip", 0, "Skip a number of rows to accommodate headers")
	flag.IntVar(&tabIndex, "tab", 0, "Spreadsheet tab to look for the specified column")

	flag.Parse()

	showBanner()

	if sheet == "" ||  column == "" || authToken == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	MakeClones(sheet, tabIndex, column, authToken, skipRows)

	// display the number of repos and length of time in the terminal.
	since := time.Since(start).Seconds()
	minutes := int(since / 60.0)
	seconds := int(since) % 60

	size := dirSize("github.com/")
	sizeStr := fmt.Sprintf("%.2f MB", size)

	if size > 1000 {
		size = size / 1000
		sizeStr = fmt.Sprintf("%.2f GB", size)
	}

	completed("Cloned %d out of %d repos in %d minutes and %d seconds", numOfReposCloned, numOfDirs, minutes, seconds)
	if len(dirsNotFound) != 0 {
		repos := "repo"
		if len(dirsNotFound) > 1 {
			repos = "repos"
		}
		info("%d %s not cloned and directories removed:", len(dirsNotFound), repos)
		for _, repo := range dirsNotFound {
			fmt.Println(repo)
		}
	}
	pathComplete("github.com/https:/github.com/ (%s)", sizeStr)
}