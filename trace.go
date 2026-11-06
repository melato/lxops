package lxops

import (
	_ "embed"
)

//go:embed trace.yaml
var traceDescriptions []byte

var TraceDevice bool

type Trace struct {
}

func (t *Trace) Funcs() map[string]func(bool) {
	return map[string]func(bool){
		"device": func(b bool) { TraceDevice = b },
	}
}

func (t *Trace) Descriptions() []byte {
	return traceDescriptions
}
