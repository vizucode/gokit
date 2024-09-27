package request

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestGetRequest(t *testing.T) {
	var ctx = context.TODO()

	r := NewRequest(nil)
	r.WithTimeout(5 * time.Second)
	r.WithBasicAuth("username", "password")

	res, sc, err := r.Request(nil, "https://jsonplaceholder.typicode.com/todos/1", "Get:JsonPlaceholder").Get(ctx)

	if err != nil {
		t.Error(err)
	}

	fmt.Println(sc)
	fmt.Println(res)
}
