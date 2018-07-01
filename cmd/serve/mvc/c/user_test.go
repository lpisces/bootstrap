package c

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo"
	"github.com/lpisces/bootstrap/cmd/serve"
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
		v, exist := doc.Find("input[name=email]").Attr("value")
		assert.Equal(t, exist, true)
		assert.Equal(t, v, "")
	}
}

func TestPostRegisterFailed(t *testing.T) {
	// Load default config
	serve.Conf = serve.DefaultConfig()

	f := make(url.Values)
	f.Set("email", "iamalazyrat@gmailcom")

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
		assert.Equal(t, doc.Find("title").Text(), serve.Conf.Site.Name+"-注册")
		v, exist := doc.Find("input[name=email]").Attr("value")
		assert.Equal(t, exist, true)
		assert.Equal(t, "iamalazyrat@gmailcom", v)

		emailInvalidFeedback := doc.Find("#email-invalid-feedback").Text()
		emailInvalidFeedback = strings.Replace(emailInvalidFeedback, "\n", "", -1)
		emailInvalidFeedback = strings.Replace(emailInvalidFeedback, "\t", "", -1)
		assert.Equal(t, "iamalazyrat@gmailcom does not validate as email", emailInvalidFeedback)

		passwordInvalidFeedback := doc.Find("#password-invalid-feedback").Text()
		passwordInvalidFeedback = strings.Replace(passwordInvalidFeedback, "\n", "", -1)
		passwordInvalidFeedback = strings.Replace(passwordInvalidFeedback, "\t", "", -1)
		assert.Equal(t, "non zero value required", passwordInvalidFeedback)
		//t.Fatal(doc.Html())
	}
}