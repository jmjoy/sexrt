package sexrt_test

import (
	"fmt"
	"github.com/jmjoy/sexrt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {

	sexrt.NF(func(ctx *sexrt.Ctx) {
		fmt.Fprintf(ctx.W, "fuck!!!<br />\n")
		fmt.Fprintf(ctx.W, "%#v", ctx.Req)
	})

	sexrt.P("index").Get().Q("a", `{^\d$}`).F(func(ctx *sexrt.Ctx) {
		fmt.Fprintf(ctx.W, "%#v\n", ctx)
		fmt.Fprintf(ctx.W, "%#v\n", ctx.Args)
	})

	sexrt.P("accept").Get().H("Accept", `{accept:html}`).F(func(ctx *sexrt.Ctx) {
		fmt.Fprintf(ctx.W, "%#v\n", ctx)
		fmt.Fprintf(ctx.W, "%#v\n", ctx.Args)
	})

	sexrt.Use()

	fmt.Println("rtFuncMap. Len: ", len(sexrt.RouteFuncMap()), "\n")

	rm := sexrt.RouteFuncMap()
	for i := range rm {
		fmt.Println(i)
	}

	http.ListenAndServe(":8080", nil)
}
