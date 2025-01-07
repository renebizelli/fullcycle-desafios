package cep

type SearchingInfo interface {
	Searching(searchedCEP string, channel chan<- *Response)
}
