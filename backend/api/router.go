package api

import (
	"github.com/bigpigeon/easy-takeout/backend/cachemanage"
	"github.com/jinzhu/gorm"
	"gopkg.in/gin-gonic/gin.v1"
)

var RequestBind = []func(*Api, *gin.Engine){}

type Api struct {
	DB        *gorm.DB
	Cache     *cachemanage.Manage
	NeedLogin bool
}

func Create(db *gorm.DB, cachemanage *cachemanage.Manage, needlogin bool) *gin.Engine {
	a := &Api{db, cachemanage, needlogin}
	e := gin.Default()
	gin.Recovery()
	for _, bind := range RequestBind {
		bind(a, e)
	}
	e.StaticFS("/", gin.Dir("./public", true))
	return e
}

func (api *Api) Close() {
	api.DB.Close()
	api.Cache.Close()
}
