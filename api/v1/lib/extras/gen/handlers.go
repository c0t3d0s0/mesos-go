// +build ignore

package main

import (
	"os"
	"text/template"
)

func main() {
	Run(handlersTemplate, nil, os.Args...)
}

var handlersTemplate = template.Must(template.New("").Parse(`package {{.Package}}

// go generate {{.Args}}
// GENERATED CODE FOLLOWS; DO NOT EDIT.

import (
	"context"
{{range .Imports}}
	{{ printf "%q" . -}}
{{end}}
)

type (
	// Handler is invoked upon the occurrence of some scheduler event that is generated
	// by some other component in the Mesos ecosystem (e.g. master, agent, executor, etc.)
	Handler interface {
		HandleEvent(context.Context, {{.EventType}}) error
	}

	// HandlerFunc is a functional adaptation of the Handler interface
	HandlerFunc func(context.Context, {{.EventType}}) error

	// Handlers executes an event Handler according to the event's type
	Handlers map[scheduler.Event_Type]Handler

	// HandlerFuncs executes an event HandlerFunc according to the event's type
	HandlerFuncs map[scheduler.Event_Type]HandlerFunc
)

// HandleEvent implements Handler for HandlerFunc
func (f HandlerFunc) HandleEvent(ctx context.Context, e {{.EventType}}) error { return f(ctx, e) }

var noopHandler = func(_ context.Context, _ {{.EventType}}) error { return nil }

// NoopHandler returns a HandlerFunc that does nothing and always returns nil
func NoopHandler() HandlerFunc { return noopHandler }

// HandleEvent implements Handler for Handlers
func (hs Handlers) HandleEvent(ctx context.Context, e {{.EventType}}) (err error) {
	if h := hs[e.GetType()]; h != nil {
		return h.HandleEvent(ctx, e)
	}
	return nil
}

// HandleEvent implements Handler for HandlerFuncs
func (hs HandlerFuncs) HandleEvent(ctx context.Context, e {{.EventType}}) (err error) {
	if h := hs[e.GetType()]; h != nil {
		return h.HandleEvent(ctx, e)
	}
	return nil
}

// Otherwise returns a HandlerFunc that attempts to process an event with the Handlers map; unmatched event types are
// processed by the given HandlerFunc. A nil HandlerFunc parameter is effecitvely a noop.
func (hs Handlers) Otherwise(f HandlerFunc) HandlerFunc {
	if f == nil {
		return hs.HandleEvent
	}
	return func(ctx context.Context, e {{.EventType}}) error {
		if h := hs[e.GetType()]; h != nil {
			return h.HandleEvent(ctx, e)
		}
		return f(ctx, e)
	}
}

// Otherwise returns a HandlerFunc that attempts to process an event with the HandlerFuncs map; unmatched event types
// are processed by the given HandlerFunc. A nil HandlerFunc parameter is effecitvely a noop.
func (hs HandlerFuncs) Otherwise(f HandlerFunc) HandlerFunc {
	if f == nil {
		return hs.HandleEvent
	}
	return func(ctx context.Context, e {{.EventType}}) error {
		if h := hs[e.GetType()]; h != nil {
			return h.HandleEvent(ctx, e)
		}
		return f(ctx, e)
	}
}

var (
	_ = Handler(Handlers(nil))
	_ = Handler(HandlerFunc(nil))
	_ = Handler(HandlerFuncs(nil))
)
`))