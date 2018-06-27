package main

// according to http://jsonapi.org/format/#errors
type jsonResponce struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type instance struct {
	ID string `json:"InstanceID"`
}

type updateResponce struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
	Step   string `json:"step,omitempty"`
}

type versionResponce struct {
	Version       string `json:"version"`
	ReleaseDate   string `json:"release_date,omitempty"`
	DisableUpdate bool   `json:"disable_update,omitempty"`
}
