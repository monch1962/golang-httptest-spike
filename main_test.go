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

/*func containsHeader(arr http.Header, h Header) bool {
	for _, a := range arr {
    	if a == h {
        	return true
      	}
   	}
   	return false
}*/

func examineResponse(r *http.Response,t *testing.T, expectBody string, expectHeaders []Header, expectStatus int) {
	t.Logf("Response code: %d\n", r.StatusCode)
	t.Logf("Response headers: %v\n", r.Header)
	respHeadersJSON, _ := json.MarshalIndent(r.Header, "", "  ")
	t.Logf("Response headers (JSON): %v\n", string(respHeadersJSON))
	buf, _ := ioutil.ReadAll(r.Body)
	t.Logf("Response body: %s\n", buf)
	if r.StatusCode != expectStatus {
		t.Fatalf("Expected response code %d, received response code %d\n", expectStatus, r.StatusCode)
	}
	if expectBody != "" && string(buf) != expectBody {
		t.Fatalf("Expected body '%v', received body '%v'\n", expectBody, string(buf))
	}
	//for _, h := range expectHeaders {
		//if !containsHeader(r.Header, h) {
		//	t.Fatalf("Expected header '%v' not found\n", h)
		//}
	//}
}

func TestEndpoint(t *testing.T) {
	baseURL := getBaseURL(t)
	client := constructHTTPClient(t)

	req, _ := http.NewRequest("GET", baseURL + "/posts/1", nil)

	req.Header.Set("abc", "123")
	logRequest(req, t)

	response,err := client.Do(req)
	if err != nil {
		t.Fatalf("%v\n",err.Error())
	}

	expectStatus := 200
	expectBody := "blah"
	expectHeaders := []Header{}
	header := Header{"abc", []string{"def"}}
	expectHeaders = append(expectHeaders, header)
	examineResponse(response, t, expectBody, expectHeaders, expectStatus)
}