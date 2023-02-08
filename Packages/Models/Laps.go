package Models

/**
* This controller only for API work
 */

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// Return list of all laps
func GetLapsByRaceId(c *gin.Context) {
	var laps []Lap
	race_id := c.Params.ByName("id")
	err := GetAllLapsByRaceId(&laps, race_id)
	if err != nil {
		c.JSON(404, err)
		return
	}
	c.JSON(200, laps)
}

// Return list of all laps
func GetResultsByRaceId(c *gin.Context) {
	var laps []Lap
	race_id, _ := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	err := GetAllResultsByRaceId(&laps, uint(race_id))
	if err != nil {
		c.JSON(404, err)
		return
	}
	c.JSON(200, laps)
}

// Return list of all laps
func GetListLaps(c *gin.Context) {
	var laps []Lap
	err := GetAllLaps(&laps)
	if err != nil {
		c.JSON(404, nil)
		return
	}
	c.JSON(200, laps)
}

// Return list current raceid
func GetLastLapData(c *gin.Context) {
	var lap Lap
	err := GetLastLap(&lap)
	if err != nil {
		c.JSON(404, nil)
		return
	}
	c.JSON(200, lap)
}

func GetLap(c *gin.Context) {
	var lap Lap
	id := c.Params.ByName("id")
	if err := GetOneLap(&lap, id); err != nil {
		c.JSON(404, nil)
		return
	}
	c.JSON(200, lap)
}

// Return list of all laps
func GetLapsByTagId(c *gin.Context) {
	var laps []Lap
	id := c.Params.ByName("id")
	err := GetAllLapsByTagId(&laps, id)
	if err != nil {
		c.JSON(404, nil)
		return
	}
	c.JSON(200, laps)
}

func CreateLap(c *gin.Context) {

	var lap Lap

	// Bind and validation
	if err := c.ShouldBind(&lap); err != nil {
		c.JSON(401, nil)
		return
	}

	// Try adding faq
	err := AddNewLap(&lap)
	if err != nil {
		c.JSON(401, nil)
		return
	}

	c.JSON(200, lap)
}

func UpdateLap(c *gin.Context) {

	var lap Lap
	id := c.Params.ByName("id")

	if err := GetOneLap(&lap, id); err != nil {
		c.JSON(404, nil)
		return
	}

	if err := c.ShouldBind(&lap); err != nil {
		c.JSON(401, nil)
		return
	}

	if err := PutOneLap(&lap); err != nil {
		c.JSON(401, nil)
		return
	}

	c.JSON(200, lap)
}

func DeleteLap(c *gin.Context) {

	var lap Lap
	id := c.Params.ByName("id")

	// Check lap exists
	if err := GetOneLap(&lap, id); err != nil {
		c.JSON(404, nil)
		return
	}

	// Try to delete
	if err := DeleteOneLap(&lap, id); err != nil {
		c.JSON(401, nil)
		return
	}

	c.JSON(200, nil)
}
