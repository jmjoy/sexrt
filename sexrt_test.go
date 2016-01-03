package sexrt_test

// import (
// 	"fmt"
// 	"github.com/jmjoy/sexrt"
// 	"net/http"
// 	"testing"
// )

// func TestServer(t *testing.T) {
// 	sexrt.NF(func(ctx *sexrt.Ctx) {
// 		fmt.Fprint(ctx.W, "I didn't found!!!")
// 	})

// 	sexrt.P("home").P("user", `{name:^\w{3,5}$}`).
// 		M("GET").M("POST", "PUT").
// 		E("html").E("xml", "pdf").
// 		D(`{domain:^.*$}`).
// 		Q("a", "1").Q("b", `{^b$}`).Q("c", `{c:^\d+$}`).
// 		H("Accept", `{html}`).
// 		F(func(ctx *sexrt.Ctx) {

// 		fmt.Fprint(ctx.W, ctx.Args)
// 	})

// 	for i := range sexrt.RouteFuncMap() {
// 		fmt.Println(i)
// 	}

// 	sexrt.Use()

// 	http.ListenAndServe(":8080", nil)
// }
