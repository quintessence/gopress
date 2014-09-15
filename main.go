package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/alexcesaro/log/stdlog"
	"github.com/russross/blackfriday"
	"github.com/shurcooL/go/github_flavored_markdown"
)

func extractInputFilename(inputFile string) string {
	inputPathArray := strings.Split(inputFile, "/")
	return strings.Split(inputPathArray[len(inputPathArray)-1], ".")[0]
}

func updateDestinationDirPath(currentDestDir string, currentInputFile string, newDirFlag bool) string {
	if currentDestDir == "" || newDirFlag {
		return currentDestDir + extractInputFilename(currentInputFile)
	}

	return currentDestDir
}

func fileOrDirectoryDoesNotExist(inputPath string) bool {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return true
	}
	return false
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

	logger := stdlog.GetFromFlags()

	flag.Parse()

	logger.Info("Starting gopress")

	if fileOrDirectoryDoesNotExist(sourceFile) {
		logger.Errorf("Input file or directory does not exist: %s", sourceFile)
		return
	}

	destinationDir = updateDestinationDirPath(destinationDir, sourceFile, newDir)

	if fileOrDirectoryDoesNotExist(destinationDir) {
		logger.Warningf("Output directory unspecified or does not exist, creating new directory: %s", destinationDir)
		mkdirCommand := exec.Command("mkdir", destinationDir)
		mkdirErr := mkdirCommand.Run()
		if mkdirErr != nil {
			logger.Errorf("Could not create new directory.")
			logger.Info("Exiting gopress with errors")
			return
		}
		logger.Info("Successfully created new directory.")
	}

	copyCommand := exec.Command("cp", "-rf", sourceFile, destinationDir)
	cpErr := copyCommand.Run()

	if cpErr != nil {
		logger.Errorf("Could not copy files.")
		logger.Info("Exiting gopress with errors")
	}
	logger.Info("Successfully copied files.")

	sourceFileRead, errorReadFile := ioutil.ReadFile(sourceFile)
	if errorReadFile != nil {
		logger.Errorf("Could not read Markdown file.")
		return
	}

	outputFile := destinationDir + "/" + extractInputFilename(sourceFile) + ".html"
	htmlFile, errorCreatingFile := os.Create(outputFile)
	if errorCreatingFile != nil {
		logger.Errorf("Could not create file: ")
	}

	markdownToHTML := blackfriday.MarkdownBasic(sourceFileRead)
	_, _ = htmlFile.WriteString("<html><body>\n")
	_, errorHTML := htmlFile.Write(github_flavored_markdown.Markdown(markdownToHTML))
	_, _ = htmlFile.WriteString("</body></html>")
	if errorHTML != nil {
		logger.Errorf("Could not convert to HTML: " + sourceFile)
	}
	defer htmlFile.Close()
	htmlFile.Sync()
	logger.Info("Exiting gopress")
}
