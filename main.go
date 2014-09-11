package main

import (
	//  "github.com/russross/blackfriday"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alexcesaro/log/stdlog"
)

func updateDestinationDirPath(currentDestDir string, currentInputFile string, newDirFlag bool) string {
	if currentDestDir == "" || newDirFlag {
		inputPathArray := strings.Split(currentInputFile, "/")
		return currentDestDir + strings.Split(inputPathArray[len(inputPathArray)-1], ".")[0]
	}

	return currentDestDir
}

func fileOrDirectoryDoesNotExist(inputPath string) bool {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return true
	}
	return false
}

func verbosePrint(msgToPrint string, verbose bool) {
	if verbose {
		fmt.Printf(msgToPrint)
	}
}

func main() {
	var sourceFile string
	var destinationDir string
	var verboseMode bool
	var newDir bool

	flag.StringVar(&sourceFile, "inputFile", "", "File to be converted to HTML presentation")
	flag.StringVar(&destinationDir, "outputDir", "", "Directory where HTML presentation will be written")
	flag.BoolVar(&verboseMode, "v", false, "Verbose mode prints progress messages to screen")
	flag.BoolVar(&newDir, "newDir", false, "Creates a new directory named after the file")

	flag.Parse()

	logger := stdlog.GetFromFlags()
	logger.Info("Testing gopress writes to stdout.")

	if fileOrDirectoryDoesNotExist(sourceFile) {
		verbosePrint("Input file or directory does not exist: "+sourceFile+"\nExiting program.\n", verboseMode)
		return
	}

	destinationDir = updateDestinationDirPath(destinationDir, sourceFile, newDir)

	if fileOrDirectoryDoesNotExist(destinationDir) {
		verbosePrint("No such directory, creating new directory: "+destinationDir+"\n", verboseMode)
		mkdirCommand := exec.Command("mkdir", destinationDir)
		mkdirErr := mkdirCommand.Run()
		if mkdirErr != nil {
			verbosePrint("Could not create new directory.\nExiting program.\n", verboseMode)
			return
		}
		verbosePrint("Successfully created new directory.\n", verboseMode)
	}

	copyCommand := exec.Command("cp", "-rf", sourceFile, destinationDir)
	cpErr := copyCommand.Run()

	if cpErr != nil {
		verbosePrint("Could not copy files.\nExiting program.\n", verboseMode)
	}
	verbosePrint("Successfully copied files.\nExiting program.\n", verboseMode)
}
