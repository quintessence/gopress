package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/alexcesaro/log/stdlog"
	"github.com/qanx/gopress/mdhtml"
)

func extractInputFilename(inputFile string) string {
	inputPathArray := strings.Split(inputFile, "/")
	return strings.Split(inputPathArray[len(inputPathArray)-1], ".")[0]
}

func updateDestinationDirPath(currentDestDir string, currentInputFile string, newDirFlag bool) string {
	newDestDir, _ := filepath.Split(currentInputFile)

	if currentDestDir == "" {
		return newDestDir + extractInputFilename(currentInputFile) + "/"
	}

	if newDirFlag {
		return currentDestDir + "/" + extractInputFilename(currentInputFile)
	}

	return currentDestDir
}

func fileOrDirectoryDoesNotExist(inputPath string) bool {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return true
	}
	return false
}

func inputFileIsNotMarkdownFile(inputFile string) bool {
	return strings.ToLower(filepath.Ext(inputFile)) != ".md"
}

func replaceTildaWithHomeDir(filepath string) string {
	currentUser, _ := user.Current()
	return strings.Replace(filepath, "~", currentUser.HomeDir, 1)
}

func main() {
	var sourceFile string
	var destinationDir string
	var newDir bool

	flag.StringVar(&sourceFile, "inputFile", "NULL", "File to be converted to HTML presentation")
	flag.StringVar(&destinationDir, "outputDir", "", "Directory where HTML presentation will be written")
	flag.BoolVar(&newDir, "newDir", false, "Creates a new directory named after the file")

	logger := stdlog.GetFromFlags()

	flag.Parse()

	logger.Info("Starting gopress")

	sourceFile = replaceTildaWithHomeDir(sourceFile)
	destinationDir = replaceTildaWithHomeDir(destinationDir)

	if fileOrDirectoryDoesNotExist(sourceFile) {
		logger.Errorf("Input file or directory does not exist: %s", sourceFile)
		logger.Warning("Exited with errors.")
		return
	}

	if inputFileIsNotMarkdownFile(sourceFile) {
		logger.Errorf("Input file is not a Markdown file: %s", sourceFile)
		logger.Warning("Exited with errors.")
		return
	}

	sourceFileRead, errorReadFile := ioutil.ReadFile(sourceFile)
	if errorReadFile != nil {
		logger.Errorf("Could not read Markdown file: %s", sourceFile)
		logger.Warning("Exited with errors.")
		return
	}

	if len(sourceFileRead) == 0 {
		logger.Errorf("Markdown file empty. Please create content or use different file: %s", sourceFile)
		logger.Warning("Exited with errors.")
		return
	}

	destinationDir = updateDestinationDirPath(destinationDir, sourceFile, newDir)

	if fileOrDirectoryDoesNotExist(destinationDir) {
		logger.Warningf("Output directory unspecified or does not exist, creating new directory: %s", destinationDir)
		mkdirCommand := exec.Command("mkdir", destinationDir)
		mkdirErr := mkdirCommand.Run()
		if mkdirErr != nil {
			logger.Errorf("Could not create new directory: %s", destinationDir)
			logger.Warning("Exited with errors.")
			return
		}
		logger.Info("Successfully created new directory.")
	}

	copyCommand := exec.Command("cp", "-rf", "css", "impress_css", "js", destinationDir)
	cpErr := copyCommand.Run()

	if cpErr != nil {
		logger.Error(cpErr)
		logger.Error("Could not copy files.")
		logger.Warning("Exited with errors.")
	}
	logger.Infof("Successfully copied files to: %s", destinationDir)

	outputFile := destinationDir + "/" + extractInputFilename(sourceFile) + ".html"
	htmlFile, errorCreatingFile := os.Create(outputFile)
	if errorCreatingFile != nil {
		logger.Errorf("Could not create file: %s", outputFile)
		logger.Warning("Exited with errors.")
		return
	}

	_, errorHTML := htmlFile.WriteString(mdhtml.GenerateHTML(sourceFile))
	if errorHTML != nil {
		logger.Errorf("Could not write to HTML file: %s", outputFile)
		logger.Warning("Exited with errors.")
		return
	}

	defer htmlFile.Close()
	htmlFile.Sync()
	logger.Info("Exited with no errors.")
}
