package sexrt

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	gourl "net/url"
	"strings"
	"testing"
)

var (
	testContent = "Hello world!"

	errHehe = errors.New("hehe")

	testNotFoundHandler = func(ctx *Ctx) error {
		io.WriteString(ctx.W, "hehele")
		return errHehe
	}

	testErrorHandler = func(err error) {}

	testHandler = func(ctx *Ctx) error {
		io.WriteString(ctx.W, testContent)
		return nil
	}
)

func TestMuxNew(t *testing.T) {
	mux := NewMuxWithHandler(testNotFoundHandler, testErrorHandler)
	if !checkPointerEqual(mux.notFoundHandler, testNotFoundHandler) ||
		!checkPointerEqual(mux.errorHandler, testErrorHandler) {
		t.Fatal("handler not equal")
	}

	mux = NewMux()
	srv := httptest.NewServer(mux)
	defer srv.Close()

	t.Log("srv URL:", srv.URL)
	testHTTPResponse("GET", srv.URL, "", func(body string, resp *http.Response) {
		if resp.StatusCode != 404 {
			t.Fatal("default not found handler not correct!")
		}
	})

	mux.HandleNotFound(testNotFoundHandler)
	mux.HandleError(func(err error) {
		if err != errHehe {
			t.Fatal("user defined error handler not correct!")
		}
	})

	testHTTPResponse("GET", srv.URL, "", func(body string, resp *http.Response) {
		if body != "hehele" {
			t.Fatal("use defined not found handler not correct!")
		}
	})
}

func TestMuxRoutePath(t *testing.T) {
	mux := NewMux()
	rt := mux.NewRoute()

	srv := httptest.NewServer(mux)
	defer srv.Close()

	// test index page
	rt.Func(testHandler)

	testHTTPResponse("GET", srv.URL, "", func(body string, resp *http.Response) {
		if body != testContent {
			t.Fatal("`/`: body not correct!")
		}
	})

	// test long paths
	rt.Path("hello", `{name:\w+}`, `{\d+}`).Func(func(ctx *Ctx) error {
		_, err := io.WriteString(ctx.W, "hello:"+ctx.Args["name"])
		return err
	})

	us := []string{
		srv.URL + "/hello/jmjoy/123",
		srv.URL + "/hello/jmjoy/123/",
		srv.URL + "/hello/jmjoy/123.html",
		srv.URL + "/hello/jmjoy/123.html/",
		srv.URL + "/hello/jmjoy/123.pdf",
	}
	for _, u := range us {
		testHTTPResponse("GET", u, "", func(body string, resp *http.Response) {
			if body != "hello:jmjoy" {
				t.Fatal(u + ": body not correct!")
			}
		})
		testHTTPResponse("POST", u, "", func(body string, resp *http.Response) {
			if body != "hello:jmjoy" {
				t.Fatal(u + ": body not correct!")
			}
		})
	}

	us = []string{
		srv.URL + "/hello/jmjoy/",
		srv.URL + "/hello/jmjoy",
	}
	for _, u := range us {
		testHTTPResponse("GET", u, "", func(body string, resp *http.Response) {
			if resp.StatusCode != 404 {
				t.Fatal(u + ": can found?!")
			}
		})
	}
}

func TestMuxRouteExt(t *testing.T) {
	mux := NewMux()
	rt := mux.NewRoute()

	srv := httptest.NewServer(mux)
	defer srv.Close()

	rt.Ext("html").Func(testHandler)

	testHTTPResponse("GET", srv.URL, "", func(body string, resp *http.Response) {
		if body != testContent {
			t.Fatal("`/`: body not correct!")
		}
	})

	rt.Ext("pdf").Path("hello").Func(testHandler)
	us := []string{
		srv.URL + "/hello.html",
		srv.URL + "/hello.pdf",
		srv.URL + "/hello.pdf/",
	}
	for _, u := range us {
		testHTTPResponse("GET", u, "", func(body string, resp *http.Response) {
			if body != testContent {
				t.Fatal(u + ": body not correct!")
			}
		})
	}

	u := srv.URL + "/hello.txt"
	testHTTPResponse("GET", u, "", func(body string, resp *http.Response) {
		if resp.StatusCode != 404 {
			t.Fatal(u + ": can found?!")
		}
	})
}

func TestMuxRouteMethod(t *testing.T) {
	mux := NewMux()
	rt := mux.NewRoute()

	srv := httptest.NewServer(mux)
	defer srv.Close()

	rt.Post().Func(testHandler)

	u := srv.URL
	testHTTPResponse("POST", u, "", func(body string, resp *http.Response) {
		if body != testContent {
			t.Fatal(u + ": body not correct!")
		}
	})
	testHTTPResponse("GET", u, "", func(body string, resp *http.Response) {
		if resp.StatusCode != 404 {
			t.Fatal(u + ": can found?!")
		}
	})
}

func TestMuxRouteHost(t *testing.T) {
	mux := NewMux()
	rt := mux.NewRoute()

	srv := httptest.NewServer(mux)
	defer srv.Close()

	rt.Host("www.baidu.com").Func(testHandler)

	u := srv.URL
	testHTTPResponse("GET", u, "", func(body string, resp *http.Response) {
		if resp.StatusCode != 404 {
			t.Fatal(u + ": can found?!")
		}
	})

	url, err := gourl.Parse(u)
	if err != nil {
		panic(err)
	}
	rt.Host(url.Host).Func(testHandler)

	testHTTPResponse("GET", u, "", func(body string, resp *http.Response) {
		if body != testContent {
			t.Fatal(u + ": body not correct!")
		}
	})
}

func TestMuxRouteQuery(t *testing.T) {
	mux := NewMux()
	rt := mux.NewRoute()

	srv := httptest.NewServer(mux)
	defer srv.Close()

	rt.Query("name", `{name:^\w+$}`, "age", `{age:^\d+$}`).Func(func(ctx *Ctx) error {
		io.WriteString(ctx.W, ctx.Args["name"]+":"+ctx.Args["age"])
		return nil
	})

	u := srv.URL + "/?name=jmjoy&age=23"
	testHTTPResponse("GET", u, "", func(body string, resp *http.Response) {
		if body != "jmjoy:23" {
			t.Log("body:", body)
			t.Fatal(u + ": body not correct!")
		}
	})

	us := []string{
		srv.URL,
		srv.URL + "/?name=jmjoy",
		srv.URL + "/?age=23",
		srv.URL + "/?name=jmjoy&age=hello",
	}
	for _, u := range us {
		testHTTPResponse("GET", u, "", func(body string, resp *http.Response) {
			if resp.StatusCode != 404 {
				t.Fatal(u + ": can found?!")
			}
		})
	}
}

func TestMuxRouteHeader(t *testing.T) {
	mux := NewMux()
	rt := mux.NewRoute()

	srv := httptest.NewServer(mux)
	defer srv.Close()

	rt.Header("Accept", `{html}`, "Accept", `{\*/\*}`).Func(testHandler)

	headers := []map[string]string{
		map[string]string{"Accept": "*/*"},
		map[string]string{"Accept": "text/html;text/css"},
	}
	for _, header := range headers {
		testHTTPResponseSetHeader("GET", srv.URL, "", header, func(body string, resp *http.Response) {
			if body != testContent {
				t.Fatal("/: body not correct!")
			}
		})
	}

	testHTTPResponse("GET", srv.URL, "", func(body string, resp *http.Response) {
		if resp.StatusCode != 404 {
			t.Fatal("/: can found?!")
		}
	})
}

func checkPointerEqual(p, p0 interface{}) bool {
	return fmt.Sprintf("%p", p0) == fmt.Sprintf("%p", p)
}

func testHTTPResponse(method, u, body string, fn func(string, *http.Response)) {
	testHTTPResponseSetHeader(method, u, body, nil, fn)
}

func testHTTPResponseSetHeader(method, u, body string, header map[string]string, fn func(string, *http.Response)) {
	req, err := http.NewRequest(method, u, strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	switch method {
	case "POST", "PUT":
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fn(string(buf), resp)
}
