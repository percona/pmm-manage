package main

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// according to http://jsonapi.org/format/#errors
type jsonResponce struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
}
