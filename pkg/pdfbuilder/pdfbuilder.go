package pdfbuilder

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

func pdfTemplate() (template gofpdf.Pdf, width float64, height float64) {
	// Create a new PDF instance
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Calculate PDF width and image width
	pdfWidth, pdfHeight := pdf.GetPageSize()
	return pdf, pdfWidth, pdfHeight
}

func PDFFromImages(fileHeaders []*multipart.FileHeader) (error, []byte) {
	// Create a new PDF instance
	pdf, pdfWidth, _ := pdfTemplate()

	// Loop through the file headers and add images to the PDF
	for _, fileHeader := range fileHeaders {
		// Add page for every image header
		pdf.AddPage()

		// Open the uploaded file
		src, err := fileHeader.Open()
		if err != nil {
			panic(err)
		}
		defer src.Close()

		extension := filepath.Ext(fileHeader.Filename)
		// Create a temporary file to store the uploaded file
		tempFile, err := os.CreateTemp("", "uploaded_image_*"+extension)
		if err != nil {
			panic(err)
		}
		defer os.Remove(tempFile.Name()) // Clean up the temporary file
		defer tempFile.Close()

		// Copy the contents of the uploaded file to the temporary file
		_, err = io.Copy(tempFile, src)
		if err != nil {
			panic(err)
		}

		imageWidth, _ := getImageDimensions(tempFile.Name())
		scale := pdfWidth / imageWidth

		// Add the image to the PDF
		pdf.ImageOptions(tempFile.Name(), 0, 0, imageWidth*scale, 0, false, gofpdf.ImageOptions{}, 0, "")
	}

	var buf bytes.Buffer

	if err := pdf.Output(&buf); err != nil {
		return err, nil
	}

	// Output the PDF to a file
	err := pdf.OutputFileAndClose("output/output.pdf")
	if err != nil {
		panic(err)
	}

	fmt.Println("PDF successfully created")

	return nil, buf.Bytes()
}

func getImageDimensions(imagePath string) (float64, float64) {
	file, err := os.Open(imagePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		panic(err)
	}

	return float64(img.Width), float64(img.Height)
}
