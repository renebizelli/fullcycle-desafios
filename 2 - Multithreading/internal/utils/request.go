package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/renebizelli/goexpert/desafios/multithreading/internal/dtos"
)

func ExecRequestWithContext[T any](ctx context.Context, URL string) (*T, *dtos.RequestError) {

	req, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		return nil, &dtos.RequestError{Message: err.Error(), StatusCode: 400}
	}

	res, err := http.DefaultClient.Do(req)

	select {
	case <-ctx.Done():
		return nil, &dtos.RequestError{Message: err.Error(), StatusCode: 499}
	default:

		if err != nil {
			return nil, &dtos.RequestError{Message: err.Error(), StatusCode: 400}
		}

		if res.StatusCode != 200 {
			fmt.Println(res)
			return nil, &dtos.RequestError{Message: "Invalid request", StatusCode: res.StatusCode}
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, &dtos.RequestError{Message: err.Error(), StatusCode: 500}
		}

		var o *T
		err = json.Unmarshal(body, &o)

		if err != nil {
			return nil, &dtos.RequestError{Message: err.Error(), StatusCode: 500}
		}

		return o, nil
	}

}
