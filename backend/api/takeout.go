package api

import (
	"net/http"

	"time"

	"github.com/easy-takeout/easy-takeout/backend/definition"
	"github.com/jinzhu/gorm"
	"gopkg.in/gin-gonic/gin.v1"
)

type TakeOutItem struct {
	Name string `json:"name"`
	Num  int    `json:"num"`
}

type ReqDataTakeOut struct {
	AuthUser string        `json:"authuser" binding:"required"`
	Address  string        `json:"address" binding:"required"`
	Tag      string        `json:"tag"`
	Items    []TakeOutItem `json:"items"`
}

type ReqDataOrder struct {
	OrderId uint `json:"order_id" binding:"required"`
}

type ReqDataOrderList struct {
	Start time.Time `json:"start" binding:"required"`
	End   time.Time `json:"end" binding:"required"`
}

type RespDataOrderList struct {
	gorm.Model
	ShopId uint
	Tag    string
	User   definition.User
}

func (a *Api) Takeout(c *gin.Context) {
	var data ReqDataTakeOut
	if err := c.BindJSON(&data); err != nil {
		c.String(http.StatusForbidden, err.Error())
		c.Abort()
	} else {
		authuser, _ := c.Get("authuser")
		if data.AuthUser != authuser.(string) {
			c.String(http.StatusForbidden, "")
		} else {

		}
	}
}

func (a *Api) Order(c *gin.Context) {
	var data ReqDataOrder
	if err := c.BindJSON(&data); err != nil {
		c.String(http.StatusForbidden, err.Error())
		c.Abort()
	} else {
		var order definition.Order
		a.DB.Scopes(
			definition.PreloadOrder.All,
		).Where(&definition.Order{Model: gorm.Model{ID: data.OrderId}}).First(&order)
		c.JSON(http.StatusOK, order)
	}
}

func (a *Api) OrderList(c *gin.Context) {
	var data ReqDataOrderList
	if err := c.BindJSON(&data); err != nil {
		c.String(http.StatusForbidden, err.Error())
		c.Abort()
	} else {
		var orders []RespDataOrderList
		a.DB.Model(&definition.Order{}).Scopes(
			definition.BetweenCreateTime(data.Start, data.End),
			definition.PreloadOrder.BelongsTo,
		).Scan(&orders)
		c.JSON(200, orders)
	}
}

func init() {
	RequestBind = append(RequestBind, func(a *Api, e *gin.Engine) {
		group := e.Group("/")
		group.POST("/takeout", a.Authorized, a.Takeout)
		group.POST("/order_list", a.Authorized, a.OrderList)
	})
}
