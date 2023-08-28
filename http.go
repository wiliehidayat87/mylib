package mylib

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptrace"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// HttpDial (string)
func HttpDial(url string) error {
	timeout := 10 * time.Second
	conn, err := net.DialTimeout("tcp", url, timeout)
	if err != nil {
		fmt.Printf("Site unreachable : %s, error: %#v\n", url, err)
	}
	defer conn.Close()

	return err
}

// HttpClient (time.Duration, time.Duration, bool)
func HttpClient(p PHttp) *http.Client {
	//ref: Copy and modify defaults from https://golang.org/src/net/http/transport.go
	//Note: Clients and Transports should only be created once and reused
	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			// Modify the time to wait for a connection to establish
			Timeout:   1 * time.Second,
			KeepAlive: p.KeepAlive * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives:   p.IsDisableKeepAlive,
	}

	client := http.Client{
		Transport: &transport,
		Timeout:   p.Timeout * time.Second,
	}

	return &client
}

func (l *Utils) Get(url string, timeout time.Duration) ([]byte, string, string, error) {

	start := time.Now()

	var (
		respBody    []byte
		errHttp     error
		elapseInSec string
		elapseInMS  string
	)

	httpClient := HttpClient(PHttp{Timeout: timeout, KeepAlive: 1, IsDisableKeepAlive: true})

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "x-www-form-urlencoded")
	req.Close = true

	if err != nil {
		l.Write("error",
			fmt.Sprintf("Error Occured : %#v", err),
		)
	}

	var (
		getConn   string
		dnsStart  string
		dnsDone   string
		connStart string
		connDone  string
		gotConn   string
	)

	clientTrace := &httptrace.ClientTrace{
		GetConn:  func(hostPort string) { getConn = fmt.Sprintf("Starting to create conn : [%s], ", hostPort) },
		DNSStart: func(info httptrace.DNSStartInfo) { dnsStart = fmt.Sprintf("starting to look up dns : [%s], ", info) },
		DNSDone:  func(info httptrace.DNSDoneInfo) { dnsDone = fmt.Sprintf("done looking up dns : [%#v], ", info) },
		ConnectStart: func(network, addr string) {
			connStart = fmt.Sprintf("starting tcp connection : [%s, %s], ", network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			connDone = fmt.Sprintf("tcp connection created [%s, %s, %#v], ", network, addr, err)
		},
		GotConn: func(info httptrace.GotConnInfo) { gotConn = fmt.Sprintf("conn was reused: [%#v]", info) },
	}
	clientTraceCtx := httptrace.WithClientTrace(req.Context(), clientTrace)
	req = req.WithContext(clientTraceCtx)

	response, err := httpClient.Do(req)
	if err != nil {
		l.Write("error",
			fmt.Sprintf("Error sending request to API endpoint : %#v", err),
		)
	}

	// Close the connection to reuse it
	defer response.Body.Close()

	respBody, err = io.ReadAll(response.Body)
	if err != nil {
		l.Write("error",
			fmt.Sprintf("Couldn't parse response body : %#v", err),
		)
	}

	elapse := time.Since(start)

	elapseInSec = fmt.Sprintf("%f", elapse.Seconds())
	elapseInMS = strconv.FormatInt(elapse.Milliseconds(), 10)

	l.Write("info",
		fmt.Sprintf("Hit: %s, Response: %s, Elapse: %s second, %s milisecond, live trace : %s", url, string(respBody), elapseInSec, elapseInMS, Concat(getConn, dnsStart, dnsDone, connStart, connDone, gotConn)),
	)

	req = nil
	httpClient = nil

	return respBody, elapseInSec, elapseInMS, errHttp
}

func (l *Utils) Post(url string, headers map[string]string, body []byte, timeout time.Duration) ([]byte, string, string, error) {

	start := time.Now()

	var (
		respBody    []byte
		elapseInSec string
		elapseInMS  string
	)

	httpClient := HttpClient(PHttp{Timeout: timeout, KeepAlive: 1, IsDisableKeepAlive: true})

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	//req.Header.Set("Content-Type", content_type)

	if len(headers) != 0 {
		for k, v := range headers {

			if k == "Basic-Auth" {
				auth := strings.Split(v, ":")
				req.SetBasicAuth(auth[0], auth[1])
			} else {
				req.Header.Set(k, v)
			}
		}
	}

	req.Close = true

	if err != nil {
		l.Write("error",
			fmt.Sprintf("Error Occured : %#v", err),
		)
	}

	var (
		getConn   string
		dnsStart  string
		dnsDone   string
		connStart string
		connDone  string
		gotConn   string
	)

	clientTrace := &httptrace.ClientTrace{
		GetConn:  func(hostPort string) { getConn = fmt.Sprintf("Starting to create conn : [%s], ", hostPort) },
		DNSStart: func(info httptrace.DNSStartInfo) { dnsStart = fmt.Sprintf("starting to look up dns : [%s], ", info) },
		DNSDone:  func(info httptrace.DNSDoneInfo) { dnsDone = fmt.Sprintf("done looking up dns : [%#v], ", info) },
		ConnectStart: func(network, addr string) {
			connStart = fmt.Sprintf("starting tcp connection : [%s, %s], ", network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			connDone = fmt.Sprintf("tcp connection created [%s, %s, %#v], ", network, addr, err)
		},
		GotConn: func(info httptrace.GotConnInfo) { gotConn = fmt.Sprintf("conn was reused: [%#v]", info) },
	}
	clientTraceCtx := httptrace.WithClientTrace(req.Context(), clientTrace)
	req = req.WithContext(clientTraceCtx)

	response, err := httpClient.Do(req)
	if err != nil {
		l.Write("error",
			fmt.Sprintf("Error sending request to API endpoint : %#v", err),
		)
	}

	// Close the connection to reuse it
	defer response.Body.Close()

	respBody, err = io.ReadAll(response.Body)
	if err != nil {
		l.Write("error",
			fmt.Sprintf("Couldn't parse response body : %#v", err),
		)
	}

	elapse := time.Since(start)

	elapseInSec = fmt.Sprintf("%f", elapse.Seconds())
	elapseInMS = strconv.FormatInt(elapse.Milliseconds(), 10)

	l.Write("info",
		fmt.Sprintf("Hit: %s, Request: %s, Response: %s, Elapse: %s second, %s milisecond, live trace : %s", url, string(body), string(respBody), elapseInSec, elapseInMS, Concat(getConn, dnsStart, dnsDone, connStart, connDone, gotConn)),
	)

	req = nil
	httpClient = nil

	return respBody, elapseInSec, elapseInMS, err
}

func (l *Utils) Upload(url string, headers map[string]string, extraParams map[string]string, filepath string, timeout time.Duration) {

	req, err := l.newfileUploadRequest(url, extraParams, "file", filepath)
	if err != nil {
		l.Write("error", fmt.Sprintf("Error writing tmp file : %v, URL : %s, filePath : %s", err, url, filepath))
	}

	if len(headers) != 0 {
		for k, v := range headers {

			if k == "Basic-Auth" {
				auth := strings.Split(v, ":")
				req.SetBasicAuth(auth[0], auth[1])
			} else {
				req.Header.Set(k, v)
			}
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		l.Write("error", fmt.Sprintf("Error upload file : %v, URL : %s", err, url))
	} else {
		/* body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			l.Write("error", fmt.Sprintf("Error reading response : %v, URL : %s", err, url))
		} */
		//resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			l.Write("error",
				fmt.Sprintf("Couldn't parse response body : %#v", err),
			)
		}
		defer resp.Body.Close()

		if err == nil {
			l.Write("info", fmt.Sprintf("Success upload URL : %s, status : %d, header : %#v, response : %#v", url, resp.StatusCode, resp.Header, string(respBody)))
		}
	}
}

func (l *Utils) newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		l.Write("error", fmt.Sprintf("Failed open file : %v, path : %s", err, path))
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		l.Write("error", fmt.Sprintf("Failed create form : %v, param name : %s, base_path : %s", err, paramName, filepath.Base(path)))
		//return nil, err
	}
	_, err = io.Copy(part, file)

	if err != nil {
		l.Write("error", fmt.Sprintf("Failed copy : %v, part : %#v, file : %#v", err, part, file))
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
