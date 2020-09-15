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
	"fmt"

	"github.com/remeh/sizedwaitgroup"
)

func setConcurrency() (int) {
	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))
	if err != nil {
		concurrency= 1
	}
	return concurrency
}

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

type ReqHeader struct {
	Key string
	Value string
}

type RespHeader struct {
	Key string
	Value []string
}

func containsHeader(returnedHeaders http.Header, h RespHeader) bool {
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

func logConcurrency(c int, t *testing.T) {
	t.Logf("Running %d tests concurrently\n", c)
}

func examineResponse(r *http.Response,t *testing.T, xpectBody interface{}, expectHeaders []RespHeader, expectStatus int) {
	expectBody, _ := xpectBody.(string)
	//expectBody := fmt.Sprintf("%s", xpectBody)

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

type PactDetail struct {
	testName string
	verb string
	endPoint string
	reqBody string
	reqHeaders []ReqHeader
	status int
	respBody interface{}
	respHeaders []RespHeader
}

func TestGenerated(t *testing.T) {
	var testCases = []PactDetail{
		{"TestA", "GET", "/posts/1", "", []ReqHeader{{"abc", "def"}}, 200, "", []RespHeader{{"Content-Type",[]string{"application/json; charset=utf-8"}}}},
		{"TestB", "GET", "/posts/2", "", []ReqHeader{{"abc", "def"}}, 201, "", []RespHeader{{"Content-Type",[]string{"application/json; charset=utf-8"}}}},
	}

	baseURL := getBaseURL(t)
	concurrency := setConcurrency()
	logConcurrency(concurrency, t)

	swg := sizedwaitgroup.New(concurrency)
	// These are shown separately when executing go test -v.

    for testNo, tt := range testCases {
		swg.Add()
		go func(t *testing.T, pactNo int, tt PactDetail) {
			defer swg.Done()
        	testname := fmt.Sprintf("%d - %s", pactNo, tt.testName)
        	t.Run(testname, func(t *testing.T) {

				client := constructHTTPClient(t)

				req,err := http.NewRequest(tt.verb, baseURL + tt.endPoint, nil)
				if err != nil {
					t.Fatalf("Unable to construct request: '%v %v'\n%v\n",tt.verb, baseURL+tt.endPoint, err.Error())
				}

				for _,h := range tt.reqHeaders {
					req.Header.Set(h.Key, h.Value)
				}
				logRequest(req, t)

				response,err := client.Do(req)
				if err != nil {
					t.Fatalf("%v\n",err.Error())
				}
				logResponse(response, t)

				expectHeaders := []RespHeader{}
				for _,h := range tt.respHeaders {
					header := RespHeader{h.Key, h.Value}
					expectHeaders = append(expectHeaders, header)
				}
				examineResponse(response, t, tt.respBody, tt.respHeaders, tt.status)
			})
		}(t, testNo, tt)
	}
	swg.Wait()
}
