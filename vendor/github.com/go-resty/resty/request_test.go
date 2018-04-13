// Copyright (c) 2015-2018 Jeevanandam M (jeeva@myjeeva.com), All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package resty

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

type AuthSuccess struct {
	ID, Message string
}

type AuthError struct {
	ID, Message string
}

func TestGet(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	resp, err := R().
		SetQueryParam("request_no", strconv.FormatInt(time.Now().Unix(), 10)).
		Get(ts.URL + "/")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "200 OK", resp.Status())
	assertNotNil(t, resp.Body())
	assertEqual(t, "TestGet: text response", resp.String())

	logResponse(t, resp)
}

func TestGetCustomUserAgent(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	resp, err := dcr().
		SetHeader(hdrUserAgentKey, "Test Custom User agent").
		SetQueryParam("request_no", strconv.FormatInt(time.Now().Unix(), 10)).
		Get(ts.URL + "/")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "200 OK", resp.Status())
	assertEqual(t, "TestGet: text response", resp.String())

	logResponse(t, resp)
}

func TestGetClientParamRequestParam(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	c := dc()
	c.SetQueryParam("client_param", "true").
		SetQueryParams(map[string]string{"req_1": "jeeva", "req_3": "jeeva3"}).
		SetDebug(true).
		SetLogger(ioutil.Discard)

	resp, err := c.R().
		SetQueryParams(map[string]string{"req_1": "req 1 value", "req_2": "req 2 value"}).
		SetQueryParam("request_no", strconv.FormatInt(time.Now().Unix(), 10)).
		SetHeader(hdrUserAgentKey, "Test Custom User agent").
		Get(ts.URL + "/")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "200 OK", resp.Status())
	assertEqual(t, "TestGet: text response", resp.String())

	logResponse(t, resp)
}

func TestGetRelativePath(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	c := dc()
	c.SetHostURL(ts.URL)

	resp, err := c.R().Get("mypage2")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "TestGet: text response from mypage2", resp.String())

	logResponse(t, resp)
}

func TestGet400Error(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	resp, err := dcr().Get(ts.URL + "/mypage")

	assertError(t, err)
	assertEqual(t, http.StatusBadRequest, resp.StatusCode())
	assertEqual(t, "", resp.String())

	logResponse(t, resp)
}

func TestPostJSONStringSuccess(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	c := dc()
	c.SetHeader(hdrContentTypeKey, jsonContentType).
		SetHeaders(map[string]string{hdrUserAgentKey: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) go-resty v0.1", hdrAcceptKey: jsonContentType})

	resp, err := c.R().
		SetBody(`{"username":"testuser", "password":"testpass"}`).
		Post(ts.URL + "/login")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())

	logResponse(t, resp)

	// PostJSONStringError
	resp, err = c.R().
		SetBody(`{"username":"testuser" "password":"testpass"}`).
		Post(ts.URL + "/login")

	assertError(t, err)
	assertEqual(t, http.StatusBadRequest, resp.StatusCode())

	logResponse(t, resp)
}

func TestPostJSONBytesSuccess(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	c := dc()
	c.SetHeader(hdrContentTypeKey, jsonContentType).
		SetHeaders(map[string]string{hdrUserAgentKey: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) go-resty v0.7", hdrAcceptKey: jsonContentType})

	resp, err := c.R().
		SetBody([]byte(`{"username":"testuser", "password":"testpass"}`)).
		Post(ts.URL + "/login")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())

	logResponse(t, resp)
}

func TestPostJSONBytesIoReader(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	c := dc()
	c.SetHeader(hdrContentTypeKey, jsonContentType)

	bodyBytes := []byte(`{"username":"testuser", "password":"testpass"}`)

	resp, err := c.R().
		SetBody(bytes.NewReader(bodyBytes)).
		Post(ts.URL + "/login")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())

	logResponse(t, resp)
}

func TestPostJSONStructSuccess(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	user := &User{Username: "testuser", Password: "testpass"}

	c := dc()
	resp, err := c.R().
		SetHeader(hdrContentTypeKey, jsonContentType).
		SetBody(user).
		SetResult(&AuthSuccess{}).
		Post(ts.URL + "/login")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())

	t.Logf("Result Success: %q", resp.Result().(*AuthSuccess))

	logResponse(t, resp)
}

func TestPostJSONStructInvalidLogin(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	c := dc()
	c.SetDebug(false)

	resp, err := c.R().
		SetHeader(hdrContentTypeKey, jsonContentType).
		SetBody(User{Username: "testuser", Password: "testpass1"}).
		SetError(AuthError{}).
		Post(ts.URL + "/login")

	assertError(t, err)
	assertEqual(t, http.StatusUnauthorized, resp.StatusCode())

	authError := resp.Error().(*AuthError)
	assertEqual(t, "unauthorized", authError.ID)
	assertEqual(t, "Invalid credentials", authError.Message)
	t.Logf("Result Error: %q", resp.Error().(*AuthError))

	logResponse(t, resp)
}

func TestPostJSONErrorRFC7807(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	c := dc()
	resp, err := c.R().
		SetHeader(hdrContentTypeKey, jsonContentType).
		SetBody(User{Username: "testuser", Password: "testpass1"}).
		SetError(AuthError{}).
		Post(ts.URL + "/login?ct=problem")

	assertError(t, err)
	assertEqual(t, http.StatusUnauthorized, resp.StatusCode())

	authError := resp.Error().(*AuthError)
	assertEqual(t, "unauthorized", authError.ID)
	assertEqual(t, "Invalid credentials", authError.Message)
	t.Logf("Result Error: %q", resp.Error().(*AuthError))

	logResponse(t, resp)
}

func TestPostJSONMapSuccess(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	c := dc()
	c.SetDebug(false)

	resp, err := c.R().
		SetBody(map[string]interface{}{"username": "testuser", "password": "testpass"}).
		SetResult(AuthSuccess{}).
		Post(ts.URL + "/login")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())

	t.Logf("Result Success: %q", resp.Result().(*AuthSuccess))

	logResponse(t, resp)
}

func TestPostJSONMapInvalidResponseJson(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	resp, err := dclr().
		SetBody(map[string]interface{}{"username": "testuser", "password": "invalidjson"}).
		SetResult(&AuthSuccess{}).
		Post(ts.URL + "/login")

	assertEqual(t, "invalid character '}' looking for beginning of object key string", err.Error())
	assertEqual(t, http.StatusOK, resp.StatusCode())

	authSuccess := resp.Result().(*AuthSuccess)
	assertEqual(t, "", authSuccess.ID)
	assertEqual(t, "", authSuccess.Message)

	t.Logf("Result Success: %q", resp.Result().(*AuthSuccess))

	logResponse(t, resp)
}

func TestPostXMLStringSuccess(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	c := dc()
	c.SetDebug(false)

	resp, err := c.R().
		SetHeader(hdrContentTypeKey, "application/xml").
		SetBody(`<?xml version="1.0" encoding="UTF-8"?><User><Username>testuser</Username><Password>testpass</Password></User>`).
		SetQueryParam("request_no", strconv.FormatInt(time.Now().Unix(), 10)).
		Post(ts.URL + "/login")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())

	logResponse(t, resp)
}

func TestPostXMLStringError(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	resp, err := dclr().
		SetHeader(hdrContentTypeKey, "application/xml").
		SetBody(`<?xml version="1.0" encoding="UTF-8"?><User><Username>testuser</Username>testpass</Password></User>`).
		Post(ts.URL + "/login")

	assertError(t, err)
	assertEqual(t, http.StatusBadRequest, resp.StatusCode())
	assertEqual(t, `<?xml version="1.0" encoding="UTF-8"?><AuthError><Id>bad_request</Id><Message>Unable to read user info</Message></AuthError>`, resp.String())

	logResponse(t, resp)
}

func TestPostXMLBytesSuccess(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	c := dc()
	c.SetDebug(false)

	resp, err := c.R().
		SetHeader(hdrContentTypeKey, "application/xml").
		SetBody([]byte(`<?xml version="1.0" encoding="UTF-8"?><User><Username>testuser</Username><Password>testpass</Password></User>`)).
		SetQueryParam("request_no", strconv.FormatInt(time.Now().Unix(), 10)).
		SetContentLength(true).
		Post(ts.URL + "/login")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())

	logResponse(t, resp)
}

func TestPostXMLStructSuccess(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	resp, err := dclr().
		SetHeader(hdrContentTypeKey, "application/xml").
		SetBody(User{Username: "testuser", Password: "testpass"}).
		SetContentLength(true).
		SetResult(&AuthSuccess{}).
		Post(ts.URL + "/login")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())

	t.Logf("Result Success: %q", resp.Result().(*AuthSuccess))

	logResponse(t, resp)
}

func TestPostXMLStructInvalidLogin(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	c := dc()
	c.SetError(&AuthError{})

	resp, err := c.R().
		SetHeader(hdrContentTypeKey, "application/xml").
		SetBody(User{Username: "testuser", Password: "testpass1"}).
		Post(ts.URL + "/login")

	assertError(t, err)
	assertEqual(t, http.StatusUnauthorized, resp.StatusCode())
	assertEqual(t, resp.Header().Get("Www-Authenticate"), "Protected Realm")

	t.Logf("Result Error: %q", resp.Error().(*AuthError))

	logResponse(t, resp)
}

func TestPostXMLStructInvalidResponseXml(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	resp, err := dclr().
		SetHeader(hdrContentTypeKey, "application/xml").
		SetBody(User{Username: "testuser", Password: "invalidxml"}).
		SetResult(&AuthSuccess{}).
		Post(ts.URL + "/login")

	assertEqual(t, "XML syntax error on line 1: element <Message> closed by </AuthSuccess>", err.Error())
	assertEqual(t, http.StatusOK, resp.StatusCode())

	t.Logf("Result Success: %q", resp.Result().(*AuthSuccess))

	logResponse(t, resp)
}

func TestPostXMLMapNotSupported(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	_, err := dclr().
		SetHeader(hdrContentTypeKey, "application/xml").
		SetBody(map[string]interface{}{"Username": "testuser", "Password": "testpass"}).
		Post(ts.URL + "/login")

	assertEqual(t, "unsupported 'Body' type/value", err.Error())
}

func TestRequestBasicAuth(t *testing.T) {
	ts := createAuthServer(t)
	defer ts.Close()

	c := dc()
	c.SetHostURL(ts.URL).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	resp, err := c.R().
		SetBasicAuth("myuser", "basicauth").
		SetResult(&AuthSuccess{}).
		Post("/login")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())

	t.Logf("Result Success: %q", resp.Result().(*AuthSuccess))
	logResponse(t, resp)
}

func TestRequestBasicAuthFail(t *testing.T) {
	ts := createAuthServer(t)
	defer ts.Close()

	c := dc()
	c.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetError(AuthError{})

	resp, err := c.R().
		SetBasicAuth("myuser", "basicauth1").
		Post(ts.URL + "/login")

	assertError(t, err)
	assertEqual(t, http.StatusUnauthorized, resp.StatusCode())

	t.Logf("Result Error: %q", resp.Error().(*AuthError))
	logResponse(t, resp)
}

func TestRequestAuthToken(t *testing.T) {
	ts := createAuthServer(t)
	defer ts.Close()

	c := dc()
	c.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetAuthToken("004DDB79-6801-4587-B976-F093E6AC44FF")

	resp, err := c.R().
		SetAuthToken("004DDB79-6801-4587-B976-F093E6AC44FF-Request").
		Get(ts.URL + "/profile")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
}

func TestFormData(t *testing.T) {
	ts := createFormPostServer(t)
	defer ts.Close()

	c := dc()
	c.SetFormData(map[string]string{"zip_code": "00000", "city": "Los Angeles"}).
		SetContentLength(true).
		SetDebug(true).
		SetLogger(ioutil.Discard)

	resp, err := c.R().
		SetFormData(map[string]string{"first_name": "Jeevanandam", "last_name": "M", "zip_code": "00001"}).
		SetBasicAuth("myuser", "mypass").
		Post(ts.URL + "/profile")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "Success", resp.String())
}

func TestMultiValueFormData(t *testing.T) {
	ts := createFormPostServer(t)
	defer ts.Close()

	v := url.Values{
		"search_criteria": []string{"book", "glass", "pencil"},
	}

	c := dc()
	c.SetContentLength(true).
		SetDebug(true).
		SetLogger(ioutil.Discard)

	resp, err := c.R().
		SetMultiValueFormData(v).
		Post(ts.URL + "/search")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "Success", resp.String())
}

func TestFormDataDisableWarn(t *testing.T) {
	ts := createFormPostServer(t)
	defer ts.Close()

	c := dc()
	c.SetFormData(map[string]string{"zip_code": "00000", "city": "Los Angeles"}).
		SetContentLength(true).
		SetDebug(true).
		SetLogger(ioutil.Discard).
		SetDisableWarn(true)

	resp, err := c.R().
		SetFormData(map[string]string{"first_name": "Jeevanandam", "last_name": "M", "zip_code": "00001"}).
		SetBasicAuth("myuser", "mypass").
		Post(ts.URL + "/profile")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "Success", resp.String())
}

func TestMultiPartUploadFile(t *testing.T) {
	ts := createFormPostServer(t)
	defer ts.Close()
	defer cleanupFiles("test-data/upload")

	basePath := getTestDataPath()

	c := dc()
	c.SetFormData(map[string]string{"zip_code": "00001", "city": "Los Angeles"})

	resp, err := c.R().
		SetFile("profile_img", basePath+"/test-img.png").
		SetContentLength(true).
		Post(ts.URL + "/upload")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
}

func TestMultiPartUploadFileError(t *testing.T) {
	ts := createFormPostServer(t)
	defer ts.Close()
	defer cleanupFiles("test-data/upload")

	basePath := getTestDataPath()

	c := dc()
	c.SetFormData(map[string]string{"zip_code": "00001", "city": "Los Angeles"})

	resp, err := c.R().
		SetFile("profile_img", basePath+"/test-img-not-exists.png").
		Post(ts.URL + "/upload")

	if err == nil {
		t.Errorf("Expected [%v], got [%v]", nil, err)
	}
	if resp != nil {
		t.Errorf("Expected [%v], got [%v]", nil, resp)
	}
}

func TestMultiPartUploadFiles(t *testing.T) {
	ts := createFormPostServer(t)
	defer ts.Close()
	defer cleanupFiles("test-data/upload")

	basePath := getTestDataPath()

	resp, err := dclr().
		SetFormData(map[string]string{"first_name": "Jeevanandam", "last_name": "M"}).
		SetFiles(map[string]string{"profile_img": basePath + "/test-img.png", "notes": basePath + "/text-file.txt"}).
		Post(ts.URL + "/upload")

	responseStr := resp.String()

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, true, strings.Contains(responseStr, "test-img.png"))
	assertEqual(t, true, strings.Contains(responseStr, "text-file.txt"))
}

func TestMultiPartIoReaderFiles(t *testing.T) {
	ts := createFormPostServer(t)
	defer ts.Close()
	defer cleanupFiles("test-data/upload")

	basePath := getTestDataPath()
	profileImgBytes, _ := ioutil.ReadFile(basePath + "/test-img.png")
	notesBytes, _ := ioutil.ReadFile(basePath + "/text-file.txt")

	// Just info values
	file := File{
		Name:      "test_file_name.jpg",
		ParamName: "test_param",
		Reader:    bytes.NewBuffer([]byte("test bytes")),
	}
	t.Logf("File Info: %v", file.String())

	resp, err := dclr().
		SetFormData(map[string]string{"first_name": "Jeevanandam", "last_name": "M"}).
		SetFileReader("profile_img", "test-img.png", bytes.NewReader(profileImgBytes)).
		SetFileReader("notes", "text-file.txt", bytes.NewReader(notesBytes)).
		Post(ts.URL + "/upload")

	responseStr := resp.String()

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, true, strings.Contains(responseStr, "test-img.png"))
	assertEqual(t, true, strings.Contains(responseStr, "text-file.txt"))
}

func TestMultiPartUploadFileNotOnGetOrDelete(t *testing.T) {
	ts := createFormPostServer(t)
	defer ts.Close()
	defer cleanupFiles("test-data/upload")

	basePath := getTestDataPath()

	_, err := dclr().
		SetFile("profile_img", basePath+"/test-img.png").
		Get(ts.URL + "/upload")

	assertEqual(t, "multipart content is not allowed in HTTP verb [GET]", err.Error())

	_, err = dclr().
		SetFile("profile_img", basePath+"/test-img.png").
		Delete(ts.URL + "/upload")

	assertEqual(t, "multipart content is not allowed in HTTP verb [DELETE]", err.Error())
}

func TestMultiPartMultipartField(t *testing.T) {
	ts := createFormPostServer(t)
	defer ts.Close()
	defer cleanupFiles("test-data/upload")

	jsonBytes := []byte(`{"input": {"name": "Uploaded document", "_filename" : ["file.txt"]}}`)

	resp, err := dclr().
		SetFormData(map[string]string{"first_name": "Jeevanandam", "last_name": "M"}).
		SetMultipartField("uploadManifest", "upload-file.json", "application/json", bytes.NewReader(jsonBytes)).
		Post(ts.URL + "/upload")

	responseStr := resp.String()

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, true, strings.Contains(responseStr, "upload-file.json"))
}

func TestGetWithCookie(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	c := dc()
	c.SetHostURL(ts.URL)
	c.SetCookie(&http.Cookie{
		Name:     "go-resty-1",
		Value:    "This is cookie 1 value",
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   36000,
		HttpOnly: true,
		Secure:   false,
	})

	resp, err := c.R().Get("mypage2")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "TestGet: text response from mypage2", resp.String())

	logResponse(t, resp)
}

func TestGetWithCookies(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	var cookies []*http.Cookie

	cookies = append(cookies, &http.Cookie{
		Name:     "go-resty-1",
		Value:    "This is cookie 1 value",
		Path:     "/",
		Domain:   "sample.com",
		MaxAge:   36000,
		HttpOnly: true,
		Secure:   false,
	})

	cookies = append(cookies, &http.Cookie{
		Name:     "go-resty-2",
		Value:    "This is cookie 2 value",
		Path:     "/",
		Domain:   "sample.com",
		MaxAge:   36000,
		HttpOnly: true,
		Secure:   false,
	})

	c := dc()
	c.SetHostURL(ts.URL).
		SetCookies(cookies)

	resp, err := c.R().Get("mypage2")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "TestGet: text response from mypage2", resp.String())

	logResponse(t, resp)
}

func TestPutPlainString(t *testing.T) {
	ts := createGenServer(t)
	defer ts.Close()

	resp, err := R().
		SetBody("This is plain text body to server").
		Put(ts.URL + "/plaintext")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "TestPut: plain text response", resp.String())
}

func TestPutJSONString(t *testing.T) {
	ts := createGenServer(t)
	defer ts.Close()

	DefaultClient.OnBeforeRequest(func(c *Client, r *Request) error {
		r.SetHeader("X-Custom-Request-Middleware", "OnBeforeRequest middleware")
		return nil
	})
	DefaultClient.OnBeforeRequest(func(c *Client, r *Request) error {
		c.SetContentLength(true)
		r.SetHeader("X-ContentLength", "OnBeforeRequest ContentLength set")
		return nil
	})

	DefaultClient.SetDebug(true).SetLogger(ioutil.Discard)

	resp, err := R().
		SetHeaders(map[string]string{hdrContentTypeKey: jsonContentType, hdrAcceptKey: jsonContentType}).
		SetBody(`{"content":"json content sending to server"}`).
		Put(ts.URL + "/json")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, `{"response":"json response"}`, resp.String())
}

func TestPutXMLString(t *testing.T) {
	ts := createGenServer(t)
	defer ts.Close()

	resp, err := R().
		SetHeaders(map[string]string{hdrContentTypeKey: "application/xml", hdrAcceptKey: "application/xml"}).
		SetBody(`<?xml version="1.0" encoding="UTF-8"?><Request>XML Content sending to server</Request>`).
		Put(ts.URL + "/xml")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, `<?xml version="1.0" encoding="UTF-8"?><Response>XML response</Response>`, resp.String())
}

func TestOnBeforeMiddleware(t *testing.T) {
	ts := createGenServer(t)
	defer ts.Close()

	c := dc()
	c.OnBeforeRequest(func(c *Client, r *Request) error {
		r.SetHeader("X-Custom-Request-Middleware", "OnBeforeRequest middleware")
		return nil
	})
	c.OnBeforeRequest(func(c *Client, r *Request) error {
		c.SetContentLength(true)
		r.SetHeader("X-ContentLength", "OnBeforeRequest ContentLength set")
		return nil
	})

	resp, err := c.R().
		SetBody("OnBeforeRequest: This is plain text body to server").
		Put(ts.URL + "/plaintext")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "TestPut: plain text response", resp.String())
}

func TestNoAutoRedirect(t *testing.T) {
	ts := createRedirectServer(t)
	defer ts.Close()

	_, err := R().Get(ts.URL + "/redirect-1")

	assertEqual(t, "Get /redirect-2: auto redirect is disabled", err.Error())
}

func TestHTTPAutoRedirectUpTo10(t *testing.T) {
	ts := createRedirectServer(t)
	defer ts.Close()

	c := dc()
	c.SetHTTPMode()
	_, err := c.R().Get(ts.URL + "/redirect-1")

	assertEqual(t, "Get /redirect-11: stopped after 10 redirects", err.Error())
}

func TestHostCheckRedirectPolicy(t *testing.T) {
	ts := createRedirectServer(t)
	defer ts.Close()

	c := dc().
		SetRedirectPolicy(DomainCheckRedirectPolicy("127.0.0.1"))

	_, err := c.R().Get(ts.URL + "/redirect-host-check-1")

	assertNotNil(t, err)
	assertEqual(t, true, strings.Contains(err.Error(), "redirect is not allowed as per DomainCheckRedirectPolicy"))
}

func TestHeadMethod(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	resp, err := dclr().Head(ts.URL + "/")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
}

func TestOptionsMethod(t *testing.T) {
	ts := createGenServer(t)
	defer ts.Close()

	resp, err := dclr().Options(ts.URL + "/options")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, resp.Header().Get("Access-Control-Expose-Headers"), "x-go-resty-id")
}

func TestPatchMethod(t *testing.T) {
	ts := createGenServer(t)
	defer ts.Close()

	resp, err := dclr().Patch(ts.URL + "/patch")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())

	resp.body = nil
	assertEqual(t, "", resp.String())
}

func TestRawFileUploadByBody(t *testing.T) {
	ts := createFormPostServer(t)
	defer ts.Close()

	file, _ := os.Open(getTestDataPath() + "/test-img.png")
	fileBytes, _ := ioutil.ReadAll(file)

	resp, err := dclr().
		SetBody(fileBytes).
		SetContentLength(true).
		SetAuthToken("004DDB79-6801-4587-B976-F093E6AC44FF").
		Put(ts.URL + "/raw-upload")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "image/png", resp.Request.Header.Get(hdrContentTypeKey))
}

func TestProxySetting(t *testing.T) {
	c := dc()

	transport, err := c.getTransport()

	assertNil(t, err)

	assertEqual(t, false, c.IsProxySet())
	assertNil(t, transport.Proxy)

	c.SetProxy("http://sampleproxy:8888")
	assertEqual(t, true, c.IsProxySet())
	assertNotNil(t, transport.Proxy)

	c.SetProxy("//not.a.user@%66%6f%6f.com:8888")
	assertEqual(t, false, c.IsProxySet())
	assertNil(t, transport.Proxy)

	SetProxy("http://sampleproxy:8888")
	assertEqual(t, true, IsProxySet())
	RemoveProxy()
	assertNil(t, DefaultClient.proxyURL)
	assertNil(t, transport.Proxy)
}

func TestGetClient(t *testing.T) {
	client := GetClient()
	custom := New()
	customClient := custom.GetClient()

	assertNotNil(t, client)
	assertNotNil(t, customClient)

	assertNotEqual(t, client, http.DefaultClient)
	assertNotEqual(t, customClient, http.DefaultClient)
	assertNotEqual(t, client, customClient)

	assertEqual(t, DefaultClient.httpClient, client)
}

func TestIncorrectURL(t *testing.T) {
	_, err := R().Get("//not.a.user@%66%6f%6f.com/just/a/path/also")
	assertEqual(t, true, strings.Contains(err.Error(), "parse //not.a.user@%66%6f%6f.com/just/a/path/also"))

	c := dc()
	c.SetHostURL("//not.a.user@%66%6f%6f.com")
	_, err1 := c.R().Get("/just/a/path/also")
	assertEqual(t, true, strings.Contains(err1.Error(), "parse //not.a.user@%66%6f%6f.com/just/a/path/also"))
}

func TestDetectContentTypeForPointer(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	user := &User{Username: "testuser", Password: "testpass"}

	resp, err := dclr().
		SetBody(user).
		SetResult(AuthSuccess{}).
		Post(ts.URL + "/login")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())

	t.Logf("Result Success: %q", resp.Result().(*AuthSuccess))

	logResponse(t, resp)
}

type ExampleUser struct {
	FirstName string `json:"frist_name"`
	LastName  string `json:"last_name"`
	ZipCode   string `json:"zip_code"`
}

func TestDetectContentTypeForPointerWithSlice(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	users := &[]ExampleUser{
		{FirstName: "firstname1", LastName: "lastname1", ZipCode: "10001"},
		{FirstName: "firstname2", LastName: "lastname3", ZipCode: "10002"},
		{FirstName: "firstname3", LastName: "lastname3", ZipCode: "10003"},
	}

	resp, err := dclr().
		SetBody(users).
		Post(ts.URL + "/users")

	assertError(t, err)
	assertEqual(t, http.StatusAccepted, resp.StatusCode())

	t.Logf("Result Success: %q", resp)

	logResponse(t, resp)
}

func TestDetectContentTypeForPointerWithSliceMap(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	usersmap := map[string]interface{}{
		"user1": ExampleUser{FirstName: "firstname1", LastName: "lastname1", ZipCode: "10001"},
		"user2": &ExampleUser{FirstName: "firstname2", LastName: "lastname3", ZipCode: "10002"},
		"user3": ExampleUser{FirstName: "firstname3", LastName: "lastname3", ZipCode: "10003"},
	}

	var users []map[string]interface{}
	users = append(users, usersmap)

	resp, err := dclr().
		SetBody(&users).
		Post(ts.URL + "/usersmap")

	assertError(t, err)
	assertEqual(t, http.StatusAccepted, resp.StatusCode())

	t.Logf("Result Success: %q", resp)

	logResponse(t, resp)
}

func TestDetectContentTypeForSlice(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	users := []ExampleUser{
		{FirstName: "firstname1", LastName: "lastname1", ZipCode: "10001"},
		{FirstName: "firstname2", LastName: "lastname3", ZipCode: "10002"},
		{FirstName: "firstname3", LastName: "lastname3", ZipCode: "10003"},
	}

	resp, err := dclr().
		SetBody(users).
		Post(ts.URL + "/users")

	assertError(t, err)
	assertEqual(t, http.StatusAccepted, resp.StatusCode())

	t.Logf("Result Success: %q", resp)

	logResponse(t, resp)
}

func TestMultiParamsQueryString(t *testing.T) {
	ts1 := createGetServer(t)
	defer ts1.Close()

	client := dc()
	req1 := client.R()

	client.SetQueryParam("status", "open")

	_, _ = req1.SetQueryParam("status", "pending").
		Get(ts1.URL)

	assertEqual(t, true, strings.Contains(req1.URL, "status=pending"))
	// pending overrides open
	assertEqual(t, false, strings.Contains(req1.URL, "status=open"))

	_, _ = req1.SetQueryParam("status", "approved").
		Get(ts1.URL)

	assertEqual(t, true, strings.Contains(req1.URL, "status=approved"))
	// approved overrides pending
	assertEqual(t, false, strings.Contains(req1.URL, "status=pending"))

	ts2 := createGetServer(t)
	defer ts2.Close()

	req2 := client.R()

	v := url.Values{
		"status": []string{"pending", "approved", "reject"},
	}

	_, _ = req2.SetMultiValueQueryParams(v).Get(ts2.URL)

	assertEqual(t, true, strings.Contains(req2.URL, "status=pending"))
	assertEqual(t, true, strings.Contains(req2.URL, "status=approved"))
	assertEqual(t, true, strings.Contains(req2.URL, "status=reject"))

	// because it's removed by key
	assertEqual(t, false, strings.Contains(req2.URL, "status=open"))
}

func TestSetQueryStringTypical(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	resp, err := dclr().
		SetQueryString("productId=232&template=fresh-sample&cat=resty&source=google&kw=buy a lot more").
		Get(ts.URL)

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "200 OK", resp.Status())
	assertEqual(t, "TestGet: text response", resp.String())

	resp, err = dclr().
		SetQueryString("&%%amp;").
		Get(ts.URL)

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "200 OK", resp.Status())
	assertEqual(t, "TestGet: text response", resp.String())
}

func TestOutputFileWithBaseDirAndRelativePath(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()
	defer cleanupFiles("test-data/dir-sample")

	DefaultClient = dc()
	SetRedirectPolicy(FlexibleRedirectPolicy(10))
	SetOutputDirectory(getTestDataPath() + "/dir-sample")
	SetDebug(true)

	resp, err := R().
		SetOutput("go-resty/test-img-success.png").
		Get(ts.URL + "/my-image.png")

	assertError(t, err)
	assertEqual(t, true, resp.Size() != 0)
}

func TestOutputFileWithBaseDirError(t *testing.T) {
	c := dc().SetRedirectPolicy(FlexibleRedirectPolicy(10)).
		SetOutputDirectory(getTestDataPath() + `/go-resty\0`)

	_ = c
}

func TestOutputPathDirNotExists(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()
	defer cleanupFiles("test-data/not-exists-dir")

	DefaultClient = dc()
	SetRedirectPolicy(FlexibleRedirectPolicy(10))
	SetOutputDirectory(getTestDataPath() + "/not-exists-dir")

	resp, err := R().
		SetOutput("test-img-success.png").
		Get(ts.URL + "/my-image.png")

	assertError(t, err)
	assertEqual(t, true, resp.Size() != 0)
}

func TestOutputFileAbsPath(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()
	defer cleanupFiles("test-data/go-resty")

	_, err := dcr().
		SetOutput(getTestDataPath() + "/go-resty/test-img-success-2.png").
		Get(ts.URL + "/my-image.png")

	assertError(t, err)
}

func TestContextInternal(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	r := R().
		SetQueryParam("request_no", strconv.FormatInt(time.Now().Unix(), 10))

	if r.isContextCancelledIfAvailable() {
		t.Error("isContextCancelledIfAvailable != false for vanilla R()")
	}
	r.addContextIfAvailable()

	resp, err := r.Get(ts.URL + "/")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
}

func TestSRV(t *testing.T) {
	c := dc().
		SetRedirectPolicy(FlexibleRedirectPolicy(20)).
		SetScheme("http")

	r := c.R().
		SetSRV(&SRVRecord{"xmpp-server", "google.com"})

	assertEqual(t, "xmpp-server", r.SRV.Service)
	assertEqual(t, "google.com", r.SRV.Domain)

	resp, err := r.Get("/")
	if err == nil {
		assertError(t, err)
		assertNotNil(t, resp)
		assertEqual(t, http.StatusOK, resp.StatusCode())
	}
}

func TestSRVInvalidService(t *testing.T) {
	_, err := R().
		SetSRV(&SRVRecord{"nonexistantservice", "sampledomain"}).
		Get("/")

	assertNotNil(t, err)
	assertEqual(t, true, strings.Contains(err.Error(), "no such host"))
}

func TestDeprecatedCodeCoverage(t *testing.T) {
	var user1 User
	err := Unmarshal("application/json",
		[]byte(`{"username":"testuser", "password":"testpass"}`), &user1)
	assertError(t, err)
	assertEqual(t, "testuser", user1.Username)
	assertEqual(t, "testpass", user1.Password)

	var user2 User
	err = Unmarshal("application/xml",
		[]byte(`<?xml version="1.0" encoding="UTF-8"?><User><Username>testuser</Username><Password>testpass</Password></User>`),
		&user2)
	assertError(t, err)
	assertEqual(t, "testuser", user1.Username)
	assertEqual(t, "testpass", user1.Password)
}

func TestRequestDoNotParseResponse(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	resp, err := dc().R().
		SetDoNotParseResponse(true).
		SetQueryParam("request_no", strconv.FormatInt(time.Now().Unix(), 10)).
		Get(ts.URL + "/")

	assertError(t, err)

	buf := acquireBuffer()
	defer releaseBuffer(buf)
	_, _ = io.Copy(buf, resp.RawBody())

	assertEqual(t, "TestGet: text response", buf.String())
	_ = resp.RawBody().Close()

	// Manually setting RawResponse as nil
	resp, err = dc().R().
		SetDoNotParseResponse(true).
		Get(ts.URL + "/")

	assertError(t, err)

	resp.RawResponse = nil
	assertNil(t, resp.RawBody())

	// just set test part
	SetDoNotParseResponse(true)
	assertEqual(t, true, DefaultClient.notParseResponse)
	SetDoNotParseResponse(false)
}

type noCtTest struct {
	Response string `json:"response"`
}

func TestRequestExpectContentTypeTest(t *testing.T) {
	ts := createGenServer(t)
	defer ts.Close()

	c := dc()
	resp, err := c.R().
		SetResult(noCtTest{}).
		ExpectContentType("application/json").
		Get(ts.URL + "/json-no-set")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertNotNil(t, resp.Result())
	assertEqual(t, "json response no content type set", resp.Result().(*noCtTest).Response)

	assertEqual(t, "", firstNonEmpty("", ""))
}

func TestGetPathParams(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	c := dc()
	c.SetHostURL(ts.URL).
		SetPathParams(map[string]string{
			"userId": "sample@sample.com",
		})

	resp, err := c.R().SetPathParams(map[string]string{
		"subAccountId": "100002",
	}).
		Get("v1/users/{userId}/{subAccountId}/details")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "TestPathParams: text response", resp.String())

	logResponse(t, resp)

	SetPathParams(map[string]string{
		"userId": "sample@sample.com",
	})
}

func TestReportMethodSupportsPayload(t *testing.T) {
	ts := createGenServer(t)
	defer ts.Close()

	c := dc()
	resp, err := c.R().
		SetBody("body").
		Execute("REPORT", ts.URL+"/report")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())

}

func TestRequestQueryStringOrder(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	resp, err := New().R().
		SetQueryString("productId=232&template=fresh-sample&cat=resty&source=google&kw=buy a lot more").
		Get(ts.URL + "/?UniqueId=ead1d0ed-XXX-XXX-XXX-abb7612b3146&Translate=false&tempauth=eyJ0eXAiOiJKV1QiLC...HZEhwVnJ1d0NSUGVLaUpSaVNLRG5scz0&ApiVersion=2.0")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "200 OK", resp.Status())
	assertNotNil(t, resp.Body())
	assertEqual(t, "TestGet: text response", resp.String())

	logResponse(t, resp)
}

func TestRequestOverridesClientAuthorizationHeader(t *testing.T) {
	ts := createAuthServer(t)
	defer ts.Close()

	c := dc()
	c.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetHeader("Authorization", "some token")
	SetHostURL(ts.URL + "/")

	resp, err := c.R().
		SetHeader("Authorization", "Bearer 004DDB79-6801-4587-B976-F093E6AC44FF").
		Get("/profile")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
}

func TestRequestFileUploadAsReader(t *testing.T) {
	ts := createFilePostServer(t)
	defer ts.Close()

	file, _ := os.Open(getTestDataPath() + "/test-img.png")
	defer file.Close()

	resp, err := dclr().
		SetBody(file).
		SetHeader("Content-Type", "image/png").
		Post(ts.URL + "/upload")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, true, strings.Contains(resp.String(), "File Uploaded successfully"))

	file, _ = os.Open(getTestDataPath() + "/test-img.png")
	defer file.Close()

	resp, err = dclr().
		SetBody(file).
		SetHeader("Content-Type", "image/png").
		SetContentLength(true).
		Post(ts.URL + "/upload")

	assertError(t, err)
	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, true, strings.Contains(resp.String(), "File Uploaded successfully"))
}
