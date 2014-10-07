package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/alexcesaro/log/stdlog"
	"github.com/microcosm-cc/bluemonday"
	"github.com/qanx/gopress/mdhtml"
	"github.com/russross/blackfriday"
	"github.com/shurcooL/go/github_flavored_markdown"
)

func extractInputFilename(inputFile string) string {
	inputPathArray := strings.Split(inputFile, "/")
	return strings.Split(inputPathArray[len(inputPathArray)-1], ".")[0]
}

func updateDestinationDirPath(currentDestDir string, currentInputFile string, newDirFlag bool) string {
	newDestDir, _ := filepath.Split(currentInputFile)

	if currentDestDir == "" {
		return newDestDir + extractInputFilename(currentInputFile) + "/"
	}

	if newDirFlag {
		return currentDestDir + "/" + extractInputFilename(currentInputFile)
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
	var sourceFile string
	var destinationDir string
	var newDir bool

	flag.StringVar(&sourceFile, "inputFile", "NULL", "File to be converted to HTML presentation")
	flag.StringVar(&destinationDir, "outputDir", "", "Directory where HTML presentation will be written")
	flag.BoolVar(&newDir, "newDir", false, "Creates a new directory named after the file")

	logger := stdlog.GetFromFlags()

	flag.Parse()

	logger.Info("Starting gopress")

	sourceFile = replaceTildaWithHomeDir(sourceFile)
	destinationDir = replaceTildaWithHomeDir(destinationDir)

	if fileOrDirectoryDoesNotExist(sourceFile) {
		logger.Errorf("Input file or directory does not exist: %s", sourceFile)
		logger.Info("Exited with errors.")
		return
	}

	if inputFileIsNotMarkdownFile(sourceFile) {
		logger.Errorf("Input file is not a Markdown file: %s", sourceFile)
		logger.Info("Exited with errors.")
		return
	}

	sourceFileRead, errorReadFile := ioutil.ReadFile(sourceFile)
	if errorReadFile != nil {
		logger.Errorf("Could not read Markdown file: %s", sourceFile)
		return
	}

	markdownToHTML := blackfriday.MarkdownBasic(sourceFileRead)
	if markdownToHTML == nil {
		logger.Errorf("Markdown file empty. Please create content or use different file: %s", sourceFile)
		logger.Info("Exited with errors.")
		return
	}

	destinationDir = updateDestinationDirPath(destinationDir, sourceFile, newDir)

	if fileOrDirectoryDoesNotExist(destinationDir) {
		logger.Warningf("Output directory unspecified or does not exist, creating new directory: %s", destinationDir)
		mkdirCommand := exec.Command("mkdir", destinationDir)
		mkdirErr := mkdirCommand.Run()
		if mkdirErr != nil {
			logger.Errorf("Could not create new directory: %s", destinationDir)
			logger.Info("Exited with errors.")
			return
		}
		logger.Info("Successfully created new directory.")
	}

	copyCommand := exec.Command("cp", "-rf", sourceFile, destinationDir)
	cpErr := copyCommand.Run()

	if cpErr != nil {
		logger.Error("Could not copy files.")
		logger.Info("Exited with errors.")
	}
	logger.Infof("Successfully copied files to: %s", destinationDir)

	outputFile := destinationDir + "/" + extractInputFilename(sourceFile) + ".html"
	htmlFile, errorCreatingFile := os.Create(outputFile)
	if errorCreatingFile != nil {
		logger.Errorf("Could not create file: %s", outputFile)
		logger.Info("Exited with errors.")
		return
	}

	htmlHeader := `
<!DOCTYPE html>
<html>
	<head>

		<link href="css/reset.css" rel="stylesheet" />

		<meta charset="utf-8" />
		<meta name="viewport" content="width=1024" />
		<meta name="apple-mobile-web-app-capable" content="yes" />
		<link rel="shortcut icon" href="css/favicon.png" />
		<link rel="apple-touch-icon" href="css/apple-touch-icon.png" />
		<!-- Code Prettifier: -->
		<link href="css/highlight.css" type="text/css" rel="stylesheet" />
		<script type="text/javascript" src="js/highlight.pack.js"></script>
		<script>hljs.initHighlightingOnLoad();</script>
		<link href="css/style.css" rel="stylesheet" />
		<link href="http://fonts.googleapis.com/css?family=Lato:300,900" rel="stylesheet" />
	</head>
	<body>
		<div class="fallback-message">
			<p>Your browser <b>doesn't support the features required</b> by impress.js, so you are presented with a simplified version of this presentation.</p>
			<p>For the best experience please use the latest <b>Chrome</b>, <b>Safari</b> or <b>Firefox</b> browser.</p>
		</div>
	`
	htmlCSSstyle := `
		<style>
		.slide {
			color: #00786e;
		}
		h1 {
			color: orange;
		}
		</style>
		<div style="background-color: white; height: 100%;">
			<div>
				<img style="position: absolute; bottom: 0; width: 100%" src="http://i.imgur.com/QtxV5NQ.jpg" />
			</div>
		</div>
		<div id="impress">
			<div class='step slide' >
	`

	htmlFooter := `
			</div>
		</div>
		<script src="js/impress.js"></script>
		<script>impress().init();</script>
	</body>
</html>
	`
	_, _ = htmlFile.WriteString(htmlHeader)
	_, _ = htmlFile.WriteString(htmlCSSstyle)
	_, errorHTML := htmlFile.Write(bluemonday.UGCPolicy().SanitizeBytes(github_flavored_markdown.Markdown(markdownToHTML)))
	_, _ = htmlFile.WriteString(htmlFooter)
	if errorHTML != nil {
		logger.Errorf("Could not convert to HTML: %s", sourceFile)
		logger.Info("Exited with errors.")
		return
	}

	defer htmlFile.Close()
	htmlFile.Sync()
	logger.Info("Exited with no errors.")
	fmt.Println(mdhtml.SquareInteger(2))
}
