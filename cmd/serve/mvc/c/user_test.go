package c

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo"
	"github.com/lpisces/bootstrap/cmd/serve"
	"github.com/lpisces/bootstrap/cmd/serve/mvc/m"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGetRegister(t *testing.T) {

	// Load default config
	serve.Conf = serve.DefaultConfig()

	req := httptest.NewRequest(echo.GET, "/register", nil)
	nr := httptest.NewRecorder()

	e := initTestEcho()
	ctx := e.NewContext(req, nr)
	//ctx.SetPath("/register")

	// Assertions
	if assert.NoError(t, GetRegister(ctx)) {
		assert.Equal(t, http.StatusOK, nr.Code)
		doc, err := goquery.NewDocumentFromReader(nr.Body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, doc.Find("title").Text(), serve.Conf.Site.Name+"-注册")

		// header
		assert.Equal(t, "注册", doc.Find("h4").Text())

		// email
		v, exist := doc.Find("input[name=email]").Attr("value")
		assert.Equal(t, true, exist)
		assert.Equal(t, "", v)

		// password
		v, exist = doc.Find("input[name=password]").Attr("value")
		assert.Equal(t, true, exist)
		assert.Equal(t, "", v)

		// password_confirm
		v, exist = doc.Find("input[name=password_confirm]").Attr("value")
		assert.Equal(t, true, exist)
		assert.Equal(t, "", v)
	}
}

func TestPostRegisterSucc(t *testing.T) {
	// Load default config
	serve.Conf = serve.DefaultConfig()

	// migrate db
	if err := m.Migrate(); err != nil {
		t.Fatal(err)
	}

	f := make(url.Values)
	f.Set("email", "iamalazyrat@gmail.com")
	f.Set("password", "helloworld")
	f.Set("password_confirm", "helloworld")

	defer func() {
		db, err := m.GetDB()
		if err != nil {
			t.Fatal(err)
		}
		db.LogMode(true)
		u := &m.User{}
		if db.Where("email = ?", "iamalazyrat@gmail.com").First(u).RecordNotFound() {
			t.Log(u.ID)
			t.Fatal(fmt.Errorf("write db failed"))
		}
		db.Delete(u)
	}()

	req := httptest.NewRequest(echo.POST, "/register", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	nr := httptest.NewRecorder()

	e := initTestEcho()
	ctx := e.NewContext(req, nr)

	// Assertions
	if assert.NoError(t, PostRegister(ctx)) {
		assert.Equal(t, http.StatusMovedPermanently, nr.Code)
	}
}

func TestPostRegisterFailed(t *testing.T) {
	// Load default config
	serve.Conf = serve.DefaultConfig()

	f := make(url.Values)
	f.Set("email", "iamalazyrat@gmailcom")
	f.Set("password", "iamalazyrat@gmailcom")
	f.Set("password_confirm", "iamalazyrat@gma")

	req := httptest.NewRequest(echo.POST, "/register", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	nr := httptest.NewRecorder()

	e := initTestEcho()
	ctx := e.NewContext(req, nr)

	// Assertions
	if assert.NoError(t, PostRegister(ctx)) {
		assert.Equal(t, http.StatusOK, nr.Code)

		doc, err := goquery.NewDocumentFromReader(nr.Body)
		if err != nil {
			t.Fatal(err)
		}

		// title
		assert.Equal(t, doc.Find("title").Text(), serve.Conf.Site.Name+"-注册")

		// header
		assert.Equal(t, "注册", doc.Find("h4").Text())

		// email
		v, exist := doc.Find("input[name=email]").Attr("value")
		assert.Equal(t, exist, true)
		assert.Equal(t, "iamalazyrat@gmailcom", v)

		emailInvalidFeedback := doc.Find("#email-invalid-feedback").Text()
		emailInvalidFeedback = strings.Replace(emailInvalidFeedback, "\n", "", -1)
		emailInvalidFeedback = strings.Replace(emailInvalidFeedback, "\t", "", -1)
		assert.Equal(t, "Email格式不正确", emailInvalidFeedback)

		// password
		v, exist = doc.Find("input[name=password]").Attr("value")
		assert.Equal(t, exist, true)
		assert.Equal(t, "", v)

		v, exist = doc.Find("input[name=password]").Attr("class")
		assert.Equal(t, exist, true)
		assert.Equal(t, true, strings.Contains(v, "is-invalid"))

		passwordInvalidFeedback := doc.Find("#password-invalid-feedback").Text()
		passwordInvalidFeedback = strings.Replace(passwordInvalidFeedback, "\n", "", -1)
		passwordInvalidFeedback = strings.Replace(passwordInvalidFeedback, "\t", "", -1)
		assert.Equal(t, "密码长度必须在6到15位之间", passwordInvalidFeedback)

		// password_confirm
		v, exist = doc.Find("input[name=password_confirm]").Attr("value")
		assert.Equal(t, exist, true)
		assert.Equal(t, "", v)

		v, exist = doc.Find("input[name=password_confirm]").Attr("class")
		assert.Equal(t, exist, true)
		assert.Equal(t, true, strings.Contains(v, "is-invalid"))

		passwordConfirmInvalidFeedback := doc.Find("#password-confirm-invalid-feedback").Text()
		passwordConfirmInvalidFeedback = strings.Replace(passwordConfirmInvalidFeedback, "\n", "", -1)
		passwordConfirmInvalidFeedback = strings.Replace(passwordConfirmInvalidFeedback, "\t", "", -1)
		assert.Equal(t, "两次输入密码不一致", passwordConfirmInvalidFeedback)
		//t.Fatal(doc.Html())
	}
}

func TestGetLogin(t *testing.T) {

	// Load default config
	serve.Conf = serve.DefaultConfig()

	req := httptest.NewRequest(echo.GET, "/login", nil)
	nr := httptest.NewRecorder()

	e := initTestEcho()
	ctx := e.NewContext(req, nr)

	// Assertions
	if assert.NoError(t, GetLogin(ctx)) {
		assert.Equal(t, http.StatusOK, nr.Code)
		doc, err := goquery.NewDocumentFromReader(nr.Body)
		if err != nil {
			t.Fatal(err)
		}

		// title
		assert.Equal(t, serve.Conf.Site.Name+"-登录", doc.Find("title").Text())

		// header
		assert.Equal(t, "登录", doc.Find("h4").Text())

		// email
		v, exist := doc.Find("input[name=email]").Attr("value")
		assert.Equal(t, true, exist)
		assert.Equal(t, "", v)
		assert.Equal(t, false, strings.Contains(v, "is-invalid"))

		// password
		v, exist = doc.Find("input[name=password]").Attr("value")
		assert.Equal(t, true, exist)
		assert.Equal(t, "", v)
		assert.Equal(t, false, strings.Contains(v, "is-invalid"))
	}
}

func TestPostLoginFailed(t *testing.T) {
	// Load default config
	serve.Conf = serve.DefaultConfig()

	f := make(url.Values)
	f.Set("email", "iamalazyrat@gmailcom")
	f.Set("password", "iamalazyrat@gmailcom")

	req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	nr := httptest.NewRecorder()

	e := initTestEcho()
	ctx := e.NewContext(req, nr)

	// Assertions
	if assert.NoError(t, PostLogin(ctx)) {
		assert.Equal(t, http.StatusOK, nr.Code)
		doc, err := goquery.NewDocumentFromReader(nr.Body)
		if err != nil {
			t.Fatal(err)
		}

		// title
		assert.Equal(t, serve.Conf.Site.Name+"-登录", doc.Find("title").Text())

		// header
		assert.Equal(t, "登录", doc.Find("h4").Text())

		// email
		v, exist := doc.Find("input[name=email]").Attr("value")
		assert.Equal(t, exist, true)
		assert.Equal(t, "iamalazyrat@gmailcom", v)

		// password
		v, exist = doc.Find("input[name=password]").Attr("value")
		assert.Equal(t, exist, true)
		assert.Equal(t, "", v)

		v, exist = doc.Find("input[name=password]").Attr("class")
		assert.Equal(t, exist, true)
		assert.Equal(t, true, strings.Contains(v, "is-invalid"))

		passwordInvalidFeedback := doc.Find("#password-invalid-feedback").Text()
		passwordInvalidFeedback = strings.Replace(passwordInvalidFeedback, "\n", "", -1)
		passwordInvalidFeedback = strings.Replace(passwordInvalidFeedback, "\t", "", -1)
		assert.Equal(t, "邮箱或密码错误", passwordInvalidFeedback)
	}

}
