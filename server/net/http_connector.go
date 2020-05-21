package webserver

import "net/http"

//HttpServer defines the interface to server HTTP request
type HttpServer interface {
	Handle(pattern string, handler http.Handler)
	ListenAndServe(addr string, handler http.Handler) error
}

//HttpServerWrapper is a wrapper around the http package
type HttpServerWrapper struct {
}

//Handle is a wrapper around the same method from http package
func (server *HttpServerWrapper) Handle(pattern string, handler http.Handler) {
	http.Handle(pattern, handler)
}

//ListenAndServe is a wrapper around the same method from http package
func (server *HttpServerWrapper) ListenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}
