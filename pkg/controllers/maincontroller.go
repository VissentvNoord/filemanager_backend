package controllers

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func Initialize(e *echo.Echo) {
	e.Renderer = NewTemplates()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", nil)
	})

	e.POST("/pdf", handleImagesToPDF)

	e.POST("/files", handleGetFiles)
	e.POST("/upload", handleUploadFiles)
	e.POST("/download", handleDownloadFiles)
	e.DELETE("/files", handleDeleteFiles)
}
