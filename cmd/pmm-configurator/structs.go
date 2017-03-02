package main

type sshkey struct {
	Key         string `json:"key,omitempty"`
	Type        string `json:"type,omitempty"`
	Comment     string `json:"comment,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
}

// according to http://jsonapi.org/format/#errors
type jsonResponce struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
}
