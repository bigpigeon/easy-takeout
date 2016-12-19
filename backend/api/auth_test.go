package api

import (
	"net/http"
	"net/http/httptest"

	"os"
	"strings"
	"testing"

	"github.com/easy-takeout/easy-takeout/backend/cachemanage"
	"github.com/easy-takeout/easy-takeout/backend/definition"
	"gopkg.in/gin-gonic/gin.v1"
)

func reqSignUp(
	t *testing.T,
	user, pass, pass_confirm string,
	router *gin.Engine,
) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	signUpForm := getFormData(map[string]string{
		"name":                  user,
		"password":              pass,
		"password_confirmation": pass_confirm,
	})
	req, _ := http.NewRequest("POST", "/sign_up", strings.NewReader(signUpForm))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	return w
}

func reqSignIn(
	t *testing.T,
	user, pass string,
	router *gin.Engine,
) *httptest.ResponseRecorder {

	w := httptest.NewRecorder()
	signInForm := getFormData(map[string]string{
		"name":     user,
		"password": pass,
	})
	req, _ := http.NewRequest("POST", "/sign_in", strings.NewReader(signInForm))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
	return w
}

func reqSignOut(
	t *testing.T,
	user, pass string,
	router *gin.Engine,
) *httptest.ResponseRecorder {
	//signout  must on authorized status
	w := reqSignIn(t, user, pass, router)
	// read sign in session
	promiseRespCodeRight(t, w, []int{3})
	resp := http.Response{Header: w.Header()}
	SignInCookie := []*http.Cookie{}
	for _, c := range resp.Cookies() {
		switch c.Name {
		case "user":
			if c.Value != "pigeon" {
				t.Log("cookie was error")
				t.FailNow()
			}
			SignInCookie = append(SignInCookie, c)
		case "user_session":
			SignInCookie = append(SignInCookie, c)
		}
	}
	// sign out request
	signOutRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/sign_out", nil)
	for _, c := range SignInCookie {
		req.AddCookie(c)
	}
	router.ServeHTTP(signOutRecorder, req)
	return signOutRecorder
}

func reqRePass(
	t *testing.T,
	user, pass, new_pass, new_pass_confirm string,
	router *gin.Engine,
) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	form := getFormData(map[string]string{
		"name":                      user,
		"password":                  pass,
		"new_password":              new_pass,
		"new_password_confirmation": new_pass_confirm,
	})
	req, _ := http.NewRequest("POST", "/renew_password", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	router.ServeHTTP(w, req)
	return w
}

func TestAuthNormal(t *testing.T) {
	db, err := definition.Connect("sqlite3", "test/test.db", nil)
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
	//sign up with diff password
	promiseRespCodeRight(t, reqSignUp(t, "pigeon", "abc123", "abc321", r), []int{4})
	//sign up right
	promiseRespCodeRight(t, reqSignUp(t, "pigeon", "abc123", "abc123", r), []int{3})
	//sign up with existed
	promiseRespCodeRight(t, reqSignUp(t, "pigeon", "abc123", "abc123", r), []int{4})
	//sign in with wrong password
	promiseRespCodeRight(t, reqSignIn(t, "pigeon", "abc321", r), []int{4})
	//sign in right
	promiseRespCodeRight(t, reqSignIn(t, "pigeon", "abc123", r), []int{3})
	//sign_out
	promiseRespCodeRight(t, reqSignOut(t, "pigeon", "abc123", r), []int{3})
}

func TestAuthRePass(t *testing.T) {
	db, err := definition.Connect("sqlite3", "test/test.db", nil)
	checkFailErr(t, err)
	err = definition.Migrate(db)
	checkFailErr(t, err)

	cache, err := cachemanage.Create("sqlite3", "test/test.db", nil)
	checkFailErr(t, err)
	defer func() {
		db.Close()
		os.Remove("test/test.db")
	}()

	r := Create(db, cache, true)
	reqSignUp(t, "tom", "ac3", "ac3", r)
	// renew password with same password
	promiseRespCodeRight(t, reqRePass(t, "tom", "ac3", "ac3", "ac3", r), []int{4})
	// renew password with diff new password
	promiseRespCodeRight(t, reqRePass(t, "tom", "ac3", "ac2", "ab2", r), []int{4})
	// renew password with not existed
	promiseRespCodeRight(t, reqRePass(t, "joy", "ac3", "ac2", "ac2", r), []int{4})
	// renew password
	promiseRespCodeRight(t, reqRePass(t, "tom", "ac3", "ac2", "ac2", r), []int{3})
	// a old password login
	promiseRespCodeRight(t, reqSignIn(t, "tom", "ac3", r), []int{4})
	// a new password login
	promiseRespCodeRight(t, reqSignIn(t, "tom", "ac2", r), []int{3})

}
