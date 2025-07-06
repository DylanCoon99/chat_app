package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {

	var buf bytes.Buffer

	tracer := New(&buf)

	if tracer == nil {
		t.Error("Failed to create new tracer")
	} else {
		tracer.Trace("Hello from trace package!")
		if buf.String() != "Hello from trace package!\n" {
			t.Errorf("Trace should not write %s", buf.String())
		}
	}

}



func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("something")
}