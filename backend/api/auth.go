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
		c.SetCookie("user_session", "", -1, "/", "", false, true)
		c.SetCookie("user", "", -1, "/", "", false, true)
	} else {
		c.Set("authuser", name)
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
	c.SetCookie("user_session", "", -1, "/", "", false, true)
	c.SetCookie("user", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

type ReqDataSignIn struct {
	Name string `json:"name" form:"name" binding:"required"`
	Pass string `json:"password" form:"password" binding:"required"`
}

type ReqDataSignUp struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Pass        string `json:"password" form:"password" binding:"required"`
	PassConfirm string `json:"password_confirmation" form:"password_confirmation" binding:"required"`
}

type ReqDataRePass struct {
	Name           string `json:"name" form:"name" binding:"required"`
	Pass           string `json:"password" form:"password" binding:"required"`
	NewPass        string `json:"new_password" form:"new_password" binding:"required"`
	NewPassConfirm string `json:"new_password_confirmation" form:"new_password_confirmation" binding:"required"`
}

func (a *Api) SignIn(c *gin.Context) {
	var data ReqDataSignIn
	if err := c.Bind(&data); err != nil {
		c.String(http.StatusForbidden, err.Error())
		c.Abort()
	} else {
		var user definition.User
		a.DB.Where(&definition.User{Name: data.Name, Pass: data.Pass}).First(&user)
		if user.ID == 0 {
			c.String(http.StatusUnauthorized, "name or password error")
			c.Abort()
		}
	}
}

func (a *Api) SignUp(c *gin.Context) {
	var data ReqDataSignUp
	if err := c.Bind(&data); err != nil {
		c.String(http.StatusForbidden, err.Error())
		c.Abort()
	} else {
		var user definition.User
		if data.Pass != data.PassConfirm {
			c.String(http.StatusForbidden, "Password doesn't match the confirmation")
			c.Abort()
		} else if a.DB.Where(&definition.User{Name: data.Name}).First(&user); user.ID != 0 {
			c.String(http.StatusForbidden, "Name is invalid or already taken")
			c.Abort()
		} else {
			a.DB.Create(&definition.User{Name: data.Name, Pass: data.Pass})
		}
	}

}

/*
 renew password
*/
func (a *Api) RePass(c *gin.Context) {
	var data ReqDataRePass
	if err := c.Bind(&data); err != nil {
		c.String(http.StatusForbidden, err.Error())
		c.Abort()
	} else {
		var user definition.User
		if data.NewPass != data.NewPassConfirm {
			c.JSON(http.StatusForbidden, gin.H{"message": "Password doesn't match the confirmation"})
			c.Abort()
		} else if a.DB.Where(&definition.User{Name: data.Name, Pass: data.Pass}).First(&user); user.ID == 0 {
			c.JSON(http.StatusForbidden, gin.H{"message": "Old password isn't valid"})
			c.Abort()
		} else if data.NewPass == data.Pass {
			c.JSON(http.StatusForbidden, gin.H{"message": "Old password is same the new password"})
			c.Abort()
		} else {
			a.DB.Model(&user).Update(&definition.User{Pass: data.NewPass})
		}
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
