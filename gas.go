// Package gas is a web framework.
//
// Example
//
// Your project file structure
//  |-- $GOPATH
//  |   |-- src
//  |       |--Your_Project_Name
//  |          |-- config
//  |              |-- default.yaml
//  |          |-- controllers
//  |              |-- default.go
//  |          |-- log
//  |          |-- models
//  |          |-- routers
//  |              |-- routers.go
//  |          |-- static
//  |          |-- views
//  |          |-- main.go
// main.go
//  import (
//  	"Your_Project_Name/routers"
// 	"github.com/go-gas/gas"
//  )
//
//  // Create gas object with config path
//  // default is config/default.yaml
//  g := gas.New("config/path")
//
//  // register route
//  routers.RegistRout(g.Router)
//
//  // run and listen
//  g.Run()
// routers.go
//  import (
//  	"Your_Project_Name/controllers"
//  	"github.com/go-gas/gas"
//  )
//
//  func RegistRout(r *Engine.Router)  {
//
//  	r.Get("/", controllers.IndexPage)
//  	r.Post("/post/:param", controllers.PostTest)
//
//  	rc := &controllers.RestController{}
//  	r.REST("/User", rc)
//
//  }
// controllers.go
//  package controllers
//
//  import (
//  	"github.com/go-gas/gas"
//  )
//
//  func IndexPage(ctx *Engine.Context) error {
//  	return ctx.Render("", "views/layout.html", "views/index.html")
//  }
//
//  func PostTest(ctx *Engine.Context) error {
//  	a := map[string]string{
//  		"Name": ctx.GetParam("param"),
//  	}
//
//  	return ctx.Render(a, "views/layout2.html")
//  }
//
// rest_controller.go
//  import (
//  	"github.com/go-gas/gas"
//  )
//
//  type RestController struct {
//  	gas.ControllerInterface
//  }
//
//  func (rc *RestController) Get(c *Engine.Context) error {
//
//  	return c.STRING(http.StatusOK, "Test Get")
//  }
//
//  func (rc *RestController) Post(c *Engine.Context) error {
//
//  	return c.STRING(http.StatusOK, "Test Post")
//  }
package gas

import (
	"fmt"
	"github.com/go-gas/config"
	"github.com/go-gas/gas/logger"
	"github.com/go-gas/gas/model"
	"github.com/go-gas/gas/model/MySQL"
	"github.com/valyala/fasthttp"
	"os"
	"strings"
	"sync"
)

var defaultConfig = map[interface{}]interface{}{
	"Mode":       "DEV",
	"ListenAddr": "localhost",
	"ListenPort": "8080",
	"PubDir":     "public",
	"CrtFile":    "",
	"KeyFile":    "",
	"Db": map[interface{}]interface{}{
		"SqlDriver": "MySQL",
		"Hostname":  "localhost",
		"Port":      "3306",
		"Username":  "root",
		"Password":  "",
		"Charset":   "utf8",
	},
	"session": map[interface{}]interface{}{
		"CookieName":      "gas-session",
		"SessionIDLength": 10,
		"CookieLifeTime":  0,    // cookie expire time
		"ExpireTime":      3600, // time for destroy from GC
		"GCTime":          60,   // GC frequency
		"Domain":          "",
	},
	"sessionProvider":       "memory",
	"sessionProviderConfig": map[interface{}]interface{}{},
}

type (
	Engine struct {
		Router *Router
		Config *config.Engine
		Model  *gasModel
		pool   sync.Pool
		Logger *logger.Logger
	}

	gasModel struct {
		model.Model
	}
)

// New gas Object
//
// Ex:
//  g := New()
//  g.Run()
func New(configPath ...string) *Engine {
	g := &Engine{}

	// init logger
	if _, err := os.Stat("log/system.log"); os.IsNotExist(err) {
		os.Mkdir("log", 0700)
	}

	g.Logger = logger.New("log/system.log")

	// init pool
	g.pool.New = func() interface{} {
		c := createContext(nil, g)

		return c
	}

	// load config
	g.Config = config.New(defaultConfig)
	if len(configPath) != 0 {
		for _, path := range configPath {
			g.Config.Load(path)
		}
	}

	// set router
	g.Router = newRouter(g) //&Router{g: g}

	// set default not found handler
	g.Router.SetNotFoundHandler(defaultNotFoundHandler)

	// set default panic handler
	g.Router.SetPanicHandler(defaultPanicHandler)

	// set static file path
	g.Router.StaticPath(g.Config.GetString("PubDir"))

	return g
}

func defaultNotFoundHandler(c *Context) error {
	return c.STRING(404, "Page Not Found.")
}

func defaultPanicHandler(c *Context, rcv interface{}) error {
	logStr := fmt.Sprintf("Panic occurred...rcv: %v", rcv)
	c.gas.Logger.Error(logStr)

	var output string
	if c.gas.Config.Get("Mode") == "DEV" {
		output = logStr
	} else {
		output = "Sorry...some error occurred..."
	}

	return c.STRING(500, output)
}

// Load config from file
func (g *Engine) LoadConfig(configPath string) {
	g.Config.Load(configPath)
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
func (g *Engine) Run(addr ...string) (err error) {
	listenAddr := ""
	if len(addr) == 0 {
		listenAddr = g.Config.GetString("ListenAddr") + ":" + g.Config.GetString("ListenPort")
	} else {
		listenAddr = addr[0]
	}

	fmt.Println("Server is Listen on: " + listenAddr)

	err = fasthttp.ListenAndServe(listenAddr, g.Router.Handler)
	return
}

// RunTLS attaches the router to a http.Server and starts listening and serving HTTPS (secure) requests.
func (g *Engine) RunTLS(addr ...string) (err error) {
	var certFile, keyFile, listenAddr string

	if len(addr) == 0 {
		listenAddr = g.Config.GetString("ListenAddr") + ":" + g.Config.GetString("ListenPort")
		certFile = g.Config.GetString("CertFile")
		keyFile = g.Config.GetString("KeyFile")
	} else {
		listenAddr = addr[0]
		certFile = addr[1]
		keyFile = addr[2]
	}

	fmt.Println("Server is Listen on: " + listenAddr)

	err = fasthttp.ListenAndServeTLS(listenAddr, certFile, keyFile, g.Router.Handler)
	return
}

// New model according to config settings
func (g *Engine) NewModel() model.ModelInterface {
	// get db
	// db := g.NewDb()
	c := g.Config

	//var db model.SlimDbInterface
	var m model.ModelInterface
	//var builder model.BuilderInterface

	switch strings.ToLower(c.GetString("Db.SqlDriver")) {
	case "mysql":
		//db = new(model.MysqlDb)
		//m = new(model.MySQLModel)
		//builder = new(model.MySQLBuilder)
		//m = model.New(c)
		m = MySQLModel.New(c)
	default:
		panic("Unknow Database Driver: " + c.GetString("Db.SqlDriver"))

	}

	//err := db.ConnectWithConfig(g.Config.Db)
	//if err != nil {
	//	panic(err.Error())
	//}
	//m.SetDB(db)
	//builder.SetDB(db)
	//m.SetBuilder(builder)

	return m
}
