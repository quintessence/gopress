# gopress

[![baby-gopher](images/babygopher-badge.png)](http://www.babygopher.org)

Go app inspired by [mdpress](https://github.com/egonSchiele/mdpress) - converts Markdown files to HTML/impressJs presentations.

## Install

Quick and easy installation:

```
go install github.com/quintessence/gopress
```

You will also need the `cssJS.tar` file (do **not** rename this file). There are a few different ways to do this:

**Via the command line**

```
wget https://github.com/quintessence/gopress/raw/master/cssJS.tar
```

Mac users: If you do not already have `wget`, please run:
```
brew install wget
```

**In GitHub's GUI**

You can right-click on the file link and select "Save Link As..."

**Option for Chrome users: GitHub Mate**

Install link and documentation are located at the [project repo](https://github.com/camsong/chrome-github-mate). Once installed, you can download the file by clicking on it.

## Usage

To run gopress, you only need to specify the input file. If there is no output directory specified, a new directory named after the input file will be created in your current working directory where the output files will be stored.

```bash
gopress -inputFile=/path/to/myfile.md
```

**How to use multiple files**

There are two ways to convert multiple Markdown files to HTML presentations. The first is using a comma separated
list:

```bash
gopress -inputFile=/path/to/myfile.md,/path/to/mysecondfile.md
```

The second is using the `-all` flag, which grabs all Markdown files in the specified directory:

```bash
gopress -inputFile=/path/to/files/ -all
```

**Specify an Output Directory**

If you would like to store the output to a specific subdirectory, then use the outputDir flag:

```bash
gopress -inputFile=/path/to/myfile.md -outputDir=/path/to/output
```

When an output directory is specified, by default the application will not create a new subdirectory in that directory - gopress will simply use the specified directory. If you would like it to create a subdirectory in the parent directory, then use the newDir flag:

```bash
gopress -inputFile=/path/to/myfile.md -outputDir=/path/to/output -newDir
```

This will create a new subdirectory named "myfile" in /path/to/output. You can also use the `newDir` flag without `outputDir` and gopress will create the new "myfile" directory in the input directory.

### CSS & JS

The `cssJS.tar` file contains all the CSS and JS needed to run the presentation. By default, gopress expects this file to be in the same directory as the input file(s). If it is not, you can use the `cssDir` flag to specify its location:

```
gopress -inputFile=/path/to/myfile.md -cssDir=/path/to/cssJS.tar
```

You should not change the name or untar these files (the application handles this).

**Custom CSS**

gopress will look for a `custom.css` file in the same directory as your input file.
If no `custom.css` file is present, then the default attributes will be used as defined in `css/style.css`.
Please keep in mind that the CSS from this file is being added to the header of the HTML presentation file.

### Images

Locally stored images must be in an `images` subdirectory of where your input file is stored. This is because
the images used in the presentation will be copied to the `images` subdirectory of the presentation's HTML file. Please note that only the following raster extensions are supported - there are no vector extensions currently supported:

* BMP
* BPG
* GIF
* JPG/JPEG/JFIF
* PBM, PGM, PPM, and PNM
* PNG
* TIFF
* WEBP

### Logging

By default, the logging level is set to "info". To change the logging level use the `log` flag:

```bash
gopress -inputFile=/path/to/myfile.md -log=error
gopress -inputFile=/path/to/myfile.md -log=none
```

Allowed values are: debug, info, notice, warning, error, critical, alert, emergency, and none. Logging level must be lowercase.
