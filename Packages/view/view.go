package view

import (
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"log"
	"net/http"
)

//go:embed templates/*
var templates embed.FS

type View struct {
	//box *rice.Box
}

// loadTemplate loads templates embedded
func loadTemplate() *template.Template {
	var files = []string{"index.tmpl"}

	t := template.New("")

	for _, filePath := range files {
		file, err := templates.Open(filePath)
		if err != nil {
			log.Panicln("file load error: ", err)
		}

		h, err := io.ReadAll(file)
		if err != nil {
			log.Panicln("file read error:", err)
		}

		t, err = t.New(filePath).Parse(string(h))
		if err != nil {
			log.Panicln("template parce error:", t, err)
		}
	}

	return t
}

func New(r *gin.Engine) *View {
	templates := loadTemplate()

	r.SetHTMLTemplate(templates)

	return &View{}
}

func (v *View) Homepage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", nil)
}
