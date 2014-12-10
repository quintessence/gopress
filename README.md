# gopress

[![baby-gopher](images/babygopher-badge.png)](http://www.babygopher.org)

Go app based on [mdpress](https://github.com/egonSchiele/mdpress) - converts Markdown files to HTML slides with impressJs.

## Usage

To run gopress, you only need to specify the input file. If there is no output directory specified, a new directory named after the input file will be created in your current working directory where the output files will be stored.

```bash
$ go run main.go -inputFile=/path/to/myfile.md
```

If you would like to store the output to a specific subdirectory, then use the outputDir flag:

```bash
$ go run main.go -inputFile=/path/to/myfile.md -outputDir=/path/to/output
```

When an output directory is specified, by default the application will not create a new subdirectory in that directory - gopress will simply use the specified directory. If you would like it to create a subdirectory in the parent directory, then use the newDir flag:

```bash
$ go run main.go -inputFile=/path/to/myfile.md -outputDir=/path/to/output -newDir
```

This will create a new subdirectory named "myfile" in /path/to/output.

### Custom CSS

By default, gopress will look for a `custom.css` file in the same directory as your input file.
If no `custom.css` file is present, then the default attributes will be used as defined in `css/style.css`.

### Images

Locally stored images must be in an `images` subdirectory of where your input file is stored. This is because
the images used in the presentation will be copied to the `images` subdirectory of the presentation.

### Logging

By default, the logging level is set to "info". To change the logging level use the -log parameter, e.g.:

```bash
$ go run main.go -inputFile=/path/to/myfile.md -log=error
$ go run main.go -inputFile=/path/to/myfile.md -log=none
```

Allowed values are: debug, info, notice, warning, error, critical, alert, emergency, and none. Logging level must be lowercase.
