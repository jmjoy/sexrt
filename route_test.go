package sexrt

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

var testStr0 = []string{
	"hello", `{^\d+$}`, `{num:^\d+$}`,
}

var testParsed0 = []interface{}{
	"hello",
	regexp.MustCompile(`^\d+$`),
	&namedRegexp{Name: "num", Regexp: regexp.MustCompile(`^\d+$`)},
}

var testRoute0 = &Route{
	paths:   []interface{}{"hello", "world", regexp.MustCompile(`^\d+$`)},
	methods: []interface{}{regexp.MustCompile(`^\d+$`), "GET", "POST", "PUT", "DELETE"},
	exts:    []interface{}{"html", "pdf", regexp.MustCompile(`^\d+$`)},
	hosts:   []interface{}{"baidu", regexp.MustCompile(`^\d+$`), "sina"},
	querys: map[string][]interface{}{
		"key":  []interface{}{"value"},
		"key0": []interface{}{"value0", regexp.MustCompile(`^\d+$`)},
		"key1": []interface{}{"value1", &namedRegexp{Name: "name0", Regexp: regexp.MustCompile(`^\d+$`)}},
	},
	headers: map[string][]interface{}{
		"key":  []interface{}{"value"},
		"key0": []interface{}{"value0", regexp.MustCompile(`^\d+$`)},
		"key1": []interface{}{"value1", &namedRegexp{Name: "name0", Regexp: regexp.MustCompile(`^\d+$`)}},
	},
}

func TestRouteClone(t *testing.T) {
	rt := testRoute0.clone()
	t.Logf("%p, %p", testRoute0, rt)
	if rt == testRoute0 {
		t.Fatal("not a new object")
	}
	if checkPointerEqual(rt.paths, testRoute0.paths) ||
		checkPointerEqual(rt.methods, testRoute0.methods) ||
		checkPointerEqual(rt.exts, testRoute0.exts) ||
		checkPointerEqual(rt.hosts, testRoute0.hosts) ||
		checkPointerEqual(rt.querys, testRoute0.querys) ||
		checkPointerEqual(rt.headers, testRoute0.headers) {
		t.Fatal("not a new object")
	}
	if !reflect.DeepEqual(rt, testRoute0) {
		t.Fatal("not deep equal")
	}
}

func TestRouteChainMethods(t *testing.T) {
	rt := new(Route)
	rt.Path("hello", "world").Path(`{^\d+$}`).
		Method(`{^\d+$}`).Get().Post().Put().Delete().
		Ext("html", "pdf").Ext(`{^\d+$}`).
		Host("baidu", `{^\d+$}`).Host("sina").
		Query("key", "value", "key0", "value0", "key0", `{^\d+$}`).
		Query("key1", "value1", "key1", `{name0:^\d+$}`).
		Header("key", "value", "key0", "value0", "key0", `{^\d+$}`).
		Header("key1", "value1", "key1", `{name0:^\d+$}`)

	if !reflect.DeepEqual(rt, testRoute0) {
		t.Fatal("not deep equal")
	}
}

func TestRouteFunc(t *testing.T) {
	mux := NewMux()
	rt := mux.NewRoute()
	fn := func(c *Ctx) error {
		return nil
	}
	rt.Func(fn)

	if len(mux.routeHandlerPool) != 1 {
		t.Fatal("len of routeHandlerPool isn't 1")
	}

	for testRt, testFn := range mux.routeHandlerPool {
		t.Logf("%p, %p", rt, testRt)
		if rt == testRt {
			t.Fatal("not a new object")
		}

		t.Logf("%#v, %#v", rt, testRt)
		if !reflect.DeepEqual(rt, testRt) {
			t.Fatal("not deep equal")
		}

		if !checkPointerEqual(fn, testFn) {
			t.Fatal("not equal")
		}

		return
	}
}

func TestCloneRouteSingle(t *testing.T) {
	testStr := "hello"
	str := cloneRouteSingle(testStr)
	t.Logf("%p, %p", &testStr, &str)
	if checkPointerEqual(&str, &testStr) {
		t.Fatal("not a new object")
	}
	t.Log(str)
	if str != testStr {
		t.Fatal("not equal")
	}

	testReg := regexp.MustCompile(`^\d+$`)
	reg := cloneRouteSingle(testReg)
	t.Logf("%p, %p", &testReg, &reg)
	if testReg == reg {
		t.Fatal("not a new object")
	}
	if !reflect.DeepEqual(testReg, reg) {
		t.Fatal("not deep equal")
	}

	testNamedReg := &namedRegexp{
		Name:   "hello",
		Regexp: regexp.MustCompile(`^\d+$`),
	}
	namedReg := cloneRouteSingle(testNamedReg)
	t.Logf("%p, %p", &testNamedReg, &namedReg)
	if testNamedReg == namedReg {
		t.Fatal("not a new object")
	}
	if !reflect.DeepEqual(testNamedReg, namedReg) {
		t.Fatal("not deep equal")
	}
}

func TestCloneRouteSlice(t *testing.T) {
	parsed := cloneRouteSlice(testParsed0)
	t.Logf("%p, %p", testParsed0, parsed)
	if checkPointerEqual(parsed, testParsed0) {
		t.Fatal("not a new object")
	}
	if !reflect.DeepEqual(testParsed0, parsed) {
		t.Fatal("not deep equal")
	}
}

func TestCloneRouteMap(t *testing.T) {
	testM := map[string][]interface{}{
		"key0": []interface{}{"value0"},
		"key1": []interface{}{"value10", "value11"},
	}
	m := cloneRouteMap(testM)
	t.Logf("%p, %p", testM, m)
	if checkPointerEqual(testM, m) {
		t.Fatal("not a new object")
	}
	if !reflect.DeepEqual(testM, m) {
		t.Fatal("not deep equal")
	}
}

func TestParseAppendString(t *testing.T) {
	result := parseAppendString(testStr0...)
	t.Logf("%+v, %+v", testParsed0, result)
	if !reflect.DeepEqual(result, testParsed0) {
		t.Fatal("not equal")
	}
}

func checkPointerEqual(p, p0 interface{}) bool {
	return fmt.Sprintf("%p", p0) == fmt.Sprintf("%p", p)
}
