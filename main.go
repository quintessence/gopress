package main

import (
//  "github.com/russross/blackfriday"
  "os/exec"
  "flag"
  "fmt"
  "os"
  "strings"
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

func main(){
  var sourceFile string
  var destinationDir string
  var verboseMode bool
  var newDir bool

  flag.StringVar(&sourceFile, "inputFile", "", "File to be converted to HTML presentation")
  flag.StringVar(&destinationDir, "outputDir", "", "Directory where HTML presentation will be written")
  flag.BoolVar(&verboseMode, "v", false, "Verbose mode prints progress messages to screen")
  flag.BoolVar(&newDir, "newDir", false, "Creates a new directory named after the file")

  flag.Parse()

  if fileOrDirectoryDoesNotExist(sourceFile){ 
    if verboseMode {
      fmt.Printf("No such file or directory: %s\nExiting program.\n", sourceFile)
    }
    return
  }

  destinationDir = updateDestinationDirPath(destinationDir, sourceFile, newDir)

  if fileOrDirectoryDoesNotExist(destinationDir){ 
    if verboseMode {
      fmt.Printf("No such directory, creating new directory: %s\n", destinationDir)
    }
    mkdirCommand := exec.Command("mkdir", destinationDir)
    mkdirErr := mkdirCommand.Run()
    if mkdirErr != nil {
      output, _ := mkdirCommand.Output()
      println(output)
      fmt.Printf("Could not create new directory.\nExiting program.\n")
      return
    }
    fmt.Printf("Successfully created new directory.\n")
  }

  copyCommand := exec.Command("cp", "-rf", sourceFile, destinationDir)
  cpErr := copyCommand.Run()
  
  if cpErr != nil {
   output, _ := copyCommand.Output()
   println(output)
   fmt.Printf("Could not copy files.\nExiting program.\n")
  }

  fmt.Printf("Successfully copied files.\nExiting.\n")

}
