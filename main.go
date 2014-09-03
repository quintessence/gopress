package main

//import (
//  "github.com/russross/blackfriday"
//)

import "os/exec"

func main(){
  sourceFile := "../temp/myfile.md"
  destinationDir := "../temp/myfile"
  copyCommand := exec.Command("cp", "-rf", sourceFile, destinationDir)
  err := copyCommand.Run()
  if err != nil {
    return
  }
}
