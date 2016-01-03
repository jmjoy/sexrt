package sexrt

import (
	"regexp"
	"strings"
)

type namedRegexp struct {
	Name string
	*regexp.Regexp
}

// Route is a rule for matching a request
type Route struct {
	mux *Mux

	paths   []interface{}            // request.URL splits by "/" (e.g. "/hello/world" => ["hello", "world"])
	methods []interface{}            // request Method (e.g. "GET", "POST", "PUT", "DELETE")
	exts    []interface{}            // url extension (e.g. "html", "jpg", "pdf")
	hosts   []interface{}            // the Host in request header
	querys  map[string][]interface{} // url querys pair (e.g. "?a=1" => [a: 1])
	headers map[string][]interface{} // request header pair (e.g. "Accept: XXX" => [Accept: XXX])
}

// Path add some url segment to a building route, the order is important
func (rt *Route) Path(s ...string) *Route {
	rt.paths = append(rt.paths, parseAppendString(s...)...)
	return rt
}

// Method add some reqeust method to a building route
func (rt *Route) Method(s ...string) *Route {
	rt.methods = append(rt.methods, parseAppendString(s...)...)
	return rt
}

// Get is same as route.Method("GET")
func (rt *Route) Get() *Route {
	return rt.Method("GET")
}

// Post is same as route.Method("POST")
func (rt *Route) Post() *Route {
	return rt.Method("POST")
}

// Put is same as route.Method("PUT")
func (rt *Route) Put() *Route {
	return rt.Method("PUT")
}

// Delete is same as route.Method("DELETE")
func (rt *Route) Delete() *Route {
	return rt.Method("DELETE")
}

// Ext add some url extension to a building route
func (rt *Route) Ext(s ...string) *Route {
	rt.exts = append(rt.exts, parseAppendString(s...)...)
	return rt
}

// Query add some url querys pair to a building route
func (rt *Route) Query(s ...string) *Route {
	if rt.querys == nil {
		rt.querys = make(map[string][]interface{})
	}

	for i := 0; i < len(s); i += 2 {
		rt.querys[s[i]] = append(rt.querys[s[i]], parseAppendString(s[i+1])[0])
	}

	return rt
}

// Host add some host name to a building route
func (rt *Route) Host(s ...string) *Route {
	rt.hosts = append(rt.hosts, parseAppendString(s...)...)
	return rt
}

// Header add some request header pair to a building route
func (rt *Route) Header(s ...string) *Route {
	if rt.headers == nil {
		rt.headers = make(map[string][]interface{})
	}

	for i := 0; i < len(s); i += 2 {
		rt.headers[s[i]] = append(rt.headers[s[i]], parseAppendString(s[i+1])[0])
	}

	return rt
}

// Func will always deep clone the route and registe it into relative Mux
func (rt *Route) Func(fn routeHandler) {
	newRoute := rt.clone()
	rt.mux.routeHandlerPool[newRoute] = fn
}

func (rt *Route) clone() *Route {
	return &Route{
		mux:     rt.mux,
		paths:   cloneRouteSlice(rt.paths),
		methods: cloneRouteSlice(rt.methods),
		exts:    cloneRouteSlice(rt.exts),
		hosts:   cloneRouteSlice(rt.hosts),
		querys:  cloneRouteMap(rt.querys),
		headers: cloneRouteMap(rt.headers),
	}
}

func cloneRouteSingle(item interface{}) (newItem interface{}) {
	switch item.(type) {
	case string:
		newItem = item.(string)

	case *regexp.Regexp:
		reg := *(item.(*regexp.Regexp))
		newItem = &reg

	case *namedRegexp:
		nr := item.(*namedRegexp)
		reg := *nr.Regexp
		newItem = &namedRegexp{
			Name:   nr.Name,
			Regexp: &reg,
		}

	default:
		panic("Unknow type of slice item")
	}

	return
}

func cloneRouteSlice(slice []interface{}) []interface{} {
	if slice == nil {
		return nil
	}

	newSlice := make([]interface{}, 0, len(slice))

	for _, item := range slice {
		newSlice = append(newSlice, cloneRouteSingle(item))
	}

	return newSlice
}

func cloneRouteMap(m map[string][]interface{}) map[string][]interface{} {
	if m == nil {
		return nil
	}

	newM := make(map[string][]interface{}, len(m))
	for k, slice := range m {
		for _, item := range slice {
			newM[k] = append(newM[k], cloneRouteSingle(item))
		}
	}
	return newM
}

func parseAppendString(s ...string) []interface{} {
	newSlice := make([]interface{}, 0, len(s))

	for _, str := range s {
		if !strings.HasPrefix(str, "{") || !strings.HasSuffix(str, "}") {
			// common string, use `==` to validate
			newSlice = append(newSlice, str)
			continue
		}

		// regexp string, validate by regexp
		// remove `{ }`
		str = str[1 : len(str)-1]

		// check the ":" is not at the first or last position
		if index := strings.Index(str, ":"); index > 0 && index < len(str)-1 {
			// named regexp string
			newSlice = append(newSlice, &namedRegexp{
				Name:   str[:index],
				Regexp: regexp.MustCompile(str[index+1:]),
			})
			continue
		}
		// unmamed regexp string
		newSlice = append(newSlice, regexp.MustCompile(str))
	}

	return newSlice
}
