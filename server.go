package main

import (
	"net/http"
)

func NewServer(cfg Config, svc FilevaultService) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, cfg, svc)
	return mux
}
