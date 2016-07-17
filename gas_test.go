package gas

import (
	"crypto/tls"
	"github.com/gavv/httpexpect"
	"github.com/go-gas/gas/model/MySQL"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

var (
	indexString = "indexpage"
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

func testRequest(t *testing.T, url string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	defer resp.Body.Close()

	assert.NoError(t, err)

	b, ioerr := ioutil.ReadAll(resp.Body)
	assert.NoError(t, ioerr)
	assert.Equal(t, "200 OK", resp.Status, "should get a 200")
	assert.Equal(t, indexString, string(b))
}

func indexPage(ctx *Context) error {
	return ctx.STRING(http.StatusOK, indexString)
}

func TestGas(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	// set route
	g.Router.Get("/", indexPage)

	e := newHttpExpect(t, g.Router.Handler)
	e.GET("/").
		Expect().
		Status(http.StatusOK).
		Body().Equal(indexString)

	// test 404 not found
	e.GET("/gas").
		Expect().
		Status(http.StatusNotFound).
		Body().Equal(default404Body)
}

func TestRun(t *testing.T) {
	g := New()

	// set route
	g.Router.Get("/", indexPage)

	go func() {
		assert.NoError(t, g.Run())
	}()
	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	testRequest(t, "http://localhost:8080")
}

func TestRunWithDefine(t *testing.T) {
	g := New()

	// set route
	g.Router.Get("/", indexPage)

	go func() {
		assert.NoError(t, g.Run(":9000"))
	}()
	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	testRequest(t, "http://localhost:9000")
}

func TestRunWithDefault(t *testing.T) {
	g := Default()

	// set route
	g.Router.Get("/", indexPage)

	go func() {
		assert.NoError(t, g.Run(":9001"))
	}()
	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	testRequest(t, "http://localhost:9001")

	e := newHttpExpect(t, g.Router.Handler)
	// test X-Real-IP
	e.GET("/").
		WithHeader("X-Real-IP", "192.168.1.1").
		Expect().
		Status(http.StatusOK).
		Body().Equal(indexString)

	// test X-Forwarded-For
	e.GET("/").
		WithHeader("X-Forwarded-For", "192.168.1.2").
		Expect().
		Status(http.StatusOK).
		Body().Equal(indexString)
}

func TestRunTLS(t *testing.T) {
	g := New()

	// set route
	g.Router.Get("/", indexPage)

	go func() {
		assert.NoError(t, g.RunTLS("localhost:8081", "certificate/localhost.cert", "certificate/localhost.key"))
	}()
	time.Sleep(5 * time.Millisecond)

	testRequest(t, "https://localhost:8081")
}

func TestRunTLSWithConfig(t *testing.T) {
	g := New()

	g.LoadConfig("testfiles/config_test2.yaml")

	// set route
	g.Router.Get("/", indexPage)

	go func() {
		assert.NoError(t, g.RunTLS())
	}()
	time.Sleep(5 * time.Millisecond)

	testRequest(t, "https://localhost:8089")
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
