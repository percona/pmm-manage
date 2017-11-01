package main

// according to http://jsonapi.org/format/#errors
type jsonResponce struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type updateResponce struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
}
