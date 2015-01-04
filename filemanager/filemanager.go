package filemanager

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

//ExtractInputFilename - does as described.
func ExtractInputFilename(inputFile string) string {
	inputPathArray := strings.Split(inputFile, "/")
	return strings.Split(inputPathArray[len(inputPathArray)-1], ".")[0]
}

//UpdateDestinationDirPath - does as described.
func UpdateDestinationDirPath(currentDestDir string, inputFile string, newDirFlag bool) string {
	newDestDir, _ := filepath.Split(inputFile)

	if currentDestDir == "" {
		return newDestDir + ExtractInputFilename(inputFile) + "/"
	}

	if newDirFlag {
		return currentDestDir + "/" + ExtractInputFilename(inputFile)
	}

	return currentDestDir
}

//FileOrDirectoryDoesNotExist - does as described.
func FileOrDirectoryDoesNotExist(inputPath string) bool {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return true
	}
	return false
}

//InputFileIsNotMarkdownFile - does as described.
func InputFileIsNotMarkdownFile(inputFile string) bool {
	return strings.ToLower(filepath.Ext(inputFile)) != ".md"
}

//ReplaceTildaWithHomeDir - does as described.
func ReplaceTildaWithHomeDir(filepath string) string {
	currentUser, _ := user.Current()
	return strings.Replace(filepath, "~", currentUser.HomeDir, 1)
}
