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
	//assert := assert.New(t)

	// new gas
	g := New("testfiles/config_test.yaml")

	// set route
	g.Router.Get("/", indexPage)

	e := newHttpExpect(t, g.Router.Handler)
	e.GET("/").Expect().Status(http.StatusOK).Body().Equal(indexString)

	// create request
	//req, _ := http.NewRequest("GET", "/", nil)
	//rec := httptest.NewRecorder()
	//g.Router.ServeHTTP(rec, req)
	//
	//assert.Equal(http.StatusOK, rec.Code)
	//assert.Equal(indexString, rec.Body.String())

}

//func Testgas_Static(t *testing.T) {
//	assert := assert.New(t)
//
//	// new gas
//	g := New("testfiles/config_test.yaml")
//
//	// set route
//	//g.Router.Get("/", indexPage)
//
//	// create request
//	req, _ := http.NewRequest("GET", "/testfiles/static.txt", nil)
//	rec := httptest.NewRecorder()
//	g.Router.ServeHTTP(rec, req)
//
//	assert.Equal(http.StatusOK, rec.Code)
//	assert.Equal(testStaticString, rec.Body.String())
//}
//
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
