package main

import (
//  "github.com/russross/blackfriday"
  "os/exec"
  "flag"
  "fmt"
  "os"
)

func main(){
  var sourceFile string
  var destinationDir string

  flag.StringVar(&sourceFile, "inputFile", "", "File to be converted to HTML presentation")
  flag.StringVar(&destinationDir, "outputDir", "", "Directory where HTML presentation will be written")
 
  flag.Parse()

  if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
    fmt.Printf("No such file or directory: %s\nExiting program.\n", sourceFile)
    return
  }

  if _, err := os.Stat(destinationDir); os.IsNotExist(err) {
    fmt.Printf("No such directory, creating new directory: %s\n", destinationDir)
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
