package cep

type SearchingInfo interface {
	Searching(searchedCEP string) (*Response, error)
}
