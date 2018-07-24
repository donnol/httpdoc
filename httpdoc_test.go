package httpdoc

import (
	"net/http"
	"testing"
)

func TestGenerateHTTPDoc(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	var w http.ResponseWriter
	GenerateHTTPDoc(w, r)
}

func TestWrap(t *testing.T) {
	var h http.Handler
	r := Wrap(h)
	_ = r
}
