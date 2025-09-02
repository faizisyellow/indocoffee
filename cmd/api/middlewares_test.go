package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func FirstMiddleware(next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("i go first")

		next.ServeHTTP(w, r)
	}
}

func SecondMiddleware(next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("i go second")

		next.ServeHTTP(w, r)
	}
}

func TestMiddlewareChain(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hello world")
	})

	handler := NewHandlerFunc(FirstMiddleware, SecondMiddleware)(finalHandler)
	handler.ServeHTTP(httptest.NewRecorder(), request)
}
