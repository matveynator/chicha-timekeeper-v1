package Models

/**
* This controller only for API work
*/

import(
	"github.com/gin-gonic/gin"
)

// Admin login
func LoginAdmin(c *gin.Context) {

    

}

// Return list of all admins
func GetListAdmins(c *gin.Context) {
	var admins []Admin

	err := GetAllAdmins(&admins)
	if err != nil {
		c.JSON(404, nil)
		return
	}

	c.JSON(200, admins)
}

func GetAdmin(c *gin.Context) {

	var admin Admin
	id := c.Params.ByName("id")

	if err := GetOneAdmin(&admin, id); err != nil {
		c.JSON(404, nil)
		return
	}


	c.JSON(200, admin)
}

func CreateAdmin(c *gin.Context) {

	var admin Admin

	// Bind and validation
	if err := c.ShouldBind(&admin); err != nil {
		c.JSON(401, nil)
		return
	}

	// Try adding faq
	err := AddNewAdmin(&admin)
	if err != nil {
		c.JSON(401, nil)
		return
	}

	c.JSON(200, admin)
}

func UpdateAdmin(c *gin.Context) {

	var admin Admin
	id := c.Params.ByName("id")


	if err := GetOneAdmin(&admin, id); err != nil {
		c.JSON(404, nil)
		return
	}

	if err := c.ShouldBind(&admin); err != nil {
		c.JSON(401, nil)
		return
	}

	if err := PutOneAdmin(&admin); err != nil {
		c.JSON(401, nil)
		return
	}

	c.JSON(200, admin)
}

func DeleteAdmin(c *gin.Context) {

	var admin Admin
	id := c.Params.ByName("id")

	// Check admin exists
	if err := GetOneAdmin(&admin, id); err != nil {
		c.JSON(404, nil)
		return
	}

	// Try to delete
	if err := DeleteOneAdmin(&admin, id); err != nil {
		c.JSON(401, nil)
		return
	}

	c.JSON(200, nil)
}
