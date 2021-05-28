package Models

/**
 * This controller only for API work
 */

import(
    "github.com/gin-gonic/gin"
)

// Active race (used as storage for active Race)
var ActiveRace Race

// Start race
func StartRace(c *gin.Context) {

    // Search unfinished race
    var unfinishedRace Race
    if err := GetOneUnfinishedRace(&unfinishedRace); err == nil {
        c.JSON(401, map[string]string{"error" : "Complete an active ride first"})
        return
    }

    // Start race
    var race Race
    id := c.Params.ByName("id")
    if err := GetOneRace(&race, id); err != nil {
        c.JSON(404, map[string]string{"error" : "Race not found"})
        return
    }

    // Store active race in variable
    ActiveRace = race

    c.JSON(200, nil)
}

// Finish race
func FinishRace(c *gin.Context) {

    // Search unfinished race
    var unfinishedRaces []Race
    if err := GetAllUnfinishedRace(&unfinishedRaces); err != nil {
        c.JSON(401, map[string]string{"error" : "Complete an active ride first"})
        return
    }

    // Finish race
    for _,race := range unfinishedRaces {
        race.IsActive = false

        if err := PutOneRace(&race); err != nil {
            c.JSON(500, map[string]string{"error" : "Server cannot finish an active race"})
            return
        }
    }

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
