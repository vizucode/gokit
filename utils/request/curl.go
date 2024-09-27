package request

import (
	"bytes"
	"io"
	"net/http"
)

func (r *request) do(payload []byte, method string) ([]byte, int, error) {
	req, err := http.NewRequest(method, r.url, buf(payload))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if r.header != nil {
		req.Header = r.header
	}

	// set basic auth if exists
	if r.basicAuth.set {
		req.SetBasicAuth(r.basicAuth.username, r.basicAuth.username)
	}

	// set timeout if exists
	if r.timeout.set {
		r.client.Timeout = r.timeout.duration
	}

	// do request to client
	res, err := r.client.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return body, res.StatusCode, nil
}

func buf(p []byte) io.ReadCloser {
	if p != nil {
		r := bytes.NewReader(p)
		return io.NopCloser(r)
	}

	return nil
}
