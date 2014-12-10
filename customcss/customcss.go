package customcss

import (
	"io/ioutil"
)

// CSSToHTML takes in a custom CSS file and wraps it in HTML style tags
func CSSToHTML(customcssFilePath string) string {
	customcssFileByte, _ := ioutil.ReadFile(customcssFilePath)
	customcssFile := string(customcssFileByte)
	return "<style>" + customcssFile + "</style>"
}
