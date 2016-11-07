package api

import (
	"net/http"

	"github.com/easy-takeout/easy-takeout/backend/definition"
	"gopkg.in/gin-gonic/gin.v1"
)

func (a *Api) SignIn(c *gin.Context) {
	name := c.PostForm("name")
	pass := c.PostForm("password")
	var user definition.User
	a.DB.Where(&definition.User{Name: name, Pass: pass}).First(&user)
	if user.ID != 0 {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
	}
}

func (a *Api) SignUp(c *gin.Context) {
	name := c.PostForm("name")
	pass := c.PostForm("password")

	var user definition.User

	if pass == "" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Password doesn't be blank"})
	} else if a.DB.Where(&definition.User{Name: name}).First(&user); user.ID != 0 {
		c.JSON(http.StatusForbidden, gin.H{"message": "Name is invalid or already taken"})
	} else {
		a.DB.Create(&definition.User{Name: name, Pass: pass})
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

/*
 renew password
*/
func (a *Api) RePass(c *gin.Context) {
	name := c.PostForm("name")
	pass := c.PostForm("password")
	new_pass := c.PostForm("new_password")
	new_pass_confirm := c.PostForm("new_password_confirmation")

	var user definition.User
	if new_pass == "" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Password doesn't be blank"})
	} else if new_pass != new_pass_confirm {
		c.JSON(http.StatusForbidden, gin.H{"message": "Password doesn't match the confirmation"})
	} else if a.DB.Where(&definition.User{Name: name, Pass: pass}).First(&user); user.ID == 0 {
		c.JSON(http.StatusForbidden, gin.H{"message": "Old password isn't valid"})
	} else if new_pass == pass {
		c.JSON(http.StatusForbidden, gin.H{"message": "Old password is same the new password"})
	} else {
		a.DB.Model(&user).Update(&definition.User{Pass: new_pass})
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

func init() {
	RequestBind = append(RequestBind, func(a *Api, e *gin.Engine) {
		e.POST("/sign_in", a.SignIn)
		e.POST("/sign_up", a.SignUp)
		e.POST("/renew_password", a.RePass)
	})
}
