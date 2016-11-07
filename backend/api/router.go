package api

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gin-gonic/gin.v1"
)

var RequestBind = []func(*Api, *gin.Engine){}

type Api struct {
	DB        *gorm.DB
	NeedLogin bool
}

func Create(db *gorm.DB, needlogin bool) *gin.Engine {
	a := &Api{db, needlogin}
	e := gin.Default()
	for _, bind := range RequestBind {
		bind(a, e)
	}
	return e
}
