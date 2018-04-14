// Copyright (c) 2015-2018 Jeevanandam M. (jeeva@myjeeva.com), All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package resty_test

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/proxy"

	"github.com/go-resty/resty"
)

type DropboxError struct {
	Error string
}
type AuthSuccess struct {
	/* variables */
}
type AuthError struct {
	/* variables */
}
type Article struct {
	Title   string
	Content string
	Author  string
	Tags    []string
}
type Error struct {
	/* variables */
}

//
// Package Level examples
//

func Example_get() {
	resp, err := resty.R().Get("http://httpbin.org/get")

	fmt.Printf("\nError: %v", err)
	fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
	fmt.Printf("\nResponse Status: %v", resp.Status())
	fmt.Printf("\nResponse Body: %v", resp)
	fmt.Printf("\nResponse Time: %v", resp.Time())
	fmt.Printf("\nResponse Recevied At: %v", resp.ReceivedAt())
}

func Example_enhancedGet() {
	resp, err := resty.R().
		SetQueryParams(map[string]string{
			"page_no": "1",
			"limit":   "20",
			"sort":    "name",
			"order":   "asc",
			"random":  strconv.FormatInt(time.Now().Unix(), 10),
		}).
		SetHeader("Accept", "application/json").
		SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F").
		Get("/search_result")

	printOutput(resp, err)
}

func Example_post() {
	// POST JSON string
	// No need to set content type, if you have client level setting
	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"username":"testuser", "password":"testpass"}`).
		SetResult(AuthSuccess{}). // or SetResult(&AuthSuccess{}).
		Post("https://myapp.com/login")

	printOutput(resp, err)

	// POST []byte array
	// No need to set content type, if you have client level setting
	resp1, err1 := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody([]byte(`{"username":"testuser", "password":"testpass"}`)).
		SetResult(AuthSuccess{}). // or SetResult(&AuthSuccess{}).
		Post("https://myapp.com/login")

	printOutput(resp1, err1)

	// POST Struct, default is JSON content type. No need to set one
	resp2, err2 := resty.R().
		SetBody(resty.User{Username: "testuser", Password: "testpass"}).
		SetResult(&AuthSuccess{}). // or SetResult(AuthSuccess{}).
		SetError(&AuthError{}).    // or SetError(AuthError{}).
		Post("https://myapp.com/login")

	printOutput(resp2, err2)

	// POST Map, default is JSON content type. No need to set one
	resp3, err3 := resty.R().
		SetBody(map[string]interface{}{"username": "testuser", "password": "testpass"}).
		SetResult(&AuthSuccess{}). // or SetResult(AuthSuccess{}).
		SetError(&AuthError{}).    // or SetError(AuthError{}).
		Post("https://myapp.com/login")

	printOutput(resp3, err3)
}

func Example_dropboxUpload() {
	// For example: upload file to Dropbox
	// POST of raw bytes for file upload.
	file, _ := os.Open("/Users/jeeva/mydocument.pdf")
	fileBytes, _ := ioutil.ReadAll(file)

	// See we are not setting content-type header, since go-resty automatically detects Content-Type for you
	resp, err := resty.R().
		SetBody(fileBytes).     // resty autodetects content type
		SetContentLength(true). // Dropbox expects this value
		SetAuthToken("<your-auth-token>").
		SetError(DropboxError{}).
		Post("https://content.dropboxapi.com/1/files_put/auto/resty/mydocument.pdf") // you can use PUT method too dropbox supports it

	// Output print
	fmt.Printf("\nError: %v\n", err)
	fmt.Printf("Time: %v\n", resp.Time())
	fmt.Printf("Body: %v\n", resp)
}

func Example_put() {
	// Just one sample of PUT, refer POST for more combination
	// request goes as JSON content type
	// No need to set auth token, error, if you have client level settings
	resp, err := resty.R().
		SetBody(Article{
			Title:   "go-resty",
			Content: "This is my article content, oh ya!",
			Author:  "Jeevanandam M",
			Tags:    []string{"article", "sample", "resty"},
		}).
		SetAuthToken("C6A79608-782F-4ED0-A11D-BD82FAD829CD").
		SetError(&Error{}). // or SetError(Error{}).
		Put("https://myapp.com/article/1234")

	printOutput(resp, err)
}

func Example_clientCertificates() {
	// Parsing public/private key pair from a pair of files. The files must contain PEM encoded data.
	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		log.Fatalf("ERROR client certificate: %s", err)
	}

	resty.SetCertificates(cert)
}

func Example_customRootCertificate() {
	resty.SetRootCertificate("/path/to/root/pemFile.pem")
}

//
// top level method examples
//

func ExampleNew() {
	// Creating client1
	client1 := resty.New()
	resp1, err1 := client1.R().Get("http://httpbin.org/get")
	fmt.Println(resp1, err1)

	// Creating client2
	client2 := resty.New()
	resp2, err2 := client2.R().Get("http://httpbin.org/get")
	fmt.Println(resp2, err2)
}

//
// Client object methods
//

func ExampleClient_SetCertificates() {
	// Parsing public/private key pair from a pair of files. The files must contain PEM encoded data.
	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		log.Fatalf("ERROR client certificate: %s", err)
	}

	resty.SetCertificates(cert)
}

//
// Resty Socks5 Proxy request
//

func Example_socks5Proxy() {
	// create a dailer
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9150", nil, proxy.Direct)
	if err != nil {
		log.Fatalf("Unable to obtain proxy dialer: %v\n", err)
	}

	// create a transport
	ptransport := &http.Transport{Dial: dialer.Dial}

	// set transport into resty
	resty.SetTransport(ptransport)

	resp, err := resty.R().Get("http://check.torproject.org")
	fmt.Println(err, resp)
}

func printOutput(resp *resty.Response, err error) {
	fmt.Println(resp, err)
}
