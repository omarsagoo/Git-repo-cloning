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
	fmt.Printf("[makeclones] \x1b[32;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// warning should be used to display a warning
func warning(format string, args ...interface{}) {
	fmt.Printf("[makeclones] \x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func completed(format string, args ...interface{}) {
	fmt.Printf("[COMPLETED] \x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
func pathComplete(format string, args ...interface{}) {
	fmt.Printf("[Saved to path] \x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

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
	// fmt.Println(banner)
	red := color.New(color.FgRed)

	boldRed := red.Add(color.Bold)
	boldRed.Println(banner)
	color.Blue(fmt.Sprintf("%[1]*s", (w+len(title))/2, title))
	fmt.Println()
}

func progressBar(numOfDirs int, results chan int) {
	writer := goterminal.New(os.Stdout)
	var z float64
	bar := "[" + strings.Repeat(" ", numOfDirs) + "]"
	for i := 0; i < numOfDirs; i++ {
		z += float64(<-results)
		// str := strings.Split(bar, "")
		// str[i+1] = "="
		// bar = strings.Join(str, "")
		bar = strings.Replace(bar, " ", string(4255), 1)
		writer.Clear()
		num := 100.0 / float64(numOfDirs) * z
		fmt.Fprintf(writer, "Downloading: %d/100 %s\n", int(num), bar)
		writer.Print()

	}
	writer.Reset()
}

var fileSize float64

func dirSize(dir string) float64 {

	allfiles, err := ioutil.ReadDir(dir)
	checkIfError(err)

	for _, file := range allfiles {
		if file.IsDir() {
			dirSize(dir + "/" + file.Name())
		}
		fileSize += float64(file.Size()) / 1000000.0
	}

	return fileSize
}
