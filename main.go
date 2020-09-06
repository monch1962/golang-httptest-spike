package main

import (
	"net/http"
	"net/http/httptest"
	// "testing"
	"io"
	"io/ioutil"
	"fmt"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request){
		io.WriteString(w, "ping")
	}
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))
}