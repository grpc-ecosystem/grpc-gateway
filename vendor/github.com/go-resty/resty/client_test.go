// Copyright (c) 2015-2018 Jeevanandam M (jeeva@myjeeva.com), All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package resty

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestClientBasicAuth(t *testing.T) {
	ts := createAuthServer(t)
	defer ts.Close()

	c := dc()
	c.SetBasicAuth("myuser", "basicauth").
		SetHostURL(ts.URL).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	resp, err := c.R().
		SetResult(&AuthSuccess{}).
		Post("/login")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())

	t.Logf("Result Success: %q", resp.Result().(*AuthSuccess))
	logResponse(t, resp)
}

func TestClientAuthToken(t *testing.T) {
	ts := createAuthServer(t)
	defer ts.Close()

	c := dc()
	c.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetAuthToken("004DDB79-6801-4587-B976-F093E6AC44FF").
		SetHostURL(ts.URL + "/")

	resp, err := c.R().Get("/profile")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
}

func TestOnAfterMiddleware(t *testing.T) {
	ts := createGenServer(t)
	defer ts.Close()

	c := dc()
	c.OnAfterResponse(func(c *Client, res *Response) error {
		t.Logf("Request sent at: %v", res.Request.Time)
		t.Logf("Response Recevied at: %v", res.ReceivedAt())

		return nil
	})

	resp, err := c.R().
		SetBody("OnAfterResponse: This is plain text body to server").
		Put(ts.URL + "/plaintext")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "TestPut: plain text response", resp.String())
}

func TestClientRedirectPolicy(t *testing.T) {
	ts := createRedirectServer(t)
	defer ts.Close()

	c := dc()
	c.SetHTTPMode().
		SetRedirectPolicy(FlexibleRedirectPolicy(20))

	_, err := c.R().Get(ts.URL + "/redirect-1")

	assertEqual(t, "Get /redirect-21: stopped after 20 redirects", err.Error())
}

func TestClientTimeout(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	c := dc()
	c.SetHTTPMode().
		SetTimeout(time.Duration(time.Second * 3))

	_, err := c.R().Get(ts.URL + "/set-timeout-test")
	assertEqual(t, true, strings.Contains(strings.ToLower(err.Error()), "timeout"))
}

func TestClientTimeoutWithinThreshold(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	c := dc()
	c.SetHTTPMode().
		SetTimeout(time.Duration(time.Second * 3))

	resp, err := c.R().Get(ts.URL + "/set-timeout-test-with-sequence")
	assertError(t, err)

	seq1, _ := strconv.ParseInt(resp.String(), 10, 32)

	resp, err = c.R().Get(ts.URL + "/set-timeout-test-with-sequence")
	assertError(t, err)

	seq2, _ := strconv.ParseInt(resp.String(), 10, 32)

	assertEqual(t, seq1+1, seq2)
}

func TestClientTimeoutInternalError(t *testing.T) {
	c := dc()
	c.SetHTTPMode()
	c.SetTimeout(time.Duration(time.Second * 1))

	_, _ = c.R().Get("http://localhost:9000/set-timeout-test")
}

func TestClientProxy(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	c := dc()
	c.SetTimeout(1 * time.Second)
	c.SetProxy("http://sampleproxy:8888")

	resp, err := c.R().Get(ts.URL)
	assertNotNil(t, resp)
	assertNotNil(t, err)

	// Error
	c.SetProxy("//not.a.user@%66%6f%6f.com:8888")

	resp, err = c.R().
		Get(ts.URL)
	assertNil(t, err)
	assertNotNil(t, resp)
}

func TestClientSetCertificates(t *testing.T) {
	DefaultClient = dc()
	SetCertificates(tls.Certificate{})

	transport, err := DefaultClient.getTransport()

	assertNil(t, err)
	assertEqual(t, 1, len(transport.TLSClientConfig.Certificates))
}

func TestClientSetRootCertificate(t *testing.T) {
	DefaultClient = dc()
	SetRootCertificate(getTestDataPath() + "/sample-root.pem")

	transport, err := DefaultClient.getTransport()

	assertNil(t, err)
	assertNotNil(t, transport.TLSClientConfig.RootCAs)
}

func TestClientSetRootCertificateNotExists(t *testing.T) {
	DefaultClient = dc()
	SetRootCertificate(getTestDataPath() + "/not-exists-sample-root.pem")

	transport, err := DefaultClient.getTransport()

	assertNil(t, err)
	assertNil(t, transport.TLSClientConfig)
}

func TestClientOnBeforeRequestModification(t *testing.T) {
	tc := New()
	tc.OnBeforeRequest(func(c *Client, r *Request) error {
		r.SetAuthToken("This is test auth token")
		return nil
	})

	ts := createGetServer(t)
	defer ts.Close()

	resp, err := tc.R().Get(ts.URL + "/")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "200 OK", resp.Status())
	assertNotNil(t, resp.Body())
	assertEqual(t, "TestGet: text response", resp.String())

	logResponse(t, resp)
}

func TestClientSetTransport(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()
	DefaultClient = dc()

	transport := &http.Transport{
		// something like Proxying to httptest.Server, etc...
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(ts.URL)
		},
	}
	SetTransport(transport)

	transportInUse, err := DefaultClient.getTransport()

	assertNil(t, err)

	assertEqual(t, true, transport == transportInUse)
}

func TestClientSetScheme(t *testing.T) {
	DefaultClient = dc()

	SetScheme("http")

	assertEqual(t, true, DefaultClient.scheme == "http")
}

func TestClientSetCookieJar(t *testing.T) {
	DefaultClient = dc()
	backupJar := DefaultClient.httpClient.Jar

	SetCookieJar(nil)
	assertNil(t, DefaultClient.httpClient.Jar)

	SetCookieJar(backupJar)
	assertEqual(t, true, DefaultClient.httpClient.Jar == backupJar)
}

func TestClientOptions(t *testing.T) {
	SetHTTPMode().SetContentLength(true)
	assertEqual(t, Mode(), "http")
	assertEqual(t, DefaultClient.setContentLength, true)

	SetRESTMode()
	assertEqual(t, Mode(), "rest")

	SetHostURL("http://httpbin.org")
	assertEqual(t, "http://httpbin.org", DefaultClient.HostURL)

	SetHeader(hdrContentTypeKey, jsonContentType)
	SetHeaders(map[string]string{
		hdrUserAgentKey: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) go-resty v0.1",
		"X-Request-Id":  strconv.FormatInt(time.Now().UnixNano(), 10),
	})
	assertEqual(t, jsonContentType, DefaultClient.Header.Get(hdrContentTypeKey))

	SetCookie(&http.Cookie{
		Name:     "default-cookie",
		Value:    "This is cookie default-cookie value",
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   36000,
		HttpOnly: true,
		Secure:   false,
	})
	assertEqual(t, "default-cookie", DefaultClient.Cookies[0].Name)

	var cookies []*http.Cookie
	cookies = append(cookies, &http.Cookie{
		Name:  "default-cookie-1",
		Value: "This is default-cookie 1 value",
		Path:  "/",
	})
	cookies = append(cookies, &http.Cookie{
		Name:  "default-cookie-2",
		Value: "This is default-cookie 2 value",
		Path:  "/",
	})
	SetCookies(cookies)
	assertEqual(t, "default-cookie-1", DefaultClient.Cookies[1].Name)
	assertEqual(t, "default-cookie-2", DefaultClient.Cookies[2].Name)

	SetQueryParam("test_param_1", "Param_1")
	SetQueryParams(map[string]string{"test_param_2": "Param_2", "test_param_3": "Param_3"})
	assertEqual(t, "Param_3", DefaultClient.QueryParam.Get("test_param_3"))

	rTime := strconv.FormatInt(time.Now().UnixNano(), 10)
	SetFormData(map[string]string{"r_time": rTime})
	assertEqual(t, rTime, DefaultClient.FormData.Get("r_time"))

	SetBasicAuth("myuser", "mypass")
	assertEqual(t, "myuser", DefaultClient.UserInfo.Username)

	SetAuthToken("AC75BD37F019E08FBC594900518B4F7E")
	assertEqual(t, "AC75BD37F019E08FBC594900518B4F7E", DefaultClient.Token)

	SetDisableWarn(true)
	assertEqual(t, DefaultClient.DisableWarn, true)

	SetRetryCount(3)
	assertEqual(t, 3, DefaultClient.RetryCount)

	rwt := time.Duration(1000) * time.Millisecond
	SetRetryWaitTime(rwt)
	assertEqual(t, rwt, DefaultClient.RetryWaitTime)

	mrwt := time.Duration(2) * time.Second
	SetRetryMaxWaitTime(mrwt)
	assertEqual(t, mrwt, DefaultClient.RetryMaxWaitTime)

	err := &AuthError{}
	SetError(err)
	if reflect.TypeOf(err) == DefaultClient.Error {
		t.Error("SetError failed")
	}

	SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	transport, transportErr := DefaultClient.getTransport()

	assertNil(t, transportErr)
	assertEqual(t, true, transport.TLSClientConfig.InsecureSkipVerify)

	OnBeforeRequest(func(c *Client, r *Request) error {
		c.Log.Println("I'm in Request middleware")
		return nil // if it success
	})
	OnAfterResponse(func(c *Client, r *Response) error {
		c.Log.Println("I'm in Response middleware")
		return nil // if it success
	})

	SetTimeout(time.Duration(5 * time.Second))
	SetRedirectPolicy(FlexibleRedirectPolicy(10), func(req *http.Request, via []*http.Request) error {
		return errors.New("sample test redirect")
	})
	SetContentLength(true)

	SetDebug(true)
	assertEqual(t, DefaultClient.Debug, true)

	var sl int64 = 1000000
	SetDebugBodyLimit(sl)
	assertEqual(t, DefaultClient.debugBodySizeLimit, sl)

	SetAllowGetMethodPayload(true)
	assertEqual(t, DefaultClient.AllowGetMethodPayload, true)

	SetScheme("http")
	assertEqual(t, DefaultClient.scheme, "http")

	SetCloseConnection(true)
	assertEqual(t, DefaultClient.closeConnection, true)

	SetLogger(ioutil.Discard)
}

func TestClientPreRequestHook(t *testing.T) {
	SetPreRequestHook(func(c *Client, r *Request) error {
		c.Log.Println("I'm in Pre-Request Hook")
		return nil
	})

	SetPreRequestHook(func(c *Client, r *Request) error {
		c.Log.Println("I'm Overwriting existing Pre-Request Hook")
		return nil
	})
}

func TestClientAllowsGetMethodPayload(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	c := dc()
	c.SetAllowGetMethodPayload(true)
	c.SetPreRequestHook(func(*Client, *Request) error { return nil }) // for coverage

	payload := "test-payload"
	resp, err := c.R().SetBody(payload).Get(ts.URL + "/get-method-payload-test")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, payload, resp.String())
}

func TestClientRoundTripper(t *testing.T) {
	c := NewWithClient(&http.Client{})

	rt := &CustomRoundTripper{}
	c.SetTransport(rt)

	ct, err := c.getTransport()
	assertNotNil(t, err)
	assertNil(t, ct)
	assertEqual(t, "current transport is not an *http.Transport instance", err.Error())

	c.SetTLSClientConfig(&tls.Config{})
	c.SetProxy("http://localhost:9090")
	c.RemoveProxy()
	c.SetCertificates(tls.Certificate{})
	c.SetRootCertificate(getTestDataPath() + "/sample-root.pem")
}

func TestClientNewRequest(t *testing.T) {
	c := New()
	request := c.NewRequest()

	assertNotNil(t, request)
}

func TestNewRequest(t *testing.T) {
	request := NewRequest()

	assertNotNil(t, request)
}

func TestDebugBodySizeLimit(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	var lgr bytes.Buffer
	c := dc()
	c.SetDebug(true)
	c.SetLogger(&lgr)
	c.SetDebugBodyLimit(30)

	testcases := []struct{ url, want string }{
		// Text, does not exceed limit.
		{ts.URL, "TestGet: text response"},
		// Empty response.
		{ts.URL + "/no-content", "***** NO CONTENT *****"},
		// JSON, does not exceed limit.
		{ts.URL + "/json", "{\n   \"TestGet\": \"JSON response\"\n}"},
		// Invalid JSON, does not exceed limit.
		{ts.URL + "/json-invalid", "TestGet: Invalid JSON"},
		// Text, exceeds limit.
		{ts.URL + "/long-text", "RESPONSE TOO LARGE"},
		// JSON, exceeds limit.
		{ts.URL + "/long-json", "RESPONSE TOO LARGE"},
	}

	for _, tc := range testcases {
		_, err := c.R().Get(tc.url)
		assertError(t, err)
		debugLog := lgr.String()
		if !strings.Contains(debugLog, tc.want) {
			t.Errorf("Expected logs to contain [%v], got [\n%v]", tc.want, debugLog)
		}
		lgr.Reset()
	}
}

// CustomRoundTripper just for test
type CustomRoundTripper struct {
}

// RoundTrip just for test
func (rt *CustomRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	return &http.Response{}, nil
}

func TestSetLogPrefix(t *testing.T) {
	c := New()
	c.SetLogPrefix("CUSTOM ")
	assertEqual(t, "CUSTOM ", c.logPrefix)
	assertEqual(t, "CUSTOM ", c.Log.Prefix())

	c.disableLogPrefix()
	c.enableLogPrefix()
	assertEqual(t, "CUSTOM ", c.logPrefix)
	assertEqual(t, "CUSTOM ", c.Log.Prefix())
}
