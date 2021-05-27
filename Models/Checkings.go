package Models

/**
 * This controller only for API work
 */

import(
    "github.com/gin-gonic/gin"
)


// Return list of all checkins
func GetListCheckins(c *gin.Context) {
    var checkins []Checkin

	if err := GetAllCheckins(&checkins); err != nil {
        c.JSON(404, nil)
        return
	}

    c.JSON(200, checkins)
}

func GetCheckin(c *gin.Context) {

    var checkin Checkin
    id := c.Params.ByName("id")

    if err := GetOneCheckin(&checkin, id); err != nil {
        c.JSON(404, nil)
        return
    }


    c.JSON(200, checkin)
}

func CreateCheckin(c *gin.Context) {

    var checkin Checkin

	// Bind and validation
	if err := c.ShouldBind(&checkin); err != nil {
		c.JSON(401, nil)
		return
	}

	// Try adding faq
	err := AddNewCheckin(&checkin)
	if err != nil {
		c.JSON(401, nil)
		return
	}

    c.JSON(200, checkin)
}

func UpdateCheckin(c *gin.Context) {

    var checkin Checkin
	id := c.Params.ByName("id")


	if err := GetOneCheckin(&checkin, id); err != nil {
		c.JSON(404, nil)
		return
	}

	if err := c.ShouldBind(&checkin); err != nil {
		c.JSON(401, nil)
		return
	}

	if err := PutOneCheckin(&checkin); err != nil {
        c.JSON(401, nil)
		return
	}

    c.JSON(200, checkin)
}

func DeleteCheckin(c *gin.Context) {

    var checkin Checkin
	id := c.Params.ByName("id")

	// Check checkin exists
	if err := GetOneCheckin(&checkin, id); err != nil {
        c.JSON(404, nil)
		return
	}

	// Try to delete
	if err := DeleteOneCheckin(&checkin, id); err != nil {
        c.JSON(401, nil)
		return
	}

    c.JSON(200, nil)
}
