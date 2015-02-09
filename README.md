#Sexrt

Sexrt is a sexy HTTP router for golang, the design inspiration is origin to gorilla/mux:

http://www.gorillatoolkit.org/pkg/mux

##Get started

To install:

        go get github.com/jmjoy/sexrt

Take a try:

```go
package main

import (
  "fmt"
  "github.com/jmjoy/sexrt"
  "net/http"
)

func main() {
  sexrt.Path("index").Get().Func(func(ctx *sexrt.Ctx) {
    fmt.Fprintln(ctx.W, "Hello world")
  })

  sexrt.Use()

  http.ListenAndServe(":8080", nil)
}
```

Now, you can visit it: http://localhost:8080/index

not only, you also can via http://localhost:8080/index.html or
http://localhost:8080/index.pdf or ... to visit it

it you don't like the extension, you can do it simply:

```go
sexrt.Path("index").Get().Ext("").Func(func(ctx *sexrt.Ctx) {
  fmt.Fprintln(ctx.W, "Hello world")
})
```

or add Ext() some arguments if you want to specify the extension.

if you love, you can use the "shorthand style", just like this:

```go
sexrt.P("index").E("").F(func(ctx *sexrt.Ctx) {
  fmt.Fprintln(ctx.W, "Hello world")
})
```

if you want more url segments, you can do:

    sexrt.P("home", "user", "info").F(whatFunc)

or use "chain style":

    sexrt.P("home").P("user").P("info").F(whatFunc)

many times, you would need to use a regexp to match a url:

    sexrt.P("home").P(`{^\d+$}`).F(whatFunc)

the argument surround with `{}` means it use regexp, if you want to save the matched string:

```go
sexrt.P("home").P(`{id:^\d+$}`).F(func(ctx *sexrt.Ctx) {
  fmt.Fprintf(ctx.W, "%#v", ctx.Args)
})
```

than you can get the string via "sexrt.Ctx.Args"

##More example

```go
whatFunc := func(ctx *sexrt.Ctx) {}

// all suport Validate Functions
sexrt.Path("edit").Method("POST").Ext("html").
    Domain("localhost").Query("id", `{id:^\d+$}`).
    Header("Accept", `{html}`).Func(whatFunc)

// all suport Validate Functions' shorthand
sexrt.P("edit").M("POST").E("html").
    D("localhost").Q("id", `{id:^\d+$}`).
    H("Accept", `{html}`).F(whatFunc)

// HomePage
sexrt.F(whatFunc)

// sub router
article := sexrt.Get().P("article")
article.P("add").F(whatFunc)
article.P("del").F(whatFunc)

// Restful route
sexrt.Get() // same as sexrt.Method("GET")
sexrt.Post() // same as sexrt.Method("POST")
sexrt.Put() // same as sexrt.Method("PUT")
sexrt.Delete() // same as sexrt.Method("DELETE")

// not found handle function
nfFunc := func(ctx *sexrt.Ctx) {}
sexrt.NotFunc(nfFunc)
sexrt.NF(nfFunc)
```

##Reference

example:  http://xxx

document: http://xxx
