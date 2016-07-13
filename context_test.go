package gas

import (
	"encoding/json"
	"net/http"
	"testing"
	_ "github.com/go-gas/sessions/memory"
	"github.com/stretchr/testify/assert"
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
