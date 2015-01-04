package filemanager

import (
	"io/ioutil"
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

func inputFileIsMarkdownFile(inputFile string) bool {
	return strings.ToLower(filepath.Ext(inputFile)) == ".md"
}

//ReplaceTildaWithHomeDir - does as described.
func ReplaceTildaWithHomeDir(filepath string) string {
	currentUser, _ := user.Current()
	return strings.Replace(filepath, "~", currentUser.HomeDir, 1)
}

//MakeFileList - makes list of one or more files
func MakeFileList(path string, useAllMDfiles bool) []string {
	var files []string
	if useAllMDfiles {
		path = ReplaceTildaWithHomeDir(path)
		directoryContents, _ := ioutil.ReadDir(path)
		for _, file := range directoryContents {
			if inputFileIsMarkdownFile(file.Name()) {
				files = append(files, path+file.Name())
			}
		}
	} else {
		files = strings.Split(path, ",")
		for i := 0; i < len(files); i++ {
			files[i] = ReplaceTildaWithHomeDir(files[i])
		}
	}
	return files
}
