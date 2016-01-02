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
	*Route

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
		Route:            new(Route),
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
		fn := this.matchRoute(ctx)

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

// parseRequest parse all need arugments for match route
func parseRequest(r *http.Request) (paths []string, method, ext, domain string, querys, headers map[string][]string) {

	// method
	method = r.Method

	// domain, actually host
	domain = r.Host

	// querys
	querys = r.URL.Query()

	// headers
	headers = r.Header

	return
}

// isRouteMatch check the request is match a route in global route-function map
func isRouteMatch(rt *Route, ctx *Ctx) bool {
	r := ctx.R

	paths, ext := getPathsAndExt(r.URL)

	args = make(map[string]string)

	// check paths
	if len(rt.paths) != len(paths) {
		return
	}
	for i := range rt.paths {
		y, key, value := isSingleMatch(rt.paths[i], paths[i])
		// validate failed
		if !y {
			return
		}
		// success once
		if key != "" {
			args[key] = value
		}
	}

	// check method
	if len(rt.methods) > 0 {
		if !isSliceMatch(rt.methods, method, args) {
			return
		}
	}

	// check extension
	if len(rt.exts) > 0 {
		if !isSliceMatch(rt.exts, ext, args) {
			return
		}
	}

	// check domain
	if len(rt.domains) > 0 {
		if !isSliceMatch(rt.domains, domain, args) {
			return
		}
	}

	// check querys
	if len(rt.querys) > 0 {
		if !isMapMatch(rt.querys, querys, args) {
			return
		}
	}

	// check headers
	if len(rt.headers) > 0 {
		if !isMapMatch(rt.headers, headers, args) {
			return
		}
	}

	yes = true
	return
}

// isSingleMatch use "==" or regexp to validate a single argument of request is match or not
func isSingleMatch(rtArg, reqArg string, args map[string]string) bool {
	if !strings.HasPrefix(rtArg, "{") || !strings.HasSuffix(rtArg, "}") {
		// common validate
		return rtArg == reqArg
	}

	// use regexp to validate
	rtArg = rtArg[1 : len(rtArg)-1]

	// check the ":" is not at the first or last position
	if index := strings.Index(rtArg, ":"); index > 0 && index < len(rtArg)-1 {
		// get regexp
		regStr := rtArg[index+1:]
		reg := regexp.MustCompile(regStr)

		// regexp validate success
		if reg.MatchString(reqArg) {
			yes = true
			key = rtArg[:index]
			value = reqArg
			return
		}
		// regexp validate failed
		return

		// don't contain ":", means it doesn't need to save in Args
	} else if index == -1 {
		reg := regexp.MustCompile(rtArg)
		yes = reg.MatchString(reqArg)
		return
	}

}

// isSliceMatch check if one item in the route is the request argument
func isSliceMatch(slice []string, single string, args map[string]string) bool {
	for i := range slice {
		if y, key, value := isSingleMatch(slice[i], single); y {
			// success, has one item match
			if key != "" {
				args[key] = value
			}
			return true
		}

	}
	// failed
	return false
}

// isMapMatch check if all map key of route are exists in request map , and at most one item of value(slice) is match the route
func isMapMatch(rtMap map[string]string, reqMap map[string][]string, args map[string]string) bool {
	// use to count the success counts
	flag := 0
roop:
	for k := range rtMap {
		slice, ok := reqMap[k]
		if !ok {
			return false
		}
		for i := range slice {
			if y, key, value := isSingleMatch(rtMap[k], slice[i]); y {
				// success, has one item match
				if key != "" {
					args[key] = value
				}
				flag++
				continue roop
			}
		}
	}

	// success or not
	return flag == len(rtMap)
}

func getPathsAndExt(u *url.URL) (paths []string, ext string) {
	paths0 := strings.Split(path.Clean(u.Path), "/")

	// remove empty
	paths = make([]string, 0, len(paths0))
	for i := range paths0 {
		if rawPaths[i] != "" {
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
