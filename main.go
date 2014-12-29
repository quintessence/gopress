package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alexcesaro/log/stdlog"
	"github.com/qanx/gopress/mdhtml"
)

func extractInputFilename(inputFile string) string {
	inputPathArray := strings.Split(inputFile, "/")
	return strings.Split(inputPathArray[len(inputPathArray)-1], ".")[0]
}

func updateDestinationDirPath(currentDestDir string, inputFile string, newDirFlag bool) string {
	newDestDir, _ := filepath.Split(inputFile)

	if currentDestDir == "" {
		return newDestDir + extractInputFilename(inputFile) + "/"
	}

	if newDirFlag {
		return currentDestDir + "/" + extractInputFilename(inputFile)
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
	var sourceFilePath string
	var destinationDir string
	var newDir bool

	flag.StringVar(&sourceFilePath, "inputFile", "NULL", "File to be converted to HTML presentation")
	flag.StringVar(&destinationDir, "outputDir", "", "Directory where HTML presentation will be written")
	flag.BoolVar(&newDir, "newDir", false, "Creates a new directory named after the file")

	logger := stdlog.GetFromFlags()

	flag.Parse()

	logger.Info("Starting gopress")

	sourceFilePath = replaceTildaWithHomeDir(sourceFilePath)
	destinationDir = replaceTildaWithHomeDir(destinationDir)

	if fileOrDirectoryDoesNotExist(sourceFilePath) {
		logger.Errorf("Input file or directory does not exist: %s", sourceFilePath)
		logger.Warning("Exited with errors.")
		return
	}

	if inputFileIsNotMarkdownFile(sourceFilePath) {
		logger.Errorf("Input file is not a Markdown file: %s", sourceFilePath)
		logger.Warning("Exited with errors.")
		return
	}

	sourceFileRead, errorReadFile := ioutil.ReadFile(sourceFilePath)
	if errorReadFile != nil {
		logger.Errorf("Could not read Markdown file: %s", sourceFilePath)
		logger.Warning("Exited with errors.")
		return
	}

	if len(sourceFileRead) == 0 {
		logger.Errorf("Markdown file empty. Please create content or use different file: %s", sourceFilePath)
		logger.Warning("Exited with errors.")
		return
	}

	destinationDir = updateDestinationDirPath(destinationDir, sourceFilePath, newDir)

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
		logger.Error("Could not copy CSS and JS files.")
		logger.Warning("Exited with errors.")
	}
	logger.Infof("Successfully copied files to: %s", destinationDir)

	outputFile := destinationDir + "/" + extractInputFilename(sourceFilePath) + ".html"
	htmlFile, errorCreatingFile := os.Create(outputFile)
	if errorCreatingFile != nil {
		logger.Errorf("Could not create file: %s", outputFile)
		logger.Warning("Exited with errors.")
		return
	}

	//Convert Markdown to HTML
	htmlContents := mdhtml.GenerateHTML(sourceFilePath)

	//Locate image paths specified in HTML
	findImagePaths := regexp.MustCompile("image(s)?.*(.(?i)(jp(e)?g|png|gif|bmp|tiff))")
	imagesToCopy := findImagePaths.FindAllString(htmlContents, -1)

	//Copy images to 'images' directory if there are existing images to copy. Create 'images' directory if needed.
	imagesDirectory := destinationDir + "/" + "images"
	if len(imagesToCopy) > 0 {
		//Create 'images' directory if not present.
		if fileOrDirectoryDoesNotExist(imagesDirectory) {
			mkdirCommand := exec.Command("mkdir", imagesDirectory)
			mkdirErr := mkdirCommand.Run()
			if mkdirErr != nil {
				logger.Errorf("Could not create 'images' directory: %s", imagesDirectory)
				logger.Warning("Exited with errors.")
				return
			}
			logger.Infof("Successfully created 'images' directory: %s", imagesDirectory)
		}

		//Copy images to 'images' directory
		for i := 0; i < len(imagesToCopy); i++ {
			copyImagesCommand := exec.Command("cp", filepath.Dir(sourceFilePath)+"/"+imagesToCopy[i], imagesDirectory)
			copyImagesError := copyImagesCommand.Run()
			if copyImagesError != nil {
				logger.Errorf("Could not copy image file: %s", filepath.Dir(sourceFilePath)+"/"+imagesToCopy[i])
				logger.Warning("Exited with errors.")
				return
			}
		}
		logger.Infof("Successfully copied image files to: %s", imagesDirectory)
	}

	//Write HTML to file
	_, errorHTML := htmlFile.WriteString(htmlContents)
	if errorHTML != nil {
		logger.Errorf("Could not write to HTML file: %s", outputFile)
		logger.Warning("Exited with errors.")
		return
	}

	//Close file and exit program.
	defer htmlFile.Close()
	htmlFile.Sync()
	logger.Info("Exited with no errors.")
}
