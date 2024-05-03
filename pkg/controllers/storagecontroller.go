package controllers

import (
	"fmt"

	"example.com/myproject/pkg/storagemanager"
	"github.com/labstack/echo"
)

type FilesPayload struct {
	Files []string `json:"files"`
}

func handleGetFiles(c echo.Context) error {
	connection := storagemanager.Conn.MongoConnection
	if connection == nil {
		conn, err := storagemanager.CreateConnections()
		if err != nil {
			return c.JSON(400, "Can't connect to database")
		}

		connection = conn.MongoConnection
	}

	files, err := connection.GetAllFiles()
	if err != nil {
		panic(err.Error())
	}

	return c.JSON(200, files)
}

func handleUploadFiles(c echo.Context) error {
	files := make([]storagemanager.File, 0)

	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(400, "Wrong input type")
	}

	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {

			newFile, err := storagemanager.UploadFile(c.Request().Context(), fileHeader)
			if err != nil {
				continue
			}

			files = append(files, newFile)
		}
	}

	return c.JSON(200, files)
}

func handleDownloadFiles(c echo.Context) error {
	azureConnection := storagemanager.Conn.AzureConnection
	if azureConnection == nil {
		conn, err := storagemanager.CreateConnections()
		if err != nil {
			return c.JSON(400, "Can't connect to database")
		}

		azureConnection = conn.AzureConnection
	}

	payload := new(FilesPayload)
	if err := c.Bind(payload); err != nil {
		return c.JSON(400, "Wrong input type")
	}

	fileID := payload.Files[0]

	mongoConnection, err := storagemanager.ConnectClient(storagemanager.GetCredentials().MongoURI)
	if err != nil {
		return c.JSON(400, "Can't connect to database")
	}

	fileMeta, err := mongoConnection.GetFile(fileID)
	if err != nil {
		return c.JSON(400, "File doesn't exist")
	}

	fileData, err := azureConnection.DownloadFile(fileID, "testcontainer")
	if err != nil {
		return c.JSON(400, "Error downloading file")
	}

	// Set the response headers to indicate that the response is a file
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileID))

	if fileMeta.ContentType == "application/pdf" {
		c.Response().Header().Set("Content-Type", "application/pdf")
	} else {
		c.Response().Header().Set("Content-Type", "application/octet-stream")
	}

	// Write the file data to the response body
	_, err = c.Response().Write(fileData)
	if err != nil {
		return c.JSON(500, "Error writing file data to response")
	}

	return nil
}

func handleDeleteFiles(c echo.Context) error {
	payload := new(FilesPayload)
	if err := c.Bind(payload); err != nil {
		return c.JSON(400, "Wrong input type")
	}

	for _, file := range payload.Files {
		fmt.Println(file)
		err := storagemanager.DeleteFile(file)
		if err != nil {
			return c.JSON(400, "Error deleting file")
		}
	}

	return c.JSON(200, "File deleted")
}
