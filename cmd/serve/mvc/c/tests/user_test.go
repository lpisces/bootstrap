package c

import (
	//"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo"
	"github.com/lpisces/bootstrap/cmd/serve"
	"github.com/lpisces/bootstrap/cmd/serve/mvc"
	mvcc "github.com/lpisces/bootstrap/cmd/serve/mvc/c"
	//"github.com/lpisces/bootstrap/cmd/serve/mvc/m"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	//"net/url"
	//"strings"
	"github.com/gorilla/sessions"
	//"github.com/labstack/echo-contrib/session"
	"testing"
)

func TestGetRegister(t *testing.T) {

	req := httptest.NewRequest(echo.GET, "/register", nil)
	nr := httptest.NewRecorder()

	// get echo instance
	e, err := mvc.InitEcho()
	if err != nil {
		t.Fatal(err)
	}

	// set template render, use packr if not development mode
	tt, err := mvc.InitTemplate("../../../../../cmd/serve/mvc/v")
	if err != nil {
		return
	}
	e.Renderer = tt

	// context
	ctx := e.NewContext(req, nr)
	ctx.Set("_session_store", sessions.NewCookieStore([]byte("secret")))

	// Load config
	//mvc.InitConfig(ctx)

	// migrate
	mvc.InitDB()

	//ctx.SetPath("/register")

	// Assertions
	if assert.NoError(t, mvcc.GetRegister(ctx)) {
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
