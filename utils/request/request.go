package request

import (
	"net/http"
	"time"
)

type request struct {
	url           string
	header        http.Header
	serviceTarget string
	timeout
	client    *http.Client
	basicAuth basicAuth
}

type timeout struct {
	duration time.Duration
	set      bool
}

type basicAuth struct {
	username string
	password string
	set      bool
}

func (r *request) WithTimeout(d time.Duration) {
	r.timeout.duration = d
	r.timeout.set = true
}

func (r *request) WithBasicAuth(username, password string) {
	r.basicAuth.username = username
	r.basicAuth.password = password
	r.basicAuth.set = true
}

type Client interface {
	Request(header http.Header, url string, serviceTarget string) MethodInterface
	WithTimeout(d time.Duration)
	WithBasicAuth(username, password string)
}

func NewRequest(client *http.Client) Client {
	if client == nil {
		client = http.DefaultClient
	}

	r := new(request)
	r.client = client
	return r
}

func (r *request) Request(header http.Header, url string, serviceTarget string) MethodInterface {
	r.header = header
	r.url = url
	r.serviceTarget = serviceTarget
	return r
}
