package controllers

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"

	"example.com/myproject/pkg/pdfbuilder"
	"github.com/labstack/echo"
)

func handleImagesToPDF(c echo.Context) error {
	// Parse the form data
	err := c.Request().ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		fmt.Println("Error parsing form data:", err)
		return err
	}

	// Get the map of uploaded files
	fileHeaders := c.Request().MultipartForm.File

	// Check if any files were uploaded
	if len(fileHeaders) == 0 {
		fmt.Println("No files uploaded")
		return errors.New("no files uploaded")
	}

	// Initialize a slice to store file headers
	var files []*multipart.FileHeader

	// Iterate over the map of file headers
	for _, fileHeader := range fileHeaders {
		// Iterate over the slice of file headers for the current key

		files = append(files, fileHeader...)
	}

	// Call the BuildPDF function with the slice of file headers
	err, pdfBytes := pdfbuilder.PDFFromImages(files)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Content-Type", "application/pdf")

	// Return the generated PDF file
	return c.Blob(http.StatusOK, "application/pdf", pdfBytes)
}
