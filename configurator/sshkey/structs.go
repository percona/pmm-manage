package sshkey

type Handler struct {
	KeyPath  string
	KeyOwner string
}

type Key struct {
	Key         string `json:"key,omitempty"`
	Type        string `json:"type,omitempty"`
	Comment     string `json:"comment,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
}
