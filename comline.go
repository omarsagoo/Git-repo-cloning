package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/apoorvam/goterminal"
	"github.com/fatih/color"
)

// checkIfError should be used to naively panic if an error is not nil.
func checkIfError(err error) {
	if err == nil {
		return
	}
	fmt.Printf("[makeclones] \x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("%s", err))
	//os.Exit(1)
}

// info should be used to describe the example commands that are about to run.
func info(format string, args ...interface{}) {
	fmt.Printf("\n[makeclones] \x1b[32;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// warning should be used to display a warning
func warning(format string, args ...interface{}) {
	fmt.Printf("[makeclones] \x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// completed to display a comlpleted status when all repos are fully cloned
func completed(format string, args ...interface{}) {
	fmt.Printf("\n[COMPLETED] \x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// displays the path all of the repos where saved to.
func pathComplete(format string, args ...interface{}) {
	fmt.Printf("\n[Saved to path] \x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// shows the MakeClones banner in the terminal when the program is initially run
func showBanner() {
	author := "Droxey"
	version := "v 1.0.0"

	title := fmt.Sprintf("MakeClones by %s (%s)", author, version)

	banner := `
	    __    __     ______     __  __     ______     ______     __         ______     __   __     ______     ______    
	   /\ '-./  \   /\  __ \   /\ \/ /    /\  ___\   /\  ___\   /\ \       /\  __ \   /\ '-.\ \   /\  ___\   /\  ___\   
	   \ \ \-./\ \  \ \  __ \  \ \  _'-.  \ \  __\   \ \ \____  \ \ \____  \ \ \/\ \  \ \ \-.  \  \ \  __\   \ \___  \  
	    \ \_\ \ \_\  \ \_\ \_\  \ \_\ \_\  \ \_____\  \ \_____\  \ \_____\  \ \_____\  \ \_\\'\_\  \ \_____\  \/\_____\ 
	     \/_/  \/_/   \/_/\/_/   \/_/\/_/   \/_____/   \/_____/   \/_____/   \/_____/   \/_/ \/_/   \/_____/   \/_____/ 
	                                                                                                                    
	`
	allLines := strings.Split(banner, "\n")

	w := len(allLines[2])
	red := color.New(color.FgRed)

	boldRed := red.Add(color.Bold)
	boldRed.Println(banner)
	color.Blue(fmt.Sprintf("%[1]*s", (w+len(title))/2, title))
	fmt.Println()
}

// displays a progress bar that shows the status of the repos being cloned
func progressBar(numOfDirs int, results chan int) {
	// using the goterminal module, instantiate a new progress bar writer.
	writer := goterminal.New(os.Stdout)
	// instatiate a var that stores the results from the results chan
	var z float64
	// creates a string representation of the progress bar to be used in the terminal
	bar := "[" + strings.Repeat("-", numOfDirs) + "]"

	// create a for loop that loops exactly the number of dirs time.
	// in order to unload the results channel and increment the progress bar
	for i := 0; i < numOfDirs; i++ {
		z += float64(<-results)
		bar = strings.Replace(bar, "-", string(4255), 1)
		// clear the terminal after the channel has been unloaded, keeps the progress bar in the terminal
		writer.Clear()
		// incrememnt the progress bar completion percent by 100 divded by the num of dirs.
		num := 100.0 / float64(numOfDirs) * z
		fmt.Fprintf(writer, "Downloading: %d/100 %s\n", int(num), bar)
		writer.Print()

	}
	writer.Reset()
}

var fileSize float64

// displays the size of all the files that were created.
func dirSize(dir string) float64 {

	allfiles, err := ioutil.ReadDir(dir)
	checkIfError(err)

	// recurssively checks through all of the files that were created and incrememnts the filesize
	for _, file := range allfiles {
		if file.IsDir() {
			dirSize(dir + "/" + file.Name())
		}
		fileSize += float64(file.Size()) / 1000000.0
	}

	return fileSize
}
