package view

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"

	"chicha/Models"
)

type View struct {
	static embed.FS
}

func createMyRender() multitemplate.Renderer {
	f := template.FuncMap{
		"timestampRender": timestampRender,
	}

	r := multitemplate.NewRenderer()

	r.AddFromFilesFuncs("index", f, "static/templates/index.tmpl")
	r.AddFromFilesFuncs("race", f, "static/templates/race.tmpl", "static/templates/race_table.tmpl")
	r.AddFromFilesFuncs("race_table", f, "static/templates/race_table.tmpl")
	return r
}

// loadTemplate loads templates embedded
func (v *View) loadTemplate() *template.Template {
	var files = []string{
		"static/templates/index.tmpl",
		"static/templates/race.tmpl",
		"static/templates/race_table.tmpl",
	}

	t := template.New("")

	for _, filePath := range files {
		file, err := v.static.Open(filePath)
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

// return fs for serve static files
func (v *View) getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(v.static, "static/assets")
	if err != nil {
		log.Fatal(err)
	}
	return http.FS(fsys)
}

func New(r *gin.Engine, static embed.FS) *View {
	v := &View{static: static}
	r.HTMLRender = createMyRender()
	//r.SetHTMLTemplate(v.loadTemplate())

	// endpoints
	{
		// static files
		r.StaticFS("/static/assets/", v.getFileSystem())

		r.GET("/", v.Homepage)
		r.GET("/race/:id", v.RaceView)
	}

	return v
}

func (v *View) Homepage(c *gin.Context) {
	laps := new([]Models.Lap)

	//sqlsr
	s := `
		SELECT 
		       race_id, 
		       MIN(discovery_unix_time) as discovery_unix_time, 
		       MIN(discovery_time) as discovery_time
		FROM laps 
		GROUP BY race_id
	`

	if err := Models.DB.Raw(s).Find(laps).Error; err != nil {
		c.Error(err)
		log.Println(err)
		return
	}

	c.HTML(http.StatusOK, "index", laps)
}

func (v *View) RaceView(c *gin.Context) {
	raceID := c.Params.ByName("id")
	laps := new([]Models.Lap)

	if err := Models.GetAllResultsByRaceId(laps, raceID); err != nil {
		c.Error(err)
		log.Println(err)
		return
	}

	for _, v := range *laps {
		timestampRender(v.LapTime)
	}

	reslt := gin.H{
		"RaceID": raceID,
		"Laps":   laps,
	}

	if c.Query("updtable") == "true" {
		c.HTML(http.StatusOK, "race_table", reslt)
		return
	}
	c.HTML(http.StatusOK, "race", reslt)
}

func timestampRender(ts int64) string {
	return time.UnixMilli(ts).UTC().Format("15:04:05.000")
}
