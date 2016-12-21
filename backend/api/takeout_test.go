package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"time"

	"github.com/easy-takeout/easy-takeout/backend/cachemanage"
	"github.com/easy-takeout/easy-takeout/backend/definition"
	"gopkg.in/gin-gonic/gin.v1"
)

func reqJson(
	t *testing.T,
	data map[string]interface{},
	cookies []*http.Cookie, req_url string,
	router *gin.Engine,
) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	jsondata := getJsonData(t, data)
	req, _ := http.NewRequest("POST", req_url, strings.NewReader(jsondata))
	for _, c := range cookies {
		req.AddCookie(c)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
	return w
}

func reqTakeout(
	t *testing.T,
	user, pass, authUser, address, tag string, items []TakeOutItem,
	router *gin.Engine,
) *httptest.ResponseRecorder {
	//takeout  must on authorized status
	w := reqSignIn(t, authUser, pass, router)
	promiseRespCodeRight(t, w, []int{3})
	resp := http.Response{Header: w.Header()}
	SignInCookie := resp.Cookies()

	data := map[string]interface{}{
		"authuser": authUser,
		"address":  address,
		"tag":      tag,
		"items":    items,
	}
	return reqJson(t, data, SignInCookie, "/takeout", router)
}

func reqOrder(
	t *testing.T,
	orderId uint,
	router *gin.Engine,
) *httptest.ResponseRecorder {
	data := map[string]interface{}{
		"order_id": orderId,
	}
	return reqJson(t, data, []*http.Cookie{}, "/order", router)
}

func reqOrderList(
	t *testing.T,
	start, end time.Time,
	router *gin.Engine,
) *httptest.ResponseRecorder {
	data := map[string]interface{}{
		"start": start,
		"end":   end,
	}
	return reqJson(t, data, []*http.Cookie{}, "/order_list", router)
}

func TestTakeoutNormal(t *testing.T) {
	db, err := definition.Connect("sqlite3", "test/test.db", nil)
	//	db = db.Debug()
	checkFailErr(t, err)

	err = definition.Migrate(db)
	checkFailErr(t, err)
	defer func() {
		db.Close()
		os.Remove("test/test.db")
	}()

	cache, err := cachemanage.Create("sqlite3", "test/test.db", nil)
	checkFailErr(t, err)

	r := Create(db, cache, true)
	//make a user
	promiseRespCodeRight(t, reqSignUp(t, "jia", "123", "123", r), []int{3})

	var takeoutResp RespDataTakeout
	// a normal takeout
	{
		w := reqTakeout(t,
			"jia", "123",
			"jia", "waimai.baidu.com/waimai/shop/123", "", []TakeOutItem{{"手撕鸡", 10, 2}, {"饭", 2, 1}},
			r,
		)
		promiseRespCodeRight(t, w, []int{2})

		checkFailErr(t, json.Unmarshal(w.Body.Bytes(), &takeoutResp))
	}
	// get order
	{
		w := reqOrder(t, takeoutResp.OrderId, r)
		promiseRespCodeRight(t, w, []int{2})
		//		t.Log(fmtJsonBytes(t, w.Body.Bytes()))
	}
	// get order list
	{
		w := reqOrderList(t, time.Now().AddDate(0, 0, -1), time.Now(), r)
		promiseRespCodeRight(t, w, []int{2})
		//		t.Log(fmtJsonBytes(t, w.Body.Bytes()))
	}
}
