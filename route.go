package sexrt

import (
	"fmt"
	"net/http"
)

// route is actually a queue for matching a request
type route struct {
	paths   []string          // request.URL splits by "/", like "/hello/world" => ["hello", "world"]
	methods []string          // request Method, like "GET", "POST", "PUT", "DELETE"
	exts    []string          // url extension , like "html", "jpg", "pdf"
	querys  map[string]string // url querys argument pair, like "xxx?a=1" => [a: 1]
	domains []string          // actually the Host in request header
	headers map[string]string // request header pair, link "Accept: XXX" => [Accept: XXX]
}

// routeFuncMap is the global map for route and function mapping
var routeFuncMap = make(map[*route]func(*Ctx))

// RouteFuncMap return the golbal map routeFuncMap
func RouteFuncMap() map[*route]func(*Ctx) {
	return routeFuncMap
}

// String make route implements Stringer interface
func (this *route) String() string {
	var str string
	str += fmt.Sprintf("paths: %v\n", this.paths)
	str += fmt.Sprintf("methods: %v\n", this.methods)
	str += fmt.Sprintf("exts: %v\n", this.exts)
	str += fmt.Sprintf("querys: %v\n", this.querys)
	str += fmt.Sprintf("domains: %v\n", this.domains)
	str += fmt.Sprintf("headers: %v\n", this.headers)
	return str
}

// ----------------------------------------------------------------------------
// not found handle function
// ----------------------------------------------------------------------------

// notFound is the global function for handleing not found
var notFound = func(ctx *Ctx) {
	http.NotFound(ctx.W, ctx.Req)
}

// NotFound change the function trigger while not found, the default is http.NotFound
func NotFound(function func(ctx *Ctx)) {
	notFound = function
}

// NF alias to NotFound
func NF(function func(ctx *Ctx)) {
	NotFound(function)
}

// ----------------------------------------------------------------------------
// struct route method
// ----------------------------------------------------------------------------

// Path add some url segment to a building route, the order is important
func (this *route) Path(s ...string) *route {
	this.paths = append(this.paths, s...)
	return this
}

// Method add some reqeust method to a building route
func (this *route) Method(s ...string) *route {
	this.methods = append(this.methods, s...)
	return this
}

// Get is same as route.Method("GET")
func (this *route) Get() *route {
	this.methods = append(this.methods, "GET")
	return this
}

// Post is same as route.Method("POST")
func (this *route) Post() *route {
	this.methods = append(this.methods, "POST")
	return this
}

// Put is same as route.Method("PUT")
func (this *route) Put() *route {
	this.methods = append(this.methods, "PUT")
	return this
}

// DELETE is same as route.Method("DELETE")
func (this *route) Delete() *route {
	this.methods = append(this.methods, "DELETE")
	return this
}

// Ext add some url extension to a building route
func (this *route) Ext(s ...string) *route {
	this.exts = append(this.exts, s...)
	return this
}

// Query add some url query argument pair to a building route
func (this *route) Query(key, value string) *route {
	if this.querys == nil {
		this.querys = make(map[string]string)
	}
	this.querys[key] = value
	return this
}

// Domain add some host name to a building route
func (this *route) Domain(s ...string) *route {
	this.domains = append(this.domains, s...)
	return this
}

// Header add some request header pair to a building route
func (this *route) Header(key, value string) *route {
	if this.headers == nil {
		this.headers = make(map[string]string)
	}
	this.headers[key] = value
	return this
}

// Func will always copy the route and registe it into global route-function map
func (this *route) Func(function func(*Ctx)) {
	newRoute := *this
	routeFuncMap[&newRoute] = function
}

// ----------------------------------------------------------------------------
// package function
// ----------------------------------------------------------------------------

// Path add some url segment to a new route, the order is important
func Path(s ...string) *route {
	this := new(route)
	this.paths = append(this.paths, s...)
	return this
}

// Method add some reqeust method to a new route
func Method(s ...string) *route {
	this := new(route)
	this.methods = append(this.methods, s...)
	return this
}

// Get is same as Method("GET")
func Get() *route {
	this := new(route)
	this.methods = append(this.methods, "GET")
	return this
}

// Post is same as Method("POST")
func Post() *route {
	this := new(route)
	this.methods = append(this.methods, "POST")
	return this
}

// Put is same as Method("PUT")
func Put() *route {
	this := new(route)
	this.methods = append(this.methods, "PUT")
	return this
}

// Delete is same as Method("DELETE")
func Delete() *route {
	this := new(route)
	this.methods = append(this.methods, "DELETE")
	return this
}

// Ext add some url extension to a new route
func Ext(s ...string) *route {
	this := new(route)
	this.exts = append(this.exts, s...)
	return this
}

// Query add some url query argument pair to a new route
func Query(key, value string) *route {
	this := new(route)
	this.querys = make(map[string]string)
	this.querys[key] = value
	return this
}

// Domain add some host name to a new route
func Domain(s ...string) *route {
	this := new(route)
	this.domains = append(this.domains, s...)
	return this
}

// Header add some request header pair to a new route
func Header(key, value string) *route {
	this := new(route)
	this.headers = make(map[string]string)
	this.headers[key] = value
	return this
}

// Func will build a empty url route "/" and registe it into global route-function map,
// this function may just be called one time
func Func(function func(*Ctx)) {
	this := new(route)
	newRoute := *this
	routeFuncMap[&newRoute] = function
}

// ----------------------------------------------------------------------------
// alias
// ----------------------------------------------------------------------------

// P alias to Path
func (this *route) P(s ...string) *route {
	return this.Path(s...)
}

// M alias to Method
func (this *route) M(s ...string) *route {
	return this.Method(s...)
}

// E alias to Ext
func (this *route) E(s ...string) *route {
	return this.Ext(s...)
}

// Q alias to Query
func (this *route) Q(key, value string) *route {
	return this.Query(key, value)
}

// D alias to Domain
func (this *route) D(s ...string) *route {
	return this.Domain(s...)
}

// H alias to Header
func (this *route) H(key, value string) *route {
	return this.Header(key, value)
}

// F alias Func
func (this *route) F(function func(*Ctx)) {
	this.Func(function)
}

// P alias to Path
func P(s ...string) *route {
	return Path(s...)
}

// M alias to Method
func M(s ...string) *route {
	return Method(s...)
}

// E alias to Ext
func E(s ...string) *route {
	return Ext(s...)
}

// Q alias to Query
func Q(key, value string) *route {
	return Query(key, value)
}

// D alias to Domain
func D(s ...string) *route {
	return Domain(s...)
}

// H alias to Header
func H(key, value string) *route {
	return Header(key, value)
}

// F alias Func
func F(function func(*Ctx)) {
	Func(function)
}
