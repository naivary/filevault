package main

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	maxInputLen          = 32 << 20
	formFileKey          = "file"
	formPathKey          = "path"
	headerKeyContentType = "Content-Type"
)

type ErrorResponse struct {
	Message string
	Path    string
}

func newErrorResponse(msg, path string) ErrorResponse {
	return ErrorResponse{Message: msg, Path: path}
}

func getFile(svc FilevaultService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			return
		}
		path := r.FormValue(formPathKey)
		data, err := svc.ReadFile(path)
		if err != nil {
			httperr := newErrorResponse(err.Error(), path)
			encode(w, r, http.StatusBadRequest, httperr)
			return
		}
		filename := filepath.Base(path)
		ext := filepath.Ext(filename)
		contentType := mime.TypeByExtension(ext)

		w.Header().Set(headerKeyContentType, contentType)
		w.WriteHeader(http.StatusOK)
		_, err = io.Copy(w, bytes.NewReader(data))
		if err != nil {
			httperr := newErrorResponse(err.Error(), path)
			encode(w, r, http.StatusInternalServerError, httperr)
			return
		}
	})
}

func removeFile(svc FilevaultService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			httperr := newErrorResponse(err.Error(), "")
			encode(w, r, http.StatusBadRequest, httperr)
			return
		}
		path := r.FormValue(formPathKey)
		if err := svc.RemoveFile(path); err != nil {
			httperr := newErrorResponse(err.Error(), path)
			encode(w, r, http.StatusBadRequest, httperr)
			return
		}
        w.WriteHeader(http.StatusNoContent)
	})
}

func uploadFile(svc FilevaultService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type resp struct {
			Dir      string
			Filename string
			Path     string
		}
		if err := r.ParseMultipartForm(maxInputLen); err != nil {
			httperr := newErrorResponse(err.Error(), "")
			encode(w, r, http.StatusBadRequest, httperr)
			return
		}
		file, header, err := r.FormFile(formFileKey)
		if err != nil {
			msg := fmt.Sprintf("missing form key '%s' for file", formFileKey)
			httperr := newErrorResponse(msg, "")
			encode(w, r, http.StatusBadRequest, httperr)
			return
		}
		defer file.Close()
		_, _, ok := strings.Cut(header.Filename, ".")
		if !ok {
            httperr := newErrorResponse("filename has to be in format of  <name>.<ext>. For examle filevault.json", "")
			encode(w, r, http.StatusBadRequest, httperr)
			return
		}
		dir := r.FormValue("dir")
		path := filepath.Join(dir, header.Filename)
		if err := svc.CreateFile(path, file); err != nil {
            httperr := newErrorResponse(err.Error(), "")
			encode(w, r, http.StatusBadRequest, httperr)
			return
		}
		res := resp{
			Dir:      filepath.Dir(path),
			Filename: header.Filename,
			Path:     path,
		}
		encode(w, r, http.StatusCreated, res)
	})
}

func health(cfg Config) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if ok, err := isDirExisting(cfg.Dir); err != nil  || !ok {
            httperr := newErrorResponse("directory does not exist", cfg.Dir)
			encode(w, r, http.StatusBadRequest, httperr)
        }
        w.WriteHeader(http.StatusOK)
    })
}
