package filemanager

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

//ExtractFilename - does as described.
func ExtractFilename(inputFile string) string {
	inputPathArray := strings.Split(inputFile, "/")
	return strings.Split(inputPathArray[len(inputPathArray)-1], ".")[0]
}

//UpdateDestPath - does as described.
func UpdateDestPath(currentDestDir string, inputFile string, newDirFlag bool) string {
	newDestDir, _ := filepath.Split(inputFile)

	if currentDestDir == "" {
		return newDestDir + ExtractFilename(inputFile) + "/"
	}

	if newDirFlag {
		return currentDestDir + "/" + ExtractFilename(inputFile)
	}

	return currentDestDir
}

//DoesNotExist - Checks if a file or directory does not exist.
func DoesNotExist(inputPath string) bool {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return true
	}
	return false
}

//IsNotMarkdown - does as described.
func IsNotMarkdown(inputFile string) bool {
	return strings.ToLower(filepath.Ext(inputFile)) != ".md"
}

func isMarkdown(inputFile string) bool {
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
			if isMarkdown(file.Name()) {
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
