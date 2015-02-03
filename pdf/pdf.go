package pdf

import (
	"strings"

	"code.google.com/p/gofpdf"
	"github.com/qanx/gopress/mdhtml"
)

func createPDFPageArray(file string) []string {
	return strings.SplitAfter(mdhtml.HTMLFromMarkdown(file), "<hr/>")
}

//MakePDF creates PDF document from single input file
func MakePDF(file string, output string) {
	htmlPDFarray := createPDFPageArray(file)
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetFont("Helvetica", "", 20)
	_, lineHt := pdf.GetFontSize()
	html := pdf.HTMLBasicNew()
	for _, page := range htmlPDFarray {
		pdf.AddPage()
		html.Write(lineHt, page)
	}
	pdf.OutputFileAndClose(output)
	//pdf.OutputAndClose(docWriter(pdf, 6))
}
