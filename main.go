package main

import (
//  "github.com/russross/blackfriday"
  "os/exec"
  "flag"
)

func main(){
  var sourceFile string
  var destinationDir string

  flag.StringVar(&sourceFile, "inputFile", "", "File to be converted to HTML presentation")
  flag.StringVar(&destinationDir, "outputDir", "", "Directory where HTML presentation will be written")
 
  copyCommand := exec.Command("cp", "-rf", sourceFile, destinationDir)
  err := copyCommand.Run()
  
  if err != nil {
   output, _ := copyCommand.Output()
   println(output)
   println(sourceFile, destinationDir)
  }
}
