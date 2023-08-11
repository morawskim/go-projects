package main

import (
	"html/template"
	"net/http"
)

type myFsHandler struct {
	internalHandler http.Handler
	template        *template.Template
	basePath        string
}

func (m *myFsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		m.template.Execute(w, struct {
			BasePath string
		}{m.basePath},
		)
	} else {
		m.internalHandler.ServeHTTP(w, req)
	}
}
