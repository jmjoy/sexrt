package sexrt

// Route is a rule for matching a request
type Route struct {
	mux *Mux

	paths   []string          // request.URL splits by "/" (e.g. "/hello/world" => ["hello", "world"])
	methods []string          // request Method (e.g. "GET", "POST", "PUT", "DELETE")
	exts    []string          // url extension (e.g. "html", "jpg", "pdf")
	hosts   []string          // the Host in request header
	querys  map[string]string // url querys pair (e.g. "?a=1" => [a: 1])
	headers map[string]string // request header pair (e.g. "Accept: XXX" => [Accept: XXX])
}

// Path add some url segment to a building route, the order is important
func (r *Route) Path(s ...string) *Route {
	r.paths = append(r.paths, s...)
	return r
}

// Method add some reqeust method to a building route
func (r *Route) Method(s ...string) *Route {
	r.methods = append(r.methods, s...)
	return r
}

// Get is same as route.Method("GET")
func (r *Route) Get() *Route {
	return r.Method("GET")
}

// Post is same as route.Method("POST")
func (r *Route) Post() *Route {
	return r.Method("POST")
}

// Put is same as route.Method("PUT")
func (r *Route) Put() *Route {
	return r.Method("PUT")
}

// Delete is same as route.Method("DELETE")
func (r *Route) Delete() *Route {
	return r.Method("DELETE")
}

// Ext add some url extension to a building route
func (r *Route) Ext(s ...string) *Route {
	r.exts = append(r.exts, s...)
	return r
}

// Query add some url querys pair to a building route
func (r *Route) Query(s ...string) *Route {
	if r.querys == nil {
		r.querys = make(map[string]string)
	}

	for i := 0; i < len(s); i += 2 {
		r.querys[s[i]] = s[i+1]
	}

	return r
}

// Host add some host name to a building route
func (r *Route) Host(s ...string) *Route {
	r.hosts = append(r.hosts, s...)
	return r
}

// Header add some request header pair to a building route
func (r *Route) Header(s ...string) *Route {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}

	for i := 0; i < len(s); i += 2 {
		r.headers[s[i]] = s[i+1]
	}

	return r
}

// Func will always deep clone the route and registe it into relative Mux
func (r *Route) Func(fn routeHandler) {
	newRoute := r.clone()
	r.mux.routeHandlerPool[newRoute] = fn
}

func (r *Route) clone() *Route {
	return &Route{
		mux:     r.mux,
		paths:   cloneStringSlice(r.paths),
		methods: cloneStringSlice(r.methods),
		exts:    cloneStringSlice(r.exts),
		hosts:   cloneStringSlice(r.hosts),
		querys:  cloneStringMap(r.querys),
		headers: cloneStringMap(r.headers),
	}
}

func cloneStringSlice(slice []string) []string {
	newSlice := make([]string, 0, len(slice))
	copy(newSlice, slice)
	return newSlice
}

func cloneStringMap(m map[string]string) map[string]string {
	newM := make(map[string]string, len(m))
	for k, v := range m {
		newM[k] = v
	}
	return newM
}
