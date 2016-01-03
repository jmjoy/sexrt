Sexrt [![GoDoc](https://godoc.org/github.com/jmjoy/sexrt?status.png)](https://godoc.org/github.com/jmjoy/sexrt)
====

Sexrt is a sexy HTTP router for golang, inspired by gorilla/mux:

https://github.com/gorilla/mux

## Get started

To install:

    go get github.com/jmjoy/sexrt

Take a try:

```go
package main

import (
	"io"
	"net/http"

	"github.com/jmjoy/sexrt"
)

func main() {
	mux := sexrt.NewMux()
	mux.NewRoute().Func(func(ctx *sexrt.Ctx) error {
		_, err := io.WriteString(ctx.W, "Hello, world!")
		return err
	})
	http.ListenAndServe(":8080", mux)
}
```

## Usage

```go
mux := sexrt.NewMux()
rt := mux.NewRoute()
fn := func(ctx *Ctx) error {
    _, err := io.WriteString(ctx.W, "hello:"+ctx.Args["name"])
    return err
}
rt.Path("name", `{name:\w+}`).Func(fn)
// or: rt.Path("name").Path(`{name:\w+}`).Func(fn)
```

Visit: http://localhost:8080/user/foo

Will return: hello:foo

The argument surround with `{}` means it use regexp, like `{\d+}` matches some numbers,
if you want to save the matched string, use `{<name>:<regexp>}`, than the matched string will be stored in ctx.Args

You can also visit it by: http://localhost:8080/user/foo.html or http://localhost:8080/user/foo.txt and so on...

if you don't like the extension, you can do it simply:

```go
rt.Path("name", `{name:\w+}`).Ext("").Func(fn)
```

## More example

```go
// all suport chain method
rt.Path("foo", "bar").
    Method("GET", "POST").
    Ext("html", "txt").
    Host("localhost", "localhost:8080").
    Query("id", `{id:^\d+$}`, "name", `{name:^\w+$}`).
    Header("Accept", `{html}`, "Accept", `{\*/\*}`).
    Func(fn)

// RESTful
rt.Get()    // same as sexrt.Method("GET")
rt.Post()   // same as sexrt.Method("POST")
rt.Put()    // same as sexrt.Method("PUT")
rt.Delete() // same as sexrt.Method("DELETE")
```
