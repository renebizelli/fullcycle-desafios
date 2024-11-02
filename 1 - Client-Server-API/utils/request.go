package utils

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type ErrorRequest struct {
	StatusCode int
	Err        string
}

func (e *ErrorRequest) Error() string {
	return e.Err
}

func ExecGetRequestWithContext[T any](ctx context.Context, URL string) (*T, *ErrorRequest) {

	req, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		return nil, &ErrorRequest{Err: err.Error()}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, &ErrorRequest{Err: err.Error()}
	}

	if res.StatusCode != 200 {
		return nil, &ErrorRequest{StatusCode: res.StatusCode}
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &ErrorRequest{Err: err.Error()}
	}

	var o *T
	err = json.Unmarshal(body, &o)

	if err != nil {
		return nil, &ErrorRequest{Err: err.Error()}
	}

	return o, nil

}
