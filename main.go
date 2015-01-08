package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/alexcesaro/log/stdlog"
	"github.com/qanx/gopress/filemanager"
	"github.com/qanx/gopress/mdhtml"
)

func main() {

	var sourceFilePath string
	var cssDir string
	var destinationDir string
	var newDir bool
	var allTheFiles bool
	var files []string

	flag.StringVar(&sourceFilePath, "inputFile", "NULL", "Comma separated list of file(s) to be converted to HTML presentation")
	flag.StringVar(&cssDir, "cssDir", "NULL", "Directory where CSS/JS files are located")
	flag.StringVar(&destinationDir, "outputDir", "", "Directory where HTML presentation will be written")
	flag.BoolVar(&newDir, "newDir", false, "Creates a new directory named after the file")
	flag.BoolVar(&allTheFiles, "all", false, "Used in conjunction with inputFiles. Will grab all Markdown ")

	logger := stdlog.GetFromFlags()

	flag.Parse()
	logger.Info("Starting gopress")

	if cssDir == "NULL" {
		cssDir = sourceFilePath
	}

	destinationDir = filemanager.ReplaceTildaWithHomeDir(destinationDir)
	/*
		if filemanager.DoesNotExist(sourceFilePath) {
			logger.Errorf("Input file or directory does not exist: %s", sourceFilePath)
			logger.Warning("Exited with errors.")
			return
		}
	*/

	files = filemanager.MakeFileList(sourceFilePath, allTheFiles)
	/*
		fmt.Println("Printing files:")
		for _, file := range files {
			fmt.Println(file)
		}
	*/

	/*
		if filemanager.IsNotMarkdown(sourceFilePath) {
			logger.Errorf("Input file is not a Markdown file: %s", sourceFilePath)
			logger.Warning("Exited with errors.")
			return
		}
	*/

	for _, file := range files {
		sourceFileRead, errorReadFile := ioutil.ReadFile(file)
		if errorReadFile != nil {
			logger.Errorf("Could not read Markdown file: %s", file)
			logger.Warning("Exited with errors.")
			return
		}

		if len(sourceFileRead) == 0 {
			logger.Errorf("Markdown file empty. Please create content or use different file: %s", file)
			logger.Warning("Exited with errors.")
			return
		}

		destinationDir = filemanager.UpdateDestPath(destinationDir, file, newDir)

		if filemanager.DoesNotExist(destinationDir) {
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

		if filemanager.DoesNotExist(cssDir+"css") || filemanager.DoesNotExist(cssDir+"impress_css") || filemanager.DoesNotExist(cssDir+"js") {
			logger.Errorf("CSS/JS files do not exist in specified directory: %s", cssDir)
			logger.Warning("If the CSS/JS files are not in the same directory as the Markdown files, please specify the directory with the cssDir flag.")
			logger.Warning("Exited with errors.")
			return
		}

		if filemanager.DoesNotExist(destinationDir+"/css") || filemanager.DoesNotExist(destinationDir+"/impress_css") || filemanager.DoesNotExist(destinationDir+"/js") {
			copyCommand := exec.Command("cp", "-rf", cssDir+"/css", cssDir+"/impress_css", cssDir+"/js", destinationDir)
			fmt.Println("cssDir is: " + cssDir + "/css")
			cpErr := copyCommand.Run()

			if cpErr != nil {
				logger.Error("Could not copy CSS and JS files.")
				logger.Warning("Exited with errors.")
				return
			}
			logger.Infof("Successfully copied CSS and JS files to: %s", destinationDir)
		}

		outputFile := destinationDir + "/" + filemanager.ExtractFilename(file) + ".html"
		htmlFile, errorCreatingFile := os.Create(outputFile)
		if errorCreatingFile != nil {
			logger.Errorf("Could not create file: %s", outputFile)
			logger.Warning("Exited with errors.")
			return
		}

		//Convert Markdown to HTML
		htmlContents := mdhtml.GenerateHTML(file)

		//Locate image paths specified in HTML
		findImagePaths := regexp.MustCompile("image(s)?.*(.(?i)(jp(e)?g|png|gif|bmp|tiff))")
		imagesToCopy := findImagePaths.FindAllString(htmlContents, -1)

		//Copy images to 'images' directory if there are existing images to copy. Create 'images' directory if needed.
		imagesDirectory := destinationDir + "/" + "images"
		if len(imagesToCopy) > 0 {
			//Create 'images' directory if not present.
			if filemanager.DoesNotExist(imagesDirectory) {
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
				copyImagesCommand := exec.Command("cp", filepath.Dir(file)+"/"+imagesToCopy[i], imagesDirectory)
				copyImagesError := copyImagesCommand.Run()
				if copyImagesError != nil {
					logger.Errorf("Could not copy image file: %s", filepath.Dir(file)+"/"+imagesToCopy[i])
					logger.Warning("Exited with errors.")
					return
				}
			}
			logger.Infof("Successfully copied image files for %s to: %s", filemanager.ExtractFilename(file), imagesDirectory)
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
	}

	logger.Info("Exited with no errors.")
}
