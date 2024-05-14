package main

import "net/http"

func addRoutes(mux *http.ServeMux, cfg Config, svc FilevaultService) {
	mux.Handle("POST /api/v1", uploadFile(svc))
	mux.Handle("GET /api/v1", getFile(svc))
	mux.Handle("DELETE /api/v1", removeFile(svc))
	mux.Handle("GET /api/v1/healthz", health(cfg))
}
