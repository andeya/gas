# Gas

<img src="https://raw.githubusercontent.com/go-gas/gas/master/logo.jpg" alt="go-gas" width="200px" />

[![Build Status](https://travis-ci.org/go-gas/gas.svg?branch=master)](https://travis-ci.org/go-gas/gas) [![codecov](https://codecov.io/gh/go-gas/gas/branch/master/graph/badge.svg)](https://codecov.io/gh/go-gas/gas) [![Go Report Card](https://goreportcard.com/badge/github.com/go-gas/gas)](https://goreportcard.com/report/github.com/go-gas/gas)
[![Join the chat at https://gitter.im/go-gas/gas](https://badges.gitter.im/go-gas/gas.svg)](https://gitter.im/go-gas/gas?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Gas is a high performance, full-featured(in the future), easy to use, and quick develop backend web apllication framework in Golang.
 
# Features

- Router (based on [fasthttprouter](https://github.com/buaazp/fasthttprouter) package)
- Easy to use golang template engine. (will include another template engine) 
- Context (easy to manage the request and response)
- Middleware (Global and specify routing path middleware support)
- Log package
- Read config from a yaml file [gas-config](https://github.com/go-gas/Config)
- Database model (developing)

other features are highly active development 

##### and you can see example at [gas-example](https://github.com/go-gas/example).

# Install

```
$ go get github.com/go-gas/gas
```

# Run demo

```
$ git clone https://github.com/go-gas/example.git && cd example
$ go run main.go
```

# Your project file structure

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

# Quick start

### Import
```go
import (
    "Your_Project_Name/routers"
    "github.com/go-gas/gas"
    "github.com/go-gas/gas/middleware"
)
```

### New

```go
g := gas.New() // will load "config/default.yaml"
```

or

```go
g := gas.New("config/path")
```

### Register Routes

```go
routers.RegistRout(g.Router)
```
Then in your routers.go

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

### Register middleware

##### Global middleware
If you want a middleware to be run during every request to your application,
you can use Router.Use function to register your middleware.

```go
g.Router.Use(middleware.LogMiddleware)
```

##### Assigning middleware to Route
If you want to assign middleware to specific routes,
you can set your middlewares after set route function like:

```go
r.Get("/", controllers.IndexPage, myMiddleware1, myMiddleware2)
```

##### And you can write your own middleware function

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

### The final step

Run and listen your web application
```go
g.Run()
```

or you can give listen address
```go
g.Run(":8080")
```

but I recommend setting listen address in config files.

# Benchmark

Using [go-web-framework-benchmark](https://github.com/smallnest/go-web-framework-benchmark) to benchmark with another web fframework.

<img src="https://raw.githubusercontent.com/go-gas/go-web-framework-benchmark/master/benchmark.png" alt="go-gas-benchmark" />

#### Benchmark-alloc

<img src="https://raw.githubusercontent.com/go-gas/go-web-framework-benchmark/master/benchmark_alloc.png" alt="go-gas-benchmark-alloc" />

#### Benchmark-latency

<img src="https://raw.githubusercontent.com/go-gas/go-web-framework-benchmark/master/benchmark_latency.png" alt="go-gas-benchmark-latency" />

#### Benchmark-pipeline

<img src="https://raw.githubusercontent.com/go-gas/go-web-framework-benchmark/master/benchmark-pipeline.png" alt="go-gas-benchmark-pipeline" />

## Concurrency

<img src="https://raw.githubusercontent.com/go-gas/go-web-framework-benchmark/master/concurrency.png" alt="go-gas-concurrency" />

#### Concurrency-alloc

<img src="https://raw.githubusercontent.com/go-gas/go-web-framework-benchmark/master/concurrency_alloc.png" alt="go-gas-concurrency-alloc" />

#### Concurrency-latency

<img src="https://raw.githubusercontent.com/go-gas/go-web-framework-benchmark/master/concurrency_latency.png" alt="go-gas-concurrency-latency" />

#### Concurrency-pipeline

<img src="https://raw.githubusercontent.com/go-gas/go-web-framework-benchmark/master/concurrency-pipeline.png" alt="go-gas-concurrency-pipeline" />

### Roadmap
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
- [ ] Cache
 - [ ] Memory
 - [ ] File
 - [ ] Redis
 - [ ] Memcache
- [ ] i18n
- [ ] HTTPS
- [ ] Command line tools
- [ ] Form handler (maybe next version)
- [ ] Security check features(csrf, xss filter...etc)
