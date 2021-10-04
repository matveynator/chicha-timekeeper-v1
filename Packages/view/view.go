package view

import (
	"embed"
	"html/template"
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

func (v *View) setupRenderer() multitemplate.Renderer {
	f := template.FuncMap{
		"timestampRender": timestampRender,
	}

	r := multitemplate.NewRenderer()

	index, _ := v.static.ReadFile("static/templates/index.tmpl")
	race, _ := v.static.ReadFile("static/templates/race.tmpl")
	raceTable, _ := v.static.ReadFile("static/templates/race_table.tmpl")
	raceTableView, _ := v.static.ReadFile("static/templates/race_table_view.tmpl")

	r.AddFromStringsFuncs("index", f, string(index))
	r.AddFromStringsFuncs("race", f, string(race), string(raceTable))
	r.AddFromStringsFuncs("race_table_view", f, string(raceTableView), string(raceTable))

	return r
}

//// loadTemplate loads templates embedded
//func (v *View) loadTemplate() *template.Template {
//	t := template.New("")
//
//	for _, filePath := range files {
//		file, err := v.static.Open(filePath)
//		if err != nil {
//			log.Panicln("file load error: ", err)
//		}
//
//		h, err := io.ReadAll(file)
//		if err != nil {
//			log.Panicln("file read error:", err)
//		}
//
//		t, err = t.New(filePath).Parse(string(h))
//		if err != nil {
//			log.Panicln("template parce error:", t, err)
//		}
//	}
//
//	return t
//}

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
	r.HTMLRender = v.setupRenderer()

	// endpoints
	{
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

	c.HTML(http.StatusOK, "templates/index.tmpl", laps)
}

func (v *View) RaceView(c *gin.Context) {
	raceID := c.Params.ByName("id")
	laps := new([]Models.Lap)

	if err := Models.GetAllLapsByRaceId(laps, raceID); err != nil {
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
		c.HTML(http.StatusOK, "race_table_view", reslt)
		return
	}
	c.HTML(http.StatusOK, "race", reslt)
}

func timestampRender(ts int64) string {
	return time.UnixMilli(ts).UTC().Format("15:04:05.000")
}
