package sexrt

import (
	"net/http"
	"path"
	"regexp"
	"strings"
)

type Ctx struct {
	Req  *http.Request
	W    http.ResponseWriter
	Args map[string]string
}

func Use() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// get hanleFunc and regexp args of a matchesd route
		f, args := matchRoute(r)

		ctx := &Ctx{r, w, args}
		f(ctx)
	})

}

func matchRoute(r *http.Request) (function func(*Ctx), args map[string]string) {
	// parse request arguments
	paths, method, ext, domain, querys, headers := parseRequest(r)

	// find a matched route
	for rt := range routeFuncMap {
		if is, arguments := isMatch(rt, paths, method, ext, domain, querys, headers); is {
			function = routeFuncMap[rt]
			args = arguments
			return
		}
	}

	// not found
	function = notFound
	return
}

func parseRequest(r *http.Request) (paths []string, method, ext, domain string, querys, headers map[string][]string) {
	// get a url.path slice
	rawPaths := strings.Split(path.Clean(r.URL.Path), "/")

	// remove empty single
	paths = make([]string, 0, 2)
	for i := range rawPaths {
		if rawPaths[i] != "" {
			paths = append(paths, rawPaths[i])
		}
	}

	// split basename and extension
	if length := len(paths); length >= 1 {
		// the "." can't be the first or last character in the last paths single
		if index := strings.LastIndex(paths[length-1], "."); index > 0 && index < len(paths[length-1])-1 {
			ext = paths[length-1][index+1:]
			paths[length-1] = paths[length-1][:index]
		}
	}

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

func isMatch(rt *route, paths []string, method, ext, domain string, querys, headers map[string][]string) (yes bool, args map[string]string) {
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
		for i := range rt.methods {
			y, key, value := isSingleMatch(rt.methods[i], method)
			if !y {
				return
			}
			// success once
			if key != "" {
				args[key] = value
			}
		}
	}

	// check extension
	if len(rt.exts) > 0 {
		for i := range rt.exts {
			y, key, value := isSingleMatch(rt.exts[i], ext)
			if !y {
				return
			}
			// success once
			if key != "" {
				args[key] = value
			}
		}
	}

	// check domain
	if len(rt.domains) > 0 {
		for i := range rt.domains {
			y, key, value := isSingleMatch(rt.domains[i], domain)
			if !y {
				return
			}
			// success once
			if key != "" {
				args[key] = value
			}
		}
	}

	// check querys
	if len(rt.querys) > 0 {
		for k := range rt.querys {
			arr, ok := querys[k]
			if !ok {
				return
			}
			for i := range arr {
				y, key, value := isSingleMatch(arr[i], rt.querys[k])
				if !y {
					return
				}
				// success once
				if key != "" {
					args[key] = value
				}
			}
		}
	}

	// check headers
	if len(rt.headers) > 0 {
		for k := range rt.headers {
			arr, ok := headers[k]
			if !ok {
				return
			}
			for i := range arr {
				y, key, value := isSingleMatch(arr[i], rt.headers[k])
				if !y {
					return
				}
				// success once
				if key != "" {
					args[key] = value
				}
			}
		}
	}

	yes = true
	return
}

func isSingleMatch(rtArg, reqArg string) (yes bool, key, value string) {
	// use regexp to validate
	if strings.HasPrefix(rtArg, "{") && strings.HasSuffix(rtArg, "}") {
		rtArg = strings.TrimLeft(rtArg, "{")
		rtArg = strings.TrimRight(rtArg, "}")

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

	// common validate
	yes = rtArg == reqArg
	return
}
