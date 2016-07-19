# Gas

<img src="https://raw.githubusercontent.com/go-gas/gas/master/logo.png" alt="go-gas" />

[![Build Status](https://travis-ci.org/go-gas/gas.svg?branch=master)](https://travis-ci.org/go-gas/gas) [![codecov](https://codecov.io/gh/go-gas/gas/branch/master/graph/badge.svg)](https://codecov.io/gh/go-gas/gas) [![Go Report Card](https://goreportcard.com/badge/github.com/go-gas/gas)](https://goreportcard.com/report/github.com/go-gas/gas)
[![Join the chat at https://gitter.im/go-gas/gas](https://badges.gitter.im/go-gas/gas.svg)](https://gitter.im/go-gas/gas?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Gas aims to be a high performance, full-featured, easy to use, and quick develop backend web application framework in Golang.

## Features

- Router (based on [fasthttprouter](https://github.com/buaazp/fasthttprouter) package)
- Easy to use golang template engine. (will include another template engine)
- Context (easy to manage the request, response and session)
- Middleware (Global and specify routing path middleware support)
- Logger package [gas-logger](https://github.com/go-gas/logger)
- Read config from a yaml file [gas-config](https://github.com/go-gas/config)
- Database model (developing, based on [go-gas/SQLBuilder](https://github.com/go-gas/SQLBuilder))
- Support listen HTTP/HTTPS and UNIX addr.

other features are highly active development

##### and you can see example at [gas-example](https://github.com/go-gas/example).

## Install

```
$ go get github.com/go-gas/gas
```

## How to use

### Micro service

If you want to create a micro service, you can write all of things in one package, for example:

    |-- $GOPATH
    |   |-- src
    |       |--Your_Project_Name
    |          |-- main.go
    |          |-- config.yaml

main.go

```go
package main

import (
  "github.com/go-gas/gas"
  "net/http"
)

func main() {
  g := gas.Default("config.yaml")

  g.Router.Get("/", Index)
  g.Router.Get("/user", GetUser)

  g.Run()
}

func Index(ctx *gas.Context) error {
  return ctx.HTML(http.StatusOK, "Micro service! <br> <a href=\"/user\">json response example</a>")
}

func GetUser(ctx *gas.Context) error {
  return ctx.JSON(http.StatusOK, gas.H{
    "name": "John",
    "age":  32,
  })
}
```

see [go-gas/example/micro_demo](https://github.com/go-gas/example/tree/master/micro_demo)

### Large project architecture

Write all code in one file is so dirty and hard to maintain when you build a large site include many controller, middlerware... and so on.

So, maybe we can seperate them in many packages. (directorys)

#### file structure

    |-- $GOPATH
    |   |-- src
    |       |--Your_Project_Name
    |          |-- config
    |              |-- default.yaml
    |          |-- controllers
    |              |-- default.go
    |          |-- log
    |          |-- models
    |          |-- routers
    |              |-- routers.go
    |          |-- static
    |          |-- views
    |          |-- main.go


#### main.go

##### 1. import

```go
import (
    "Your_Project_Name/routers"
    "github.com/go-gas/gas"
    "github.com/go-gas/gas/middleware"
)
```

##### 2. New or Default

`Default` returns an Engine instance with the `Logger` middleware already attached.

```go
g := gas.Default()
g.LoadConfig("your/config/path")
```

`New` returns a new blank Engine instance without any middleware attached.

```go
g := gas.New()
g.LoadConfig("your/config/path")
```

If you don't want to load any config,
you might be know that gas have a default config with

```go
var defaultConfig = map[interface{}]interface{}{
	"Mode":       "DEV",
	"ListenAddr": "localhost",
	"ListenPort": "8080",
	"PubDir":     "public",
	"Db": map[interface{}]interface{}{
		"SqlDriver": "MySQL",
		"Hostname":  "localhost",
		"Port":      "3306",
		"Username":  "root",
		"Password":  "",
		"Charset":   "utf8",
	},
}
```

or you can give config path when new gas app

```go
g := gas.New("config/path1", "config/path2")
```

##### 3. Register Routes

```go
routers.RegistRout(g.Router)
```

Then in your routers/routers.go

```go
package routers

import (
    "Your_Project_Name/controllers"
    "github.com/go-gas/gas"
)

func RegistRout(r *gas.Router)  {

    r.Get("/", controllers.IndexPage)
    r.Post("/post/:param", controllers.PostTest)

    rc := &controllers.RestController{}
    r.REST("/User", rc)

}
```

##### 4. Using gas.Context

###### Cookie

gas context provides set cookie with string or []bytes, you don't need to do type assertion.

```go
ctx.SetCookie("key", "value")
ctx.SetCookieBytes([]byte("key"), []byte("value"))
```

and you can set cookie detail using gas.CookieSettings struct
```go
type CookieSettings struct {
	PathByte []byte
	PathString string

	DomainByte []byte
	DomainString string

	Expired int
	HttpOnly bool
}
```

for example:

```go
cfg := &CookieSettings{
	PathByte: []byte("/somePath"),
	DomainByte: []byte("example.com"),
	Expired: 3600,
	HttpOnly: true,
}

ctx.SetCookieByConfig(cfg, "key", "value")
```

###### Session

First, you must import provider package before call SessionStart.

```go
import _ "github.com/go-gas/sessions/memory"
```

And the provider package will Register to session package automaticly.

Then call SessionStart, and use Get and Set to operate your session.

```go
s := ctx.SessionStart()
s.Set("key", "value")
d := s.Get("value").(string)
```

The Get function will return interface{} type, you must know the data type and do type assertion your self.

##### 5. Register middleware

###### Global middleware

If you want a middleware to be run during every request to your application,
you can use Router.Use function to register your middleware.

```go
g.Router.Use(middleware.LogMiddleware)
```

###### Assigning middleware to Route

If you want to assign middleware to specific routes,
you can set your middlewares after set route function like:

```go
r.Get("/", controllers.IndexPage, myMiddleware1, myMiddleware2)
```

And you can write your own middleware function

```go
func LogMiddleware(next gas.GasHandler) gas.GasHandler {
    return func (c *gas.Context) error  {

       // do something before next handler

       err := next(c)

       // do something after next handler

       return err
    }
}
```

or

```go
func MyMiddleware2 (ctx *gas.Context) error {
  // do something
}
```

#### The final step

Run and listen your web application with default `8080` port.

```go
g.Run()
```

you can give listen address and another port.

```go
g.Run(":8089")
```

Serving HTTPS (secure) requests.

```go
g.RunTLS(":8080", "CertFile", "CertKey")
```

HTTP requests from the given UNIX addr.

```go
g.RunUNIX("/tmp/gas.sock", 0644)
```

but I recommend setting listen address in config files.

## Benchmark

Using [go-web-framework-benchmark](https://github.com/smallnest/go-web-framework-benchmark) to benchmark with another web fframework.

<img src="https://raw.githubusercontent.com/go-gas/go-web-framework-benchmark/master/benchmark.png" alt="go-gas-benchmark" />

### Benchmark-alloc

<img src="https://raw.githubusercontent.com/go-gas/go-web-framework-benchmark/master/benchmark_alloc.png" alt="go-gas-benchmark-alloc" />

### Benchmark-latency

<img src="https://raw.githubusercontent.com/go-gas/go-web-framework-benchmark/master/benchmark_latency.png" alt="go-gas-benchmark-latency" />

### Benchmark-pipeline

<img src="https://raw.githubusercontent.com/go-gas/go-web-framework-benchmark/master/benchmark-pipeline.png" alt="go-gas-benchmark-pipeline" />

## Concurrency

<img src="https://raw.githubusercontent.com/go-gas/go-web-framework-benchmark/master/concurrency.png" alt="go-gas-concurrency" />

### Concurrency-alloc

<img src="https://raw.githubusercontent.com/go-gas/go-web-framework-benchmark/master/concurrency_alloc.png" alt="go-gas-concurrency-alloc" />

### Concurrency-latency

<img src="https://raw.githubusercontent.com/go-gas/go-web-framework-benchmark/master/concurrency_latency.png" alt="go-gas-concurrency-latency" />

### Concurrency-pipeline

<img src="https://raw.githubusercontent.com/go-gas/go-web-framework-benchmark/master/concurrency-pipeline.png" alt="go-gas-concurrency-pipeline" />

## Roadmap

- [ ] Router
 - [ ] Group Routing
- [ ] Models
 - [ ] Model fields mapping
 - [ ] ORM
 - [ ] Relation mapping
 - [x] Transaction
 - [ ] QueryBuilder
- [ ] Session
 - [ ] Filesystem
 - [ ] Database
 - [ ] Redis
 - [ ] Memcache
 - [x] In memory
- [ ] Cache
 - [ ] Memory
 - [ ] File
 - [ ] Redis
 - [ ] Memcache
- [ ] i18n
- [x] HTTPS
- [ ] Command line tools
- [ ] Form handler
- [ ] Security check features(csrf, xss filter...etc)
- [ ] Doker Image
