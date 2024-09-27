package request

import (
	"context"
	"net/http"
)

type MethodInterface interface {
	// Get is request with method GET
	Get(ctx context.Context) ([]byte, int, error)
	// Post is request with method POST
	Post(ctx context.Context, payload []byte) ([]byte, int, error)
	// Put is request with method PUT
	Put(ctx context.Context, payload []byte) ([]byte, int, error)
	// Delete is request with method DELETE
	Delete(ctx context.Context, payload []byte) ([]byte, int, error)
}

func (r *request) Get(ctx context.Context) ([]byte, int, error) {
	return r.wrapper(ctx, nil, http.MethodGet)
}

// Post is request with method POST
func (r *request) Post(ctx context.Context, payload []byte) ([]byte, int, error) {
	return r.wrapper(ctx, payload, http.MethodPost)
}

// Put is request with method PUT
func (r *request) Put(ctx context.Context, payload []byte) ([]byte, int, error) {
	return r.wrapper(ctx, payload, http.MethodPut)
}

// Delete is request with method Delete
func (r *request) Delete(ctx context.Context, payload []byte) ([]byte, int, error) {
	return r.wrapper(ctx, payload, http.MethodDelete)
}
