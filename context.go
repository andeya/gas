package gas

import (
	"encoding/json"
	"errors"
	"github.com/buaazp/fasthttprouter"
	"github.com/go-gas/gas/model"
	"github.com/valyala/fasthttp"
	"html/template"
	"github.com/go-gas/sessions"
	"time"
)

// D is JSON Data Type
type JSON map[string]interface{}

type Context struct {
	//context.Context
	*fasthttp.RequestCtx

	//RespWriter *ResponseWriter
	//Req        *fasthttp.Request
	ps *fasthttprouter.Params

	// handlerFunc CHandler

	gas *Engine //

	// DB
	isUseDB bool
	mobj    model.ModelInterface

	// cookie
	defaultCookieConfig *CookieSettings

	// session
	isUseSession bool
	sessionManager *sessions.SessionManager
	cookieHandler sessions.HTTPCookieHandlerInterface
}

type CookieSettings struct {
	PathByte []byte
	PathString string

	DomainByte []byte
	DomainString string

	Expired int
	HttpOnly bool
}

// create context
//func createContext(w *ResponseWriter, r *http.Request, g *Engine) *Context {
func createContext(r *fasthttp.RequestCtx, g *Engine) *Context {
	c := &Context{
		defaultCookieConfig: &CookieSettings{
			PathByte: []byte("/"),
			Expired: 60 * 60 * 24, // one day
			HttpOnly: true,
		},
	}
	//c.RespWriter = w
	c.RequestCtx = r
	c.gas = g

	return c
}

// reset context when get it from buffer
//func (ctx *Context) reset(w http.ResponseWriter, r *http.Request, g *Engine) {
//	ctx.Req = r
//func (ctx *Context) reset(w http.ResponseWriter, r *http.Request, g *Goslim) {
func (ctx *Context) reset(fctx *fasthttp.RequestCtx, ps *fasthttprouter.Params, g *Engine) {

	//ctx.Req = fctx.Request
	//ctx.RespWriter.reset(w)
	ctx.RequestCtx = fctx
	ctx.ps = ps
	ctx.gas = g

	ctx.mobj = nil
	ctx.isUseDB = false

	ctx.isUseSession = false
	ctx.cookieHandler = nil
}

// func (ctx *Context) Next()  {
//     ctx.handlerFunc(ctx)
// }

// Get parameter from post or get value
func (ctx *Context) GetParam(name string) string {
	//if ctx.Req.PostForm == nil || ctx.Req.Form == nil {
	//	ctx.Req.ParseForm()
	//}
	//
	//if v := ctx.Req.FormValue(name); v != "" {
	//	return v
	//}

	if fv := ctx.FormValue(name); fv != nil {
		return string(fv)
	}

	return ctx.ps.ByName(name)
}

//func (ctx *Context) GetFormValue(name string) string {
//	if fv := ctx.FormValue(name); fv != nil {
//		return string(fv)
//	}
//
//	return ""
//}

// func (ctx *Context) GetAllParams()  {
//     res := make(map[string]string, 0)

//     for key, v := ctx.Req.Form {
//         res[key] = v[0]
//     }
// }

// Render function combined data and template to show
func (ctx *Context) Render(data interface{}, tplPath ...string) error {
	if len(tplPath) == 0 {
		return errors.New("File path can not be empty")
	}

	ctx.SetContentType(TextHTMLCharsetUTF8)

	// tpls := strings.Join(tplPath, ", ")
	tmpl := template.New("gas")

	for _, tpath := range tplPath {
		tmpl = template.Must((tmpl.ParseFiles(tpath)))
	}

	err := tmpl.Execute(ctx, data)

	return err
	// if err != nil {
	//     // println(err)
	//     // panic(err)

	//     return err
	// }

	// return nil
}

// Set the response data-type to html
func (ctx *Context) HTML(code int, html string) error {

	ctx.SetContentType(TextHTMLCharsetUTF8)
	ctx.SetStatusCode(code)         // .RespWriter.WriteHeader(code)
	_, err := ctx.WriteString(html) //_, err := ctx.RespWriter.Write([]byte(html))

	return err
}

// Set the response data-type to plain text
func (ctx *Context) STRING(status int, data string) error {

	//ctx.RespWriter.Header().Set(ContentType, TextPlainCharsetUTF8)
	//ctx.RespWriter.WriteHeader(status)
	//_, err := ctx.RespWriter.Write([]byte(data))

	if ctx.IsGet() {
		ctx.SetContentType(TextPlainCharsetUTF8)
	} else {
		ctx.SetContentType(ApplicationForm)
	}
	ctx.SetStatusCode(status)
	_, err := ctx.WriteString(data)
	return err
}

// Set response data-type to json and try to json encode
func (ctx *Context) JSON(status int, data interface{}) error {

	ctx.SetContentType(ApplicationJSONCharsetUTF8)
	ctx.SetStatusCode(status)

	// encode json string
	jsonByte, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, errr := ctx.Write(jsonByte)

	return errr
}

func (ctx *Context) SetHeader(key, value string) {
	ctx.Response.Header.Set(key, value)
}

// Get model using context in controller
func (ctx *Context) GetModel() model.ModelInterface {
	m := ctx.gas.NewModel()

	if m != nil {
		ctx.isUseDB = true
		ctx.mobj = m

		return m
	}

	return nil
}

// Close db connection
func (ctx *Context) CloseDB() error {
	return ctx.mobj.Builder().GetDB().Close()
}


// ==== cookie ====

func (ctx *Context) SetCookie(key, value string) {
	c := ctx.generateCookieFromConfig(ctx.defaultCookieConfig)

	c.SetKey(key)
	c.SetValue(value)

	ctx.Request.Header.SetCookie(key, value)
	ctx.Response.Header.SetCookie(c)
}

func (ctx *Context) SetCookieBytes(key, value []byte) {
	c := ctx.generateCookieFromConfig(ctx.defaultCookieConfig)

	c.SetKeyBytes(key)
	c.SetValueBytes(value)

	ctx.Request.Header.SetCookieBytesKV(key, value)
	ctx.Response.Header.SetCookie(c)
}

func (ctx *Context) SetCookieByConfig(cfg *CookieSettings, key, value string) {
	c := ctx.generateCookieFromConfig(cfg)

	c.SetKey(key)
	c.SetValue(value)

	ctx.Request.Header.SetCookie(key, value)
	ctx.Response.Header.SetCookie(c)
}

func (ctx *Context) SetCookieByConfigWithBytes(cfg *CookieSettings, key, value []byte) {
	c := ctx.generateCookieFromConfig(cfg)

	c.SetKeyBytes(key)
	c.SetValueBytes(value)

	ctx.Request.Header.SetCookieBytesKV(key, value)
	ctx.Response.Header.SetCookie(c)
}

func (ctx *Context) generateCookieFromConfig(cfg *CookieSettings) *fasthttp.Cookie {
	c := fasthttp.AcquireCookie()
	c.Reset()

	if len(cfg.PathByte) != 0 {
		c.SetPathBytes(cfg.PathByte)
	} else if cfg.PathString != "" {
		c.SetPath(cfg.PathString)
	}

	if len(cfg.DomainByte) != 0 {
		c.SetDomainBytes(cfg.DomainByte)
	} else if cfg.DomainString != "" {
		c.SetDomain(cfg.DomainString)
	}

	if cfg.Expired != 0 {
		c.SetExpire(time.Now().Add(time.Duration(cfg.Expired) * time.Second))
	}

	c.SetHTTPOnly(cfg.HttpOnly)

	return c
}

func (ctx *Context) GetCookie(key string) []byte {
	return ctx.Request.Header.Cookie(key)
}


// ==== session management  ====

func (ctx *Context) SessionStart() sessions.SessionInterface {
	// read session provider from config
	if ctx.sessionManager == nil {
		sc := &sessions.SessionConfig{}
		ctx.sessionManager = sessions.New(ctx.gas.Config.GetString("sessionProvider"), ctx.gas.Config.GetStruct("session", sc).(*sessions.SessionConfig))
	}

	cookieHandler := sessions.NewFasthttpCookieHandler(ctx.RequestCtx)

	s, _ := ctx.sessionManager.SessionStart(cookieHandler)

	ctx.isUseSession = true
	ctx.cookieHandler = cookieHandler

	return s
}

func (ctx *Context) SessionDestroy() {
	if ctx.sessionManager == nil {
		ctx.SessionStart()
	}

	if ctx.cookieHandler == nil {
		cookieHandler := sessions.NewFasthttpCookieHandler(ctx.RequestCtx)

		ctx.isUseSession = true
		ctx.cookieHandler = cookieHandler
	}

	ctx.sessionManager.Destroy(ctx.cookieHandler)
}

func (ctx *Context) SessionEnd() {
	if ctx.cookieHandler != nil {
		sessions.RecycleFasthttpCookieHandler(ctx.cookieHandler.(*sessions.FasthttpCookieHandler))
	}
}
