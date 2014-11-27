package slack

import (
	"net/http"
)

type Response struct {
	Text string `json:"text"`
}

type Handler interface {
	http.Handler
}
