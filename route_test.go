package sexrt_test

import (
	"github.com/jmjoy/sexrt"
	"testing"
)

func TestAllFunction(t *testing.T) {
	rt := sexrt.P("one").P("two").P("three", "four", "five").
		M("one").M("two").M("three", "four", "five").
		E("one").E("two").E("three", "four", "five").
		D("one").D("two").D("three", "four", "five").
		Get().Get().Post().Post().Put().Put().Delete().Delete().
		Q("key1", "value1").Q("key2", "value2").
		H("key1", "value1").H("key2", "value2")

	t.Log("\n\n", rt, "\n\n")

}

func TestF(t *testing.T) {
	whatFunc := func(ctx *sexrt.Ctx) {}
	sexrt.F(whatFunc)
	sexrt.P("index").F(whatFunc)

	for i, v := range sexrt.RouteFuncMap() {
		t.Log("\n++++++\n", i, "\n-------\n", v)
	}
}
