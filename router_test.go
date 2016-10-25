package gas

import (
	"net/http"
	"testing"
)

func TestRouter_Static(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	e := newHttpExpect(t, g.Router.Handler)
	e.GET("/testfiles/static.txt").Expect().
		Status(http.StatusOK).
		Body().Equal("This is a static file")
}

func TestRouter_Get(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.Get("/test", func(c *Context) error {
		return c.STRING(http.StatusOK, "TEST")
	})

	e := newHttpExpect(t, g.Router.Handler)
	e.GET("/test").Expect().Status(http.StatusOK).Body().Equal("TEST")
}

func TestRouter_Post(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.Post("/test", func(c *Context) error {
		return c.STRING(http.StatusOK, c.GetParam("Test"))
	})

	e := newHttpExpect(t, g.Router.Handler)
	ee := e.POST("/test").WithFormField("Test", "POSTDATA").Expect()
	ee.Status(http.StatusOK)
	ee.Body().Equal("POSTDATA")
}

func TestRouter_Put(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.Put("/test", func(c *Context) error {
		return c.STRING(http.StatusOK, c.GetParam("Test"))
	})

	e := newHttpExpect(t, g.Router.Handler)
	ee := e.PUT("/test").WithFormField("Test", "POSTDATA").Expect()
	ee.Status(http.StatusOK)
	ee.Body().Equal("POSTDATA")
}

func TestRouter_Patch(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.Patch("/", func(c *Context) error {
		return c.STRING(http.StatusOK, c.GetParam("Test"))
	})

	e := newHttpExpect(t, g.Router.Handler)
	ee := e.PATCH("/").WithFormField("Test", "POSTDATA").Expect()
	ee.Status(http.StatusOK)
	ee.Body().Equal("POSTDATA")
}

func TestRouter_Delete(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.Delete("/", func(c *Context) error {
		return c.STRING(http.StatusOK, "Deleted")
	})

	e := newHttpExpect(t, g.Router.Handler)
	ee := e.DELETE("/").Expect()
	ee.Status(http.StatusOK)
	ee.Body().Equal("Deleted")
}

func TestRouter_Options(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.Options("/", func(c *Context) error {
		return c.STRING(http.StatusOK, "Option")
	})

	e := newHttpExpect(t, g.Router.Handler)
	ee := e.OPTIONS("/").Expect()
	ee.Status(http.StatusOK)
	ee.Body().Equal("Option")
}

func TestRouter_Head(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.Head("/", func(c *Context) error {
		return c.STRING(http.StatusOK, "Head")
	})

	e := newHttpExpect(t, g.Router.Handler)
	ee := e.HEAD("/").Expect()
	ee.Status(http.StatusOK)
	ee.Body().Equal("Head")
}

type testController struct {
	ControllerInterface
}

func (cn *testController) Get(c *Context) error {
	return c.STRING(http.StatusOK, "Get Test")
}
func (cn *testController) Post(c *Context) error {
	return c.STRING(http.StatusOK, "Post Test"+c.GetParam("Test"))
}

func TestRouter_REST(t *testing.T) {
	var c = &testController{}

	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.REST("/User", c)

	e := newHttpExpect(t, g.Router.Handler)

	ee1 := e.GET("/User").Expect()
	ee1.Status(http.StatusOK)
	ee1.Body().Equal("Get Test")

	ee2 := e.POST("/User").WithFormField("Test", "POSTED").Expect()
	ee2.Status(http.StatusOK)
	ee2.Body().Equal("Post TestPOSTED")
}

func TestRouter_SetMiddlewareFunc(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.Get("/test", func(c *Context) error {
		return c.STRING(http.StatusOK, "TEST")
	}, testMiddleware1)

	e := newHttpExpect(t, g.Router.Handler)
	e.GET("/test").WithFormField("Test", "Go").
		Expect().Status(http.StatusOK).Body().Equal("TEST")

	e.GET("/test").WithFormField("Test", "DontGo").
		Expect().Status(http.StatusForbidden).Body().Equal("ERROR")
}

func testMiddleware1(next GasHandler) GasHandler {
	return func(ctx *Context) error {
		if ctx.GetParam("Test") == "Go" {
			return next(ctx)
		}

		return ctx.STRING(http.StatusForbidden, "ERROR")
	}
}

func testMiddleware2(ctx *Context) error {
	if ctx.GetParam("Test") == "Go" {
		ctx.Request.PostArgs().Add("FromMiddleware", "200")
		return ctx.STRING(http.StatusOK, "OK")
	}

	ctx.Request.PostArgs().Add("FromMiddleware", "NO")
	return ctx.STRING(http.StatusForbidden, "ERROR-")
}

func TestRouter_SetGasHandlerAsMiddleware(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.Get("/test", func(c *Context) error {
		_, err := c.WriteString(c.GetParam("FromMiddleware"))
		return err
	}, testMiddleware2)

	e := newHttpExpect(t, g.Router.Handler)
	e.GET("/test").WithFormField("Test", "Go").
		Expect().Status(http.StatusOK).Body().Equal("OK200")

	e.GET("/test").WithFormField("Test", "DontGo").
		Expect().Status(http.StatusForbidden).Body().Equal("ERROR-NO")
}
