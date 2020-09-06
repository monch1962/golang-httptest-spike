package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
	"net/http/httputil"
	//"fmt"
)

func setTimeouts() (time.Duration, time.Duration, time.Duration) {
	httpTimeout, err := strconv.Atoi(os.Getenv("HTTP_TIMEOUT"))
	if err != nil {
		httpTimeout = 30000
	}
	tcpConnectTimeout, err := strconv.Atoi(os.Getenv("TCP_CONNECT_TIMEOUT"))
	if err != nil {
		tcpConnectTimeout = 30000
	}
	tlsConnectTimeout, err := strconv.Atoi(os.Getenv("TLS_CONNECT_TIMEOUT"))
	if err != nil {
		tlsConnectTimeout = 30000
	}
	return time.Duration(httpTimeout)*time.Millisecond, 
		time.Duration(tcpConnectTimeout)*time.Millisecond, 
		time.Duration(tlsConnectTimeout)*time.Millisecond
}

func constructHTTPClient(t *testing.T) http.Client {
	httpTimeout, tcpConnectTimeout, tlsConnectTimeout := setTimeouts()

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: tcpConnectTimeout,
		}).Dial,
		TLSHandshakeTimeout: tlsConnectTimeout,
	}
	client := http.Client{
		Timeout: httpTimeout,
		Transport: netTransport,
	}
	return client
}

func logRequest(req *http.Request, t *testing.T) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
  		t.Fatal(err)
	}
	t.Log(string(requestDump))
}

func getBaseURL(t *testing.T) string {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		t.Fatal("No BASE_URL environment variable supplied")
	}
	return baseURL
}

type Header struct {
	Key string
	Value []string
}

func containsHeader(returnedHeaders http.Header, h Header) bool {
	for k,v:= range returnedHeaders {
    	if k == h.Key && v[0] == h.Value[0] {
        	return true
      	}
   	}
	return false
}

func logResponse(r *http.Response, t *testing.T) {
	t.Logf("Response code: %d\n", r.StatusCode)
	t.Logf("Response headers: %v\n", r.Header)
	respHeadersJSON, _ := json.MarshalIndent(r.Header, "", "  ")
	t.Logf("Response headers (JSON): %v\n", string(respHeadersJSON))
	buf, _ := ioutil.ReadAll(r.Body)
	t.Logf("Response body: %s\n", buf)
}

func examineResponse(r *http.Response,t *testing.T, expectBody string, expectHeaders []Header, expectStatus int) {
	buf, _ := ioutil.ReadAll(r.Body)
	if r.StatusCode != expectStatus {
		t.Fatalf("Expected response code %d, received response code %d\n", expectStatus, r.StatusCode)
	}
	if expectBody != "" && string(buf) != expectBody {
		t.Fatalf("Expected body '%v', received body '%v'\n", expectBody, string(buf))
	}
	for _, h := range expectHeaders {
		if !containsHeader(r.Header, h) {
			t.Fatalf("Expected header '%v: %v' not found\n", h.Key, h.Value[0])
		}
	}
}

func TestEndpoint(t *testing.T) {
	baseURL := getBaseURL(t)
	client := constructHTTPClient(t)

	endPoint := "/posts/1"
	verb := "GET"

	req,err := http.NewRequest(verb, baseURL + endPoint, nil)
	if err != nil {
		t.Fatalf("Unable to construct request: '%v %v'\n%v\n",verb, baseURL+endPoint, err.Error())
	}

	req.Header.Set("abc", "123")
	req.Header.Set("def", "456")
	logRequest(req, t)

	response,err := client.Do(req)
	if err != nil {
		t.Fatalf("%v\n",err.Error())
	}
	logResponse(response, t)

	expectStatus := 200
	expectBody := ""
	expectHeaders := []Header{}
	//header := Header{"abc", []string{"def"}}
	header := Header{"Content-Type", []string{"application/json; charset=utf-8"}}
	expectHeaders = append(expectHeaders, header)
	examineResponse(response, t, expectBody, expectHeaders, expectStatus)
}