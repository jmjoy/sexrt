/*
Sexrt is a sexy HTTP router for golang.
https://github.com/jmjoy/sexrt

Features:

  * Six ways to design your route rules, they are: Path, Method, Ext, Domain, Query, Header.
  * All route rules method has shorthand alias.
  * All route rules method support regexp, although some method do not need.
  * The handle function just has one argument, which wraps http.ResponseWriter, http.Request and Regexp matched string.
  * Support sub router

Usage:

To visit http://localhost:8080

  sexrt.Func(func (ctx *sexrt.Ctx) {
    fmt.Fprint(ctx.W, ctx.Args)
  })
  sexrt.Use()
  http.ListenAndServe(":8080", nil)

To visit http://localhost:8080/home/user

  sexrt.Path("home", "user").Func(whatFunc)

or

  sexrt.Path("home").Path("user").Func(whatFunc)

but below urls also can match this route, such as:

http://localhost:8080/home/user.html
http://localhost:8080/home/user.xml
http://localhost:8080/home/user.pdf

To limit some extension:

  // only "html"
  sexrt.Path("home").Path("user").Ext("html").Func(whatFunc)
  // only "html", "xml"
  sexrt.Path("home").Path("user").Ext("html", "xml").Func(whatFunc)
  // only "html", "xml", "pdf"
  sexrt.Path("home").Path("user").Ext("html").Ext("xml").Ext("pdf").Func(whatFunc)
  // can't have extension
  sexrt.Path("home").Path("user").Ext("").Func(whatFunc)

To visit http://localhost:8080, but only with Request Method POST

  sexrt.Method("POST").Func(whatFunc)

GET and POST all can visit:

  sexrt.Method("GET").Method("POST").Func(whatFunc)

Restful route:

  // same as sexrt.Method("GET").Func(whatFunc)
  sexrt.Get().Func(whatFunc)
  // same as sexrt.Method("POST").Func(whatFunc)
  sexrt.Post().Func(whatFunc)
  // same as sexrt.Method("PUT").Func(whatFunc)
  sexrt.Put().Func(whatFunc)
  // same as sexrt.Method("DELETE").Func(whatFunc)
  sexrt.Delete().Func(whatFunc)

To limit the host:

  sexrt.Domain("www.google.com").Func(whatFunc)
  sexrt.Domain("www.google.com", "www.baidu.com").Func(whatFunc)
  sexrt.Domain("www.google.com").Domain("www.baidu.com").Func(whatFunc)

To limit the url query arugment, to visit http://localhost:8080/index.html?a=1&b=2

  sexrt.Path("index").Query("a", "1").Query("b", "2").Func(whatFunc)

To limit the request header

  sexrt.Header("Accept", "{html}").Func(whatFunc)

To use regexp, the argument surround with {} means it use regexp:

  // http://localhost:8080/article/123
  sexrt.Path("article").Path(`{^\d+$}`).Func(whatFunc)

if you want to save the matched string:

  // http://localhost:8080/article/123
  sexrt.Path("article").Path(`{id:^\d+$}`).Func(func (ctx *sexrt.Ctx) {
    fmt.Fprint(ctx.W, ctx.Args["id"])
  })

Path, Method, Ext, Domain, the second argument of Query and Header support regexp

There is shorthand alias for all method:

  // same as sexrt.Path("index").Func(whatFunc)
  sexrt.P("index").F(whatFunc)

shorthand alias map:

  Path    =>  P
  Method  =>  M
  Ext     =>  E
  Domain  =>  D
  Query   =>  Q
  Header  =>  H
  Func    =>  F

Sexrt support sub router:

  useinfo := sexrt.P("userinfo")
  sexrt.Get().F(aFunc)
  sexrt.Post().F(bFunc)

If not route is matched a reqeust, NotFound handle function is called, the default is http.NotFound

You can change the default by:

  sexrt.NotFound(nfFunc)

or

  sexrt.NF(nfFunc)

*/
package sexrt
