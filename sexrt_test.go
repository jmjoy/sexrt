package sexrt_test

import (
	"fmt"
	"github.com/jmjoy/sexrt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {

	sexrt.NF(func(ctx *sexrt.Ctx) {
		fmt.Fprintf(ctx.W, "fuck!!!")
	})

	// there is a bug in regexp
	sexrt.E("").P("index").Get().Q("a", `{^\d$}`).F(func(ctx *sexrt.Ctx) {
		fmt.Fprintf(ctx.W, "%#v", ctx)
	})

	sexrt.P("think").H("Accept", `{html}`).F(func(ctx *sexrt.Ctx) {
		fmt.Fprintf(ctx.W, "%#v", ctx)
	})

	sexrt.Use()

	rm := sexrt.RouteFuncMap()
	for i := range rm {
		fmt.Println(i)
	}

	http.ListenAndServe(":8080", nil)
}
