package main

import (
	"fmt"
	"os"
)

func customCSSFileExists(inputPath string) bool {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func main() {
	customCSSFile := "custom.css"
	fmt.Println("hello world!")
	if customCSSFileExists(customCSSFile) {
		fmt.Println("I found the custom.css file!")
	} else {
		fmt.Println("No custom.css file.")
	}

}
