package trace


import (
	"io"
	"fmt"
)


// Trace is an object capable of tracing events in code
// So type "Tracer" is any type that implements the Trace method
type Tracer interface {
	Trace(...interface{})
}


type tracer struct {
	out io.Writer
}


func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}


func New(w io.Writer) Tracer {
	return &tracer{out: w}
}



type nilTracer struct {

}

func (t nilTracer) Trace(a...interface{}) {

}


func Off() Tracer {
	return &nilTracer{}
}