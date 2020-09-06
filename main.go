package main

import (
	"net/http"
	"net/http/httptest"
	"io"
	"io/ioutil"
	"fmt"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"alive": true}`)
	}
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header)
	fmt.Println(string(body))
}