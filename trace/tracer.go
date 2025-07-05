package trace


import (

)


// Trace is an object capable of tracing events in code
type Trace interface {
	Trace(...interface{})
}