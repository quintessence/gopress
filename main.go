package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/alexcesaro/log/stdlog"
	"github.com/russross/blackfriday"
	"github.com/shurcooL/go/github_flavored_markdown"
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
	logger.Info("Starting gopress")

	if fileOrDirectoryDoesNotExist(sourceFile) {
		verbosePrint("Input file or directory does not exist: "+sourceFile+"\nExiting program.\n", verboseMode)
		logger.Errorf("Input file or directory does not exist: %s", sourceFile)
		return
	}

	destinationDir = updateDestinationDirPath(destinationDir, sourceFile, newDir)

	if fileOrDirectoryDoesNotExist(destinationDir) {
		verbosePrint("No such directory, creating new directory: "+destinationDir+"\n", verboseMode)
		logger.Warningf("Output directory unspecified or does not exist, creating new directory: %s", destinationDir)
		mkdirCommand := exec.Command("mkdir", destinationDir)
		mkdirErr := mkdirCommand.Run()
		if mkdirErr != nil {
			verbosePrint("Could not create new directory.\nExiting program.\n", verboseMode)
			logger.Errorf("Could not create new directory.")
			logger.Info("Exiting gopress with errors")
			return
		}
		verbosePrint("Successfully created new directory.\n", verboseMode)
		logger.Info("Successfully created new directory.")
	}

	copyCommand := exec.Command("cp", "-rf", sourceFile, destinationDir)
	cpErr := copyCommand.Run()

	if cpErr != nil {
		verbosePrint("Could not copy files.\nExiting program.\n", verboseMode)
		logger.Errorf("Could not copy files.")
		logger.Info("Exiting gopress with errors")
	}
	verbosePrint("Successfully copied files.\nExiting program.\n", verboseMode)
	logger.Info("Successfully copied files.")

	sourceFileRead, errorReadFile := ioutil.ReadFile(sourceFile)
	if errorReadFile != nil {
		logger.Errorf("Could not read Markdown file.")
		return
	}

	//Now for some Markdown/HTML
	//Testing: https://github.com/shurcooL/go/blob/master/github_flavored_markdown/main_test.go
	var writer io.Writer = os.Stdout
	htmlBytes := blackfriday.MarkdownBasic(sourceFileRead)
	writer.Write(github_flavored_markdown.Markdown(htmlBytes))
	_, errorHTML := os.Stdout.Write(github_flavored_markdown.Markdown(htmlBytes))
	if errorHTML != nil {
		logger.Errorf("Could not convert to HTML: " + sourceFile)
	}
	logger.Info("Exiting gopress")
}
