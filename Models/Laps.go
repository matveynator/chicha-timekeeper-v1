package Models

/**
* This controller only for API work
*/

import(
	"fmt"
	"github.com/gin-gonic/gin"
)

// Return list of all laps
func GetLapsByRaceId(c *gin.Context) {
        var laps []LapSmall
        race_id := c.Params.ByName("id")
        err := GetAllLapsByRaceId(&laps, race_id)
	fmt.Println(err)
        if err != nil {
		c.JSON(404, err)
                return
        }

	//for _, lap := range laps {
    	//	// element is the element from someSlice for where we are
	//	fmt.Printf("%d, %d, %s, %s \n", lap.RaceID, lap.LapNumber, lap.DiscoveryTime, lap.TagID)
	//fmt.Printf("%#v", lap)
    	//}
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
	var laps Lap

	err := GetLastLap(&laps)
	if err != nil {
		c.JSON(404, nil)
		return
	}

	c.JSON(200, laps)
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
