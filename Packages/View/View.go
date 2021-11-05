package view

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"
	"sort"
	"strconv"

	"html/template"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"

	"chicha/Packages/Models"
	"chicha/Packages/View/SSE"
)

type View struct {
	static embed.FS
}

func (v *View) setupRenderer() multitemplate.Renderer {
	f := template.FuncMap{
		"timestampRender":      timestampRender,
		"millisDurationRender": millisDurationRender,
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

// return fs for serve static files
func (v *View) getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(v.static, "static/assets")
	if err != nil {
		log.Fatal(err)
	}
	return http.FS(fsys)
}

func New(r *gin.Engine, static embed.FS, ch <-chan struct{}) *View {
	v := &View{static: static}
	r.HTMLRender = v.setupRenderer()

	// endpoints
	{
		// static files
		r.StaticFS("/static/assets/", v.getFileSystem())

		r.GET("/", v.Homepage)
		r.GET("/race/:id", v.RaceView)
		r.GET("/race/:id/:tag", v.RaceRiderView)

		rStream := r.Group("/race-stream")
		sse.Setup(rStream, ch)
	}

	return v
}

func (v *View) Homepage(c *gin.Context) {
	laps := new([]Models.Lap)
	lap := new(Models.Lap)
	//language=SQL
	s := `
	SELECT
	race_id, 
	MIN(discovery_unix_time) as discovery_unix_time, 
	MIN(discovery_time) as discovery_time,
	MAX(lap_is_current) as lap_is_current
	FROM laps
	GROUP BY race_id
	ORDER BY race_id desc
	`

	if err := Models.DB.Raw(s).Find(laps).Error; err != nil {
		c.Error(err)
		log.Println("",err)
		//return
	}
	if err := Models.DB.Raw(s).First(lap).Error; err != nil {
		c.Error(err)
		log.Println(err)
		//return
	}

	c.HTML(http.StatusOK, "index", gin.H{
		"currentRace": lap,
		"raceList":    laps,
	})
}


func (v *View) RaceRiderView(c *gin.Context) {
	raceID, _ := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	tagID := c.Params.ByName("tag")

	laps := Models.GetLapsForWeb(uint(raceID))

	var sLaps []gin.H

	sort.Slice(laps, func(i, j int) bool {
		//sort by minimal CurrentRacePosition
		return laps[i].LapNumber < laps[j].LapNumber
	})

	for _, v := range laps {

		if v.TagID == tagID {

			var stl string
			if v.BetterOrWorseLapTime > 0 {
				// if (+) sec worse then prev
				// orange
				stl = "orange"
			} else if v.BetterOrWorseLapTime < 0 && v.BestLapPosition != 1 {
				// if (-) sec better then prev
				// green
				stl = "green"
			} else if ( v.BestLapPosition == 1 && v.BestLapNumber == v.LapNumber )  {
				// if best race lap
				// violet
				stl = "violet"
			}

			sLaps = append(sLaps, gin.H{
				"Lap":   v,
				"Style": stl,
			})

		}
	}

	reslt := gin.H{
		"RaceID": raceID,
		"Laps":   sLaps,
	}

	if c.Query("updtable") == "true" {
		c.HTML(http.StatusOK, "race_table_view", reslt)
		return
	}
	c.HTML(http.StatusOK, "race", reslt)
}



func (v *View) RaceView(c *gin.Context) {
	raceID, _ := strconv.ParseInt(c.Params.ByName("id"), 10, 64)

	laps := Models.GetLapsForWeb(uint(raceID))

	var sLaps []gin.H

	sort.Slice(laps, func(i, j int) bool {
		//sort by minimal CurrentRacePosition
		return laps[i].CurrentRacePosition < laps[j].CurrentRacePosition
	})

	for _, v := range laps {

		if v.LapIsCurrent == 1 {

			var stl string
			if v.BetterOrWorseLapTime > 0 {
				// if (+) sec worse then prev
				// orange
				stl = "orange"
			} else if v.BetterOrWorseLapTime < 0 {
				// if (-) sec better then prev
				// green
				stl = "green"
			} else if ( v.BestLapPosition == 1 && v.BestLapNumber == v.LapNumber )  {
				// if best race lap
				// violet
				stl = "violet"
			}

			sLaps = append(sLaps, gin.H{
				"Lap":   v,
				"Style": stl,
			})

 		}
	}

	reslt := gin.H{
		"RaceID": raceID,
		"Laps":   sLaps,
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

func millisDurationRender(ts int64) string {
	//return float64(ts)/1000
	//return time.Duration(ts) * time.Millisecond
	duration := time.Duration(ts) * time.Millisecond
	if ts > 0 {
		return fmt.Sprintf("+%s", duration.String())
	} else {
		return duration.String()
	}
}
