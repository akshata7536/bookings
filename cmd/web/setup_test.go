package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type myHandler struct{}

// ServeHTTP implements http.Handler.
func (mh *myHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	panic("unimplemented")
}

func (mh *myHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {

}
