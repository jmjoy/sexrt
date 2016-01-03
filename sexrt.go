package sexrt

import (
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
)

type routeHandler func(*Ctx) error

// Ctx is a Context contains http request, response and regexp arguments.
type Ctx struct {
	R    *http.Request
	W    http.ResponseWriter
	Args map[string]string // regexp arguments
}

// Mux is a http.Handler implementer
type Mux struct {
	*http.ServeMux

	routeHandlerPool map[*Route]routeHandler

	notFoundHandler routeHandler
	errorHandler    func(error)
}

// NewMuxWithHandler
func NewMuxWithHandler(notFoundHandler routeHandler, errorHandler func(error)) *Mux {
	if notFoundHandler == nil {
		// default Not Found handler
		notFoundHandler = func(ctx *Ctx) error {
			http.NotFound(ctx.W, ctx.R)
			return nil
		}
	}

	if errorHandler == nil {
		// default error handler
		errorHandler = func(err error) {
			panic(err)
		}
	}

	mux := &Mux{
		ServeMux:         http.NewServeMux(),
		routeHandlerPool: make(map[*Route]routeHandler),
		notFoundHandler:  notFoundHandler,
		errorHandler:     errorHandler,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := &Ctx{
			R:    r,
			W:    w,
			Args: make(map[string]string),
		}

		// get handler and regexp args of a matchesd route
		fn := mux.matchRoute(ctx)

		if err := fn(ctx); err != nil {
			mux.errorHandler(err)
		}
	})

	return mux
}

func NewMux() *Mux {
	return NewMuxWithHandler(nil, nil)
}

func (mux *Mux) NewRoute() *Route {
	return &Route{mux: mux}
}

func (mux *Mux) HandleNotFound(notFoundHandler routeHandler) {
	mux.notFoundHandler = notFoundHandler
}

func (mux *Mux) HandleError(errorHandler func(error)) {
	mux.errorHandler = errorHandler
}

// matchRoute find a route which match the request
func (mux *Mux) matchRoute(ctx *Ctx) routeHandler {
	// find a matched route
	for rt := range mux.routeHandlerPool {
		if is := isRouteMatch(rt, ctx); is {
			return mux.routeHandlerPool[rt]
		}
	}

	// not found
	return mux.notFoundHandler
}

// isRouteMatch check the request is match a route in global route-function map
func isRouteMatch(rt *Route, ctx *Ctx) (is bool) {
	r := ctx.R
	args := ctx.Args

	// check method
	if len(rt.methods) > 0 {
		if !isSliceMatch(rt.methods, r.Method, args) {
			return
		}
	}

	// check host
	if len(rt.hosts) > 0 {
		if !isSliceMatch(rt.hosts, r.Host, args) {
			return
		}
	}

	// parse paths and ext
	paths, ext := getPathsAndExt(r.URL)

	// check paths
	if len(rt.paths) != len(paths) {
		return
	}
	for i := range rt.paths {
		if !isSingleMatch(rt.paths[i], paths[i], args) {
			return
		}
	}

	// check extension
	if len(rt.exts) > 0 && len(paths) > 0 {
		if !isSliceMatch(rt.exts, ext, args) {
			return
		}
	}

	// check querys
	if len(rt.querys) > 0 {
		if !isMapMatch(rt.querys, r.URL.Query(), args) {
			return
		}
	}

	// check headers
	if len(rt.headers) > 0 {
		if !isMapMatch(rt.headers, r.Header, args) {
			return
		}
	}

	return true
}

// isSingleMatch use "==" or regexp to validate a single argument of request is match or not
func isSingleMatch(item interface{}, single string, args map[string]string) bool {
	switch item.(type) {
	case string:
		return single == item.(string)

	case *regexp.Regexp:
		return item.(*regexp.Regexp).MatchString(single)

	case *namedRegexp:
		nr := item.(*namedRegexp)

		if !nr.MatchString(single) {
			return false
		}
		args[nr.Name] = single
		return true

	default:
		panic("Unknow type of slice item")
	}
}

// isSliceMatch check if one item in the route is the request argument
func isSliceMatch(rtSlice []interface{}, single string, args map[string]string) bool {
	for i := range rtSlice {
		if isSingleMatch(rtSlice[i], single, args) {
			return true
		}
	}
	return false
}

// isMapMatch check if all map key of route are exists in request map, and at most one item of value(slice) is match the route
func isMapMatch(rtMap map[string][]interface{}, reqMap map[string][]string, args map[string]string) bool {
loop:
	for k := range rtMap {
		slice, ok := reqMap[k]
		if !ok {
			return false
		}

		for i := range slice {
			if isSliceMatch(rtMap[k], slice[i], args) {
				continue loop
			}
		}

		return false
	}

	return true
}

func getPathsAndExt(u *url.URL) (paths []string, ext string) {
	paths0 := strings.Split(path.Clean(u.Path), "/")

	// remove empty
	paths = make([]string, 0, len(paths0))
	for i := range paths0 {
		if paths0[i] != "" {
			paths = append(paths, paths0[i])
		}
	}

	// split basename and extension
	if length := len(paths); length >= 1 {
		// the "." can't be the first or last character in the last paths item
		lastPath := paths[length-1]
		index := strings.LastIndex(lastPath, ".")
		if index > 0 && index < len(lastPath)-1 {
			ext = lastPath[index+1:]
			paths[length-1] = lastPath[:index]
		}
	}

	return
}
