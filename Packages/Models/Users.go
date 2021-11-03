package Models

/**
 * This controller only for API work
 */

import(
    "github.com/gin-gonic/gin"
)


// Return list of all users
func GetListUsers(c *gin.Context) {
    var users []User

	if err := GetAllUsers(&users); err != nil {
        c.JSON(404, nil)
        return
	}

    c.JSON(200, users)
}

func GetUser(c *gin.Context) {

    var user User
    id := c.Params.ByName("id")

    if err := GetOneUser(&user, id); err != nil {
        c.JSON(404, nil)
        return
    }


    c.JSON(200, user)
}

func CreateUser(c *gin.Context) {

    var user User

	// Bind and validation
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(401, nil)
		return
	}

	// Try adding faq
	err := AddNewUser(&user)
	if err != nil {
		c.JSON(401, nil)
		return
	}

    c.JSON(200, user)
}

func UpdateUser(c *gin.Context) {

    var user User
	id := c.Params.ByName("id")


	if err := GetOneUser(&user, id); err != nil {
		c.JSON(404, nil)
		return
	}

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(401, nil)
		return
	}

	if err := PutOneUser(&user); err != nil {
        c.JSON(401, nil)
		return
	}

    c.JSON(200, user)
}

func DeleteUser(c *gin.Context) {

    var user User
	id := c.Params.ByName("id")

	// Check user exists
	if err := GetOneUser(&user, id); err != nil {
        c.JSON(404, nil)
		return
	}

	// Try to delete
	if err := DeleteOneUser(&user, id); err != nil {
        c.JSON(401, nil)
		return
	}

    c.JSON(200, nil)
}
