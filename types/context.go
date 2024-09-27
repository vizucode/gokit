package types

import (
	"bytes"
	"context"
)

type EventContext struct {
	ctx          context.Context
	workerType   string
	handlerRoute string
	header       map[string]string
	key          string
	err          error
	buff         *bytes.Buffer
}

// SetContext setter context
func (e *EventContext) SetContext(ctx context.Context) {
	e.ctx = ctx
}

// SetWorkerType setter worker type
func (e *EventContext) SetWorkerType(wt string) {
	e.workerType = wt
}

// SetHandlerRoute setter handler route
func (e *EventContext) SetHandlerRoute(h string) {
	e.handlerRoute = h
}

// SetHeader setter header
func (e *EventContext) SetHeader(header map[string]string) {
	e.header = header
}

// SetKey setter key context
func (e *EventContext) SetKey(key string) {
	e.key = key
}

// SetError setter error
func (e *EventContext) SetError(err error) {
	e.err = err
}

// Context get current context
func (e *EventContext) Context() context.Context {
	return e.ctx
}

// WorkerType get worker type
func (e *EventContext) WorkerType() string {
	return e.workerType
}

// Header get header
func (e *EventContext) Header() map[string]string {
	return e.header
}

// Key get key
func (e *EventContext) Key() string {
	return e.key
}

// Message context
func (e *EventContext) Message() []byte {
	return e.buff.Bytes()
}

// Err get error
func (e *EventContext) Err() error {
	return e.err
}

// Read buffer
func (e *EventContext) Read(p []byte) (int, error) {
	return e.buff.Read(p)
}

// Write buffer
func (e *EventContext) Write(p []byte) (int, error) {
	if e.buff == nil {
		e.buff = &bytes.Buffer{}
	}

	return e.buff.Write(p)
}

// WriteString buffer
func (e *EventContext) WriteString(p string) (int, error) {
	if e.buff == nil {
		e.buff = &bytes.Buffer{}
	}

	return e.buff.WriteString(p)
}
