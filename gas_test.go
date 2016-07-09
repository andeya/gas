package gas

import (
	"github.com/gavv/httpexpect"
	"github.com/go-gas/gas/model/MySQL"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"net/http"
	"testing"
)

var (
	indexString      = "indexpage"
	testStaticString = "This is a static file"
)

func newHttpExpect(t *testing.T, h fasthttp.RequestHandler) *httpexpect.Expect {
	// create fasthttp.RequestHandler
	//handler := g.Router.Handler

	// create httpexpect instance that will call fasthtpp.RequestHandler directly
	e := httpexpect.WithConfig(httpexpect.Config{
		Reporter: httpexpect.NewAssertReporter(t),
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(h),
			Jar:       httpexpect.NewJar(),
		},
	})

	return e
}

func Testgas(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	// set route
	g.Router.Get("/", indexPage)

	e := newHttpExpect(t, g.Router.Handler)
	e.GET("/").
		Expect().
		Status(http.StatusOK).
		Body().Equal(indexString)
}

func indexPage(ctx *Context) error {
	return ctx.STRING(http.StatusOK, indexString)
}

func TestGas_NewModel(t *testing.T) {
	as := assert.New(t)

	// new gas
	g := New("testfiles/config_test.yaml")
	m := g.NewModel()

	as.IsType(&MySQLModel.MySQLModel{}, m)
}

func BenchmarkGas(b *testing.B) {
	b.ReportAllocs()

	// new gas
	g := New("testfiles/config_test.yaml")

	// set route
	g.Router.Get("/", indexPage)

	req := fasthttp.Request{}
	req.SetRequestURI("/")
	req.Header.SetMethod("GET")


	for i := 0; i < b.N; i++ {
		ctx := fasthttp.RequestCtx{
			Request: req,
		}

		g.Router.Handler(&ctx)
	}
}
