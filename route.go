package sexrt

import (
	"fmt"
	"net/http"
)

type route struct {
	paths   []string
	methods []string
	exts    []string
	querys  map[string]string
	domains []string
	headers map[string]string
}

var routeFuncMap = make(map[*route]func(*Ctx))

func RouteFuncMap() map[*route]func(*Ctx) {
	return routeFuncMap
}

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
//
// ----------------------------------------------------------------------------

var notFound = func(ctx *Ctx) {
	http.NotFound(ctx.W, ctx.Req)
}

func NotFound(function func(ctx *Ctx)) {
	notFound = function
}

func NF(function func(ctx *Ctx)) {
	NotFound(function)
}

// ----------------------------------------------------------------------------
//
// ----------------------------------------------------------------------------

func (this *route) Path(s ...string) *route {
	this.paths = append(this.paths, s...)
	return this
}

func (this *route) Method(s ...string) *route {
	this.methods = append(this.methods, s...)
	return this
}

func (this *route) Get() *route {
	this.methods = append(this.methods, "GET")
	return this
}

func (this *route) Post() *route {
	this.methods = append(this.methods, "POST")
	return this
}

func (this *route) Put() *route {
	this.methods = append(this.methods, "PUT")
	return this
}

func (this *route) Delete() *route {
	this.methods = append(this.methods, "DELETE")
	return this
}

func (this *route) Ext(s ...string) *route {
	this.exts = append(this.exts, s...)
	return this
}

func (this *route) Query(key, value string) *route {
	if this.querys == nil {
		this.querys = make(map[string]string)
	}
	this.querys[key] = value
	return this
}

func (this *route) Domain(s ...string) *route {
	this.domains = append(this.domains, s...)
	return this
}

func (this *route) Header(key, value string) *route {
	if this.headers == nil {
		this.headers = make(map[string]string)
	}
	this.headers[key] = value
	return this
}

func (this *route) Func(function func(*Ctx)) {
	newRoute := *this
	routeFuncMap[&newRoute] = function
}

// ----------------------------------------------------------------------------
//
// ----------------------------------------------------------------------------

func Path(s ...string) *route {
	this := new(route)
	this.paths = append(this.paths, s...)
	return this
}

func Method(s ...string) *route {
	this := new(route)
	this.methods = append(this.methods, s...)
	return this
}

func Get() *route {
	this := new(route)
	this.methods = append(this.methods, "GET")
	return this
}

func Post() *route {
	this := new(route)
	this.methods = append(this.methods, "POST")
	return this
}

func Put() *route {
	this := new(route)
	this.methods = append(this.methods, "PUT")
	return this
}

func Delete() *route {
	this := new(route)
	this.methods = append(this.methods, "DELETE")
	return this
}

func Ext(s ...string) *route {
	this := new(route)
	this.exts = append(this.exts, s...)
	return this
}

func Query(key, value string) *route {
	this := new(route)
	this.querys = make(map[string]string)
	this.querys[key] = value
	return this
}

func Domain(s ...string) *route {
	this := new(route)
	this.domains = append(this.domains, s...)
	return this
}

func Header(key, value string) *route {
	this := new(route)
	this.headers = make(map[string]string)
	this.headers[key] = value
	return this
}

func Func(function func(*Ctx)) {
	this := new(route)
	newRoute := *this
	routeFuncMap[&newRoute] = function
}

// ----------------------------------------------------------------------------
// alias
// ----------------------------------------------------------------------------

func (this *route) P(s ...string) *route {
	return this.Path(s...)
}

func (this *route) M(s ...string) *route {
	return this.Method(s...)
}

func (this *route) E(s ...string) *route {
	return this.Ext(s...)
}

func (this *route) Q(key, value string) *route {
	return this.Query(key, value)
}

func (this *route) D(s ...string) *route {
	return this.Domain(s...)
}

func (this *route) H(key, value string) *route {
	return this.Header(key, value)
}

func (this *route) F(function func(*Ctx)) {
	this.Func(function)
}

func P(s ...string) *route {
	return Path(s...)
}

func M(s ...string) *route {
	return Method(s...)
}

func E(s ...string) *route {
	return Ext(s...)
}

func Q(key, value string) *route {
	return Query(key, value)
}

func D(s ...string) *route {
	return Domain(s...)
}

func H(key, value string) *route {
	return Header(key, value)
}

func F(function func(*Ctx)) {
	Func(function)
}
