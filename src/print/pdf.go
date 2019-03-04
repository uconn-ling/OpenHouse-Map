package print

import (
  "bytes"
  "fmt"
  "io/ioutil"
  // "github.com/fogleman/gg"
  "github.com/jung-kurt/gofpdf"
)

var fileName string = "pdf/hello.pdf"
var imgName string = "../data/country_USA/USA-flag.png"

func CreatePdf () {
  pdf := gofpdf.New("P", "mm", "Letter", "")
  pdf.AddPage()
  pdf.SetFont("Arial", "B", 16)
  pdf.Cell(40, 10, "Hello, world")
  png, err := ioutil.ReadFile(imgName)
  // png, err := gg.LoadImage(imgName)
  if err == nil {
		rdr := bytes.NewReader(png)
		pdf = gofpdf.New("P", "mm", "A4", "")
		pdf.AddPage()
		_ = pdf.RegisterImageOptionsReader("sweden", gofpdf.ImageOptions{ImageType: "png", ReadDpi: true}, rdr)
		err = pdf.Error()
  }
  err = pdf.OutputFileAndClose(fileName)
  if err != nil {
    fmt.Println(err)
  }
}
