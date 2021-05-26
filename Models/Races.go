package Models

/**
 * This controller only for API work
 */

import(
    "github.com/gin-gonic/gin"
)

// Start race
func StartRace(c *gin.Context) {

    // Search unfinished race

    // Start race

    c.JSON(200, nil)
}

// Finish race
func FinishRace(c *gin.Context) {

    // Search unfinished race

    // Finish race

    c.JSON(200, nil)
}

// Return list of all races
func GetListRaces(c *gin.Context) {
    var races []Race

	if err := GetAllRaces(&races); err != nil {
        c.JSON(404, nil)
        return
	}

    c.JSON(200, races)
}

func GetRace(c *gin.Context) {

    var race Race
    id := c.Params.ByName("id")

    if err := GetOneRace(&race, id); err != nil {
        c.JSON(404, nil)
        return
    }


    c.JSON(200, race)
}

func CreateRace(c *gin.Context) {

    var race Race

	// Bind and validation
	if err := c.ShouldBind(&race); err != nil {
		c.JSON(401, nil)
		return
	}

	// Try adding faq
	err := AddNewRace(&race)
	if err != nil {
		c.JSON(401, nil)
		return
	}

    c.JSON(200, race)
}

func UpdateRace(c *gin.Context) {

    var race Race
	id := c.Params.ByName("id")


	if err := GetOneRace(&race, id); err != nil {
		c.JSON(404, nil)
		return
	}

	if err := c.ShouldBind(&race); err != nil {
		c.JSON(401, nil)
		return
	}

	if err := PutOneRace(&race); err != nil {
        c.JSON(401, nil)
		return
	}

    c.JSON(200, race)
}

func DeleteRace(c *gin.Context) {

    var race Race
	id := c.Params.ByName("id")

	// Check race exists
	if err := GetOneRace(&race, id); err != nil {
        c.JSON(404, nil)
		return
	}

	// Try to delete
	if err := DeleteOneRace(&race, id); err != nil {
        c.JSON(401, nil)
		return
	}

    c.JSON(200, nil)
}
