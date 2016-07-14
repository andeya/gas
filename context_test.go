package gas

import (
	"encoding/json"
	"net/http"
	"testing"
	_ "github.com/go-gas/sessions/memory"
	"github.com/stretchr/testify/assert"
	"time"
)

var (
	jsonMap = JSON{
		"Test": "index page",
	}

	tstr = "Test String"

	testHTML = `<html>
    <head>
        <title>index page</title>
    </head>

    <body>
        <b>This is index page</b>
    </body>
</html>`
)

func TestRender(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	// set route
	g.Router.Get("/", func(ctx *Context) error {
		return ctx.Render(jsonMap, "testfiles/layout.html", "testfiles/index.html")
	})

	// create fasthttp.RequestHandler
	handler := g.Router.Handler

	// create httpexpect instance that will call fasthtpp.RequestHandler directly
	e := newHttpExpect(t, handler)

	// run tests
	e.GET("/").
		Expect().
		Status(http.StatusOK).
		ContentType("text/html", "utf-8").
		Body().Equal(testHTML)

}

func TestHeader(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	// set route
	g.Router.Get("/", func(ctx *Context) error {
		ctx.SetHeader("Version", "1.0")
		return ctx.STRING(http.StatusOK, "Test Header")
	})

	// create fasthttp.RequestHandler
	handler := g.Router.Handler

	// create httpexpect instance that will call fasthtpp.RequestHandler directly
	e := newHttpExpect(t, handler)

	// run tests
	e.GET("/").
		Expect().
		Status(http.StatusOK).
		ContentType("text/plain", "utf-8").
		Header("Version").Equal("1.0")

	e.GET("/").
		Expect().
		Body().Equal("Test Header")
}

func TestHTML(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	// set route
	g.Router.Get("/", func(ctx *Context) error {
		return ctx.HTML(http.StatusOK, testHTML)
	})

	// create fasthttp.RequestHandler
	handler := g.Router.Handler

	// create httpexpect instance that will call fasthtpp.RequestHandler directly
	e := newHttpExpect(t, handler)

	// run tests
	e.GET("/").
		Expect().
		Status(http.StatusOK).
		ContentType("text/html", "utf-8").
		Body().Equal(testHTML)

}

func TestSTRINGResponse(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	// set route
	g.Router.Get("/", func(ctx *Context) error {
		return ctx.STRING(http.StatusOK, tstr)
	})

	// create fasthttp.RequestHandler
	handler := g.Router.Handler

	// create httpexpect instance that will call fasthtpp.RequestHandler directly
	e := newHttpExpect(t, handler)

	// run tests
	e.GET("/").
		Expect().
		Status(http.StatusOK).
		ContentType("text/plain", "utf-8").
		Body().Equal(tstr)

}

func TestJSONResponse(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	// set route
	g.Router.Get("/", func(ctx *Context) error {
		return ctx.JSON(http.StatusOK, jsonMap)
	})

	// create fasthttp.RequestHandler
	handler := g.Router.Handler

	// create httpexpect instance that will call fasthtpp.RequestHandler directly
	e := newHttpExpect(t, handler)

	js, _ := json.Marshal(jsonMap)

	// run tests
	e.GET("/").
		Expect().
		Status(http.StatusOK).
		ContentType("application/json", "utf-8").
		Body().Equal(string(js))

}

func TestContext_SetCookie(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")


	// set route
	g.Router.Get("/", func(ctx *Context) error {
		ctx.SetCookie("cookie-key", "cookie-value")

		return ctx.STRING(http.StatusOK, "set cookie")
	})

	// create fasthttp.RequestHandler
	handler := g.Router.Handler

	// create httpexpect instance that will call fasthtpp.RequestHandler directly
	e := newHttpExpect(t, handler)

	// run tests
	cookie := e.GET("/").
	Expect().
	Status(http.StatusOK).
	Cookie("cookie-key")
	cookie.Value().Equal("cookie-value")
}

func TestContext_SetCookieBytes(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")


	// set route
	g.Router.Get("/", func(ctx *Context) error {
		ctx.SetCookieBytes([]byte("cookie-key2"), []byte("cookie-value2"))

		return ctx.STRING(http.StatusOK, "set cookie byte")
	})

	// create fasthttp.RequestHandler
	handler := g.Router.Handler

	// create httpexpect instance that will call fasthtpp.RequestHandler directly
	e := newHttpExpect(t, handler)

	// run tests
	cookie := e.GET("/").
	Expect().
	Status(http.StatusOK).
	Cookie("cookie-key2")
	cookie.Value().Equal("cookie-value2")
}

func TestContext_SetCookieByConfig(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")


	// set route
	g.Router.Get("/", func(ctx *Context) error {
		cfg := &CookieSettings{
			PathByte: []byte("/123"),
			DomainByte: []byte("example.com"),
			Expired: 123456,
			HttpOnly: false,
		}

		ctx.SetCookieByConfig(cfg, "cookie-key3", "cookie-value3")

		return ctx.STRING(http.StatusOK, "set cookie")
	})

	// create fasthttp.RequestHandler
	handler := g.Router.Handler

	// create httpexpect instance that will call fasthtpp.RequestHandler directly
	e := newHttpExpect(t, handler)

	// run tests
	n := time.Now()
	cookie := e.GET("/").
	Expect().
	Status(http.StatusOK).
	Cookie("cookie-key3")
	cookie.Value().Equal("cookie-value3")
	cookie.Domain().Equal("example.com")
	cookie.Path().Equal("/123")
	cookie.Expires().InRange(n, n.Add(time.Second * 123459))
	httponly := cookie.Raw().HttpOnly
	assert.Equal(t, false, httponly)
}

func TestContext_SetCookieByConfigWithBytes(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")


	// set route
	g.Router.Get("/", func(ctx *Context) error {
		cfg := &CookieSettings{
			PathByte: []byte("/123"),
			DomainByte: []byte("example.com"),
			Expired: 123456,
			HttpOnly: false,
		}

		ctx.SetCookieByConfigWithBytes(cfg, []byte("cookie-key4"), []byte("cookie-value4"))

		return ctx.STRING(http.StatusOK, "set cookie")
	})

	// create fasthttp.RequestHandler
	handler := g.Router.Handler

	// create httpexpect instance that will call fasthtpp.RequestHandler directly
	e := newHttpExpect(t, handler)

	// run tests
	n := time.Now()
	cookie := e.GET("/").
	Expect().
	Status(http.StatusOK).
	Cookie("cookie-key4")
	cookie.Value().Equal("cookie-value4")
	cookie.Domain().Equal("example.com")
	cookie.Path().Equal("/123")
	cookie.Expires().InRange(n, n.Add(time.Second * 123459))
	httponly := cookie.Raw().HttpOnly
	assert.Equal(t, false, httponly)
}

func TestContext_SessionStart(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")


	// set route
	g.Router.Get("/", func(ctx *Context) error {
		ctx.SessionStart()

		assert.NotNil(t, ctx.sessionManager)

		return ctx.STRING(http.StatusOK, "session start")
	})

	// create fasthttp.RequestHandler
	handler := g.Router.Handler

	// create httpexpect instance that will call fasthtpp.RequestHandler directly
	e := newHttpExpect(t, handler)

	// run tests
	e.GET("/").
	Expect().
	Status(http.StatusOK).
	Body().Equal("session start")
}

func TestContext_SessionDestroy(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")


	// set route
	g.Router.Get("/", func(ctx *Context) error {
		ctx.sessionManager = nil
		ctx.cookieHandler = nil

		ctx.SessionDestroy()

		assert.NotNil(t, ctx.sessionManager)
		assert.NotNil(t, ctx.cookieHandler)

		return ctx.STRING(http.StatusOK, "session destroy")
	})

	// create fasthttp.RequestHandler
	handler := g.Router.Handler

	// create httpexpect instance that will call fasthtpp.RequestHandler directly
	e := newHttpExpect(t, handler)

	// run tests
	e.GET("/").
	Expect().
	Status(http.StatusOK).
	Body().Equal("session destroy")
}
