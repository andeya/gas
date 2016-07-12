package gas

import (
	"encoding/json"
	"errors"
	"github.com/buaazp/fasthttprouter"
	"github.com/go-gas/gas/model"
	"github.com/valyala/fasthttp"
	"html/template"
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
}

// create context
//func createContext(w *ResponseWriter, r *http.Request, g *Engine) *Context {
func createContext(r *fasthttp.RequestCtx, g *Engine) *Context {
	c := &Context{}
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
	jsonstr, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, errr := ctx.Write(jsonstr)

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
