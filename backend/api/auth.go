package api

import (
	"net/http"

	"github.com/easy-takeout/easy-takeout/backend/definition"
	"gopkg.in/gin-gonic/gin.v1"
)

func (a *Api) AuthRequired(c *gin.Context) {
	name, _ := c.Cookie("user")
	session, _ := c.Cookie("user_session")
	s, err := a.Cache.GetLoginSession(name)
	if err != nil {
		panic(err)
	}

	if s == "" || s != session {
		c.AbortWithStatus(http.StatusForbidden)
	}
}

func (a *Api) Authorized(c *gin.Context) {
	name := c.PostForm("name")
	session, err := a.Cache.GenerateLoginSession(name)
	if err != nil {
		panic(err)
	}
	c.SetCookie("user_session", session, 0, "/", "", false, true)
	c.SetCookie("user", name, 0, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func (a *Api) UnAuthorized(c *gin.Context) {
	c.SetCookie("user_session", "", -1, "/", "", true, true)
	c.SetCookie("user", "", -1, "/", "", true, true)
	//	c.Redirect(http.StatusFound, "/")
}

func (a *Api) SignIn(c *gin.Context) {
	name := c.PostForm("name")
	pass := c.PostForm("password")
	var user definition.User
	a.DB.Where(&definition.User{Name: name, Pass: pass}).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "name or password error"})
		c.Abort()
	}
}

func (a *Api) SignUp(c *gin.Context) {
	name := c.PostForm("name")
	pass := c.PostForm("password")

	var user definition.User

	if pass == "" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Password doesn't be blank"})
		c.Abort()
	} else if a.DB.Where(&definition.User{Name: name}).First(&user); user.ID != 0 {
		c.JSON(http.StatusForbidden, gin.H{"message": "Name is invalid or already taken"})
		c.Abort()
	} else {
		a.DB.Create(&definition.User{Name: name, Pass: pass})
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
		c.Abort()
	} else if new_pass != new_pass_confirm {
		c.JSON(http.StatusForbidden, gin.H{"message": "Password doesn't match the confirmation"})
		c.Abort()
	} else if a.DB.Where(&definition.User{Name: name, Pass: pass}).First(&user); user.ID == 0 {
		c.JSON(http.StatusForbidden, gin.H{"message": "Old password isn't valid"})
		c.Abort()
	} else if new_pass == pass {
		c.JSON(http.StatusForbidden, gin.H{"message": "Old password is same the new password"})
		c.Abort()
	} else {
		a.DB.Model(&user).Update(&definition.User{Pass: new_pass})
	}
}

func init() {
	RequestBind = append(RequestBind, func(a *Api, e *gin.Engine) {
		group := e.Group("/")
		group.POST("/sign_in", a.SignIn, a.Authorized)
		group.POST("/sign_up", a.SignUp, a.Authorized)
		group.POST("/renew_password", a.RePass, a.Authorized)
		group.POST("/sign_out", a.AuthRequired, a.UnAuthorized)
	})
}
