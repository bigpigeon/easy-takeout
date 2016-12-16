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
		authuser, exist := c.Get("authuser")
		if exist == false || data.AuthUser != authuser.(string) {
			c.String(http.StatusForbidden, "")
		} else {
			// update shop
			var shop definition.Shop
			a.DB.Where(&definition.Shop{
				Address: data.Address,
			}).First(&shop)
			TimeCurr := time.Now()
			if shop.Address == "" {
				shop.Address = data.Address
				a.DB.FirstOrCreate(&shop)
				// TODO access address to get shop information
			} else if shop.UpdatedAt.AddDate(0, 0, 1).Before(TimeCurr) {
				// TODO access address to get shop information
			}

			// update order
			order := definition.Order{
				Tag:      data.Tag,
				ShopAddr: data.Address,
			}
			a.DB.FirstOrCreate(&order)

			// update user item
			user_item := definition.UserItem{
				UserName: data.AuthUser,
				OrderId:  order.ID,
			}
			a.DB.FirstOrCreate(&user_item)

			// update item cell
			for _, item := range data.Items {
				item_cell := definition.UserItemCell{
					UserItemId: user_item.ID,
					Name:       item.Name,
				}
				a.DB.FirstOrCreate(&item_cell).Scopes(item_cell.Incr(item.Num))
			}

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
		group.POST("/order", a.Order)
		group.POST("/order_list", a.OrderList)
	})
}
