# gopress

[![baby-gopher](https://raw2.github.com/drnic/babygopher-site/gh-pages/images/babygopher-badge.png)](http://www.babygopher.org)

Go app based on [mdpress](https://github.com/egonSchiele/mdpress) - converts Markdown files to HTML slides with impressJs.

## Logging

By default, the logging level is set to "info". To change the logging level use the -log parameter, e.g.:

```bash
$ go run main.go -inputFile=/path/to/myfile.md -log=error
$ go run main.go -inputFile=/path/to/myfile.md -log=none
```

Allowed values are: debug, info, notice, warning, error, critical, alert, emergency, and none. Logging level must be lowercase.
