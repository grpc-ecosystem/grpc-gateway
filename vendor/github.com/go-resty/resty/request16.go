// +build !go1.7

// Copyright (c) 2015-2018 Jeevanandam M (jeeva@myjeeva.com)
// 2016 Andrew Grigorev (https://github.com/ei-grad)
// All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package resty

import (
	"bytes"
	"net/http"
	"net/url"
	"time"
)

// Request type is used to compose and send individual request from client
// go-resty is provide option override client level settings such as
//		Auth Token, Basic Auth credentials, Header, Query Param, Form Data, Error object
// and also you can add more options for that particular request
//
type Request struct {
	URL        string
	Method     string
	QueryParam url.Values
	FormData   url.Values
	Header     http.Header
	UserInfo   *User
	Token      string
	Body       interface{}
	Result     interface{}
	Error      interface{}
	Time       time.Time
	RawRequest *http.Request
	SRV        *SRVRecord

	client              *Client
	bodyBuf             *bytes.Buffer
	isMultiPart         bool
	isFormData          bool
	setContentLength    bool
	isSaveResponse      bool
	outputFile          string
	multipartFiles      []*File
	multipartFields     []*multipartField
	notParseResponse    bool
	fallbackContentType string
	pathParams          map[string]string
}

func (r *Request) addContextIfAvailable() {
	// nothing to do for golang<1.7
}

func (r *Request) isContextCancelledIfAvailable() bool {
	// just always return false golang<1.7
	return false
}
