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
		return c.STRING(200, "TEST")
	})

	e := newHttpExpect(t, g.Router.Handler)
	e.GET("/test").Expect().Status(http.StatusOK).Body().Equal("TEST")
}

func TestRouter_Post(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.Post("/test", func(c *Context) error {
		return c.STRING(200, c.GetParam("Test"))
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
		return c.STRING(200, c.GetParam("Test"))
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
		return c.STRING(200, c.GetParam("Test"))
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
		return c.STRING(200, "Deleted")
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
		return c.STRING(200, "Option")
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
		return c.STRING(200, "Head")
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
	return c.STRING(200, "Get Test")
}
func (cn *testController) Post(c *Context) error {
	return c.STRING(200, "Post Test"+c.GetParam("Test"))
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
