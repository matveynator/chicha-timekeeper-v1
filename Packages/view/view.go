package view

import (
	"embed"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"chicha/Models"
)

//go:embed templates/*
var templates embed.FS

type View struct{}

// loadTemplate loads templates embedded
func loadTemplate() *template.Template {
	var files = []string{"templates/index.tmpl"}

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
	v := new(View)
	r.SetHTMLTemplate(loadTemplate())

	// endpoints
	{
		r.GET("/", v.Homepage)
	}

	return v
}

func (v *View) Homepage(c *gin.Context) {
	laps := new([]Models.Lap)
	if err := Models.GetAllLaps(laps); err != nil {
		c.Error(err)
		log.Println(err)
		return
	}

	log.Println("all races: ", laps)

	c.HTML(http.StatusOK, "templates/index.tmpl", laps)
}
