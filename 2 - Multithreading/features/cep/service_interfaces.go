package cep

import "context"

type SearchingInfo interface {
	Searching(ctx context.Context, searchedCEP string, channel chan<- *Response)
}
