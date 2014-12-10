package customcss

import (
	"io/ioutil"
)

// CSSToHTML takes in a custom CSS file and wraps it in HTML style tags
func CSSToHTML(customcssFilePath string) string {
	customcssFileReadByte, _ := ioutil.ReadFile(customcssFilePath)
	customcssFileRead := string(customcssFileReadByte)
	return "<style>" + customcssFileRead + "</style>"
}
