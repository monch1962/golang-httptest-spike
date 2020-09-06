package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
	"net"
	"os"
	"strconv"
)

//var baseURL string = "https://jsonplaceholder.typicode.com"
//var httpTimeout time.Duration = func _() int64 {
//	strconv.Atoi(os.Getenv("HTTP_TIMEOUT")) * time.Millisecond
//}
//var tcpConnectTimeout time.Duration = strconv.Atoi(os.Getenv("TCP_CONNECT_TIMEOUT")) * time.Millisecond
//var tlsConnectTimeout time.Duration = strconv.Atoi(os.Getenv("TLS_CONNECT_TIMEOUT")) * time.Millisecond

func getParameters() (string, time.Duration, time.Duration, time.Duration) {
	baseURL := os.Getenv("BASE_URL")

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
	return baseURL, time.Duration(httpTimeout)*time.Millisecond, 
		time.Duration(tcpConnectTimeout)*time.Millisecond, 
		time.Duration(tlsConnectTimeout)*time.Millisecond
}

func TestNothing(t *testing.T) {
	t.Fail()
}

/*func BenchmarkRandInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = 1
	}
}

func ExampleHello() {
	fmt.Println("hello")
	// Output: hello
}
*/

func TestEndpoint(t *testing.T) {
	baseURL, httpTimeout, tcpConnectTimeout, tlsConnectTimeout := getParameters()
	if baseURL == "" {
		t.Error("No BASE_URL environment variable supplied")
		t.Fail()
		return
	}
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: tcpConnectTimeout,
		}).Dial,
		TLSHandshakeTimeout: tlsConnectTimeout,
	}
	var netClient = &http.Client{
		Timeout: httpTimeout,
		Transport: netTransport,
	}
	response,err := netClient.Get(baseURL + "/posts/1")
	if err != nil {
		t.Errorf("%v\n",err.Error())
		t.Fail()
		return
	}
	buf, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("%s\n", buf)
}