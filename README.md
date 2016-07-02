# Gas

<img src="https://raw.githubusercontent.com/go-gas/gas/master/logo.jpg" alt="go-gas" width="200px" />

[![Build Status](https://api.travis-ci.org/go-gas/gas.svg?branch=master)](https://api.travis-ci.org/go-gas/gas.svg) [![codecov](https://codecov.io/gh/go-gas/gas/branch/master/graph/badge.svg)](https://codecov.io/gh/go-gas/gas) [![Go Report Card](https://goreportcard.com/badge/github.com/go-gas/gas)](https://goreportcard.com/report/github.com/go-gas/gas)

Gas is a Web framework written in Go. And this is not a total complete project.

I just did a minimum workable architechture.

The workable feature is:

- Router (based on [fasthttprouter](https://github.com/buaazp/fasthttprouter) package)
- Context (It can easy to render view and print json string)
- Middleware
- Logger middleware
- Log package
- Read config from a yaml file

And Model is not complete yet. Just finished MySQL sql Builder

##### all feature you can see in Example Directory.

# Install
```
$ go get github.com/go-gas/gas
```

# Run demo
```
$ cd $GOPATH/src/github.com/go-gas/gas/Example
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
```go
g.Router.Use(middleware.LogMiddleware)
```

And you can write your own middleware function

```go
func LogMiddleware(next gas.CHandler) gas.CHandler {
    return func (c *gas.Context) error  {

       // do something before next handler

       err := next(c)

       // do something after next handler

       return err
    }
}
```

### And done

```go
g.Run()
```

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

## Benchmark conclusion

[Iris](https://github.com/kataras/iris) is still the fastest web framework.

But gas is very new, so in the future

I wish this framework might not be a fastest but it will very fast and full featured.

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
