/*
 * Unannotated Echo
 *
 * Unannotated Echo Service Similar to echo_service.proto but without annotations. See unannotated_echo_service.yaml for the equivalent of the annotations in gRPC API configuration format.  Echo Service API consists of a single service which returns a message.
 *
 * API version: 1.0
 * Contact: none@example.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package unannotatedecho

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"fmt"
	"github.com/antihax/optional"
)

// Linger please
var (
	_ context.Context
)

type UnannotatedEchoServiceApiService service

/* 
UnannotatedEchoServiceApiService Summary: Echo rpc
Description Echo
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param id Id represents the message identifier.
 * @param num Int value field
 * @param optional nil or *UnannotatedEchoServiceEchoOpts - Optional Parameters:
     * @param "Duration" (optional.String) - 
     * @param "LineNum" (optional.String) - 
     * @param "Lang" (optional.String) - 
     * @param "StatusProgress" (optional.String) - 
     * @param "StatusNote" (optional.String) - 
     * @param "En" (optional.String) - 
     * @param "NoProgress" (optional.String) - 
     * @param "NoNote" (optional.String) - 

@return ExamplepbUnannotatedSimpleMessage
*/

type UnannotatedEchoServiceEchoOpts struct { 
	Duration optional.String
	LineNum optional.String
	Lang optional.String
	StatusProgress optional.String
	StatusNote optional.String
	En optional.String
	NoProgress optional.String
	NoNote optional.String
}

func (a *UnannotatedEchoServiceApiService) UnannotatedEchoServiceEcho(ctx context.Context, id string, num string, localVarOptionals *UnannotatedEchoServiceEchoOpts) (ExamplepbUnannotatedSimpleMessage, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Post")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue ExamplepbUnannotatedSimpleMessage
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/v1/example/echo/{id}"
	localVarPath = strings.Replace(localVarPath, "{"+"id"+"}", fmt.Sprintf("%v", id), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	localVarQueryParams.Add("num", parameterToString(num, ""))
	if localVarOptionals != nil && localVarOptionals.Duration.IsSet() {
		localVarQueryParams.Add("duration", parameterToString(localVarOptionals.Duration.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.LineNum.IsSet() {
		localVarQueryParams.Add("lineNum", parameterToString(localVarOptionals.LineNum.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Lang.IsSet() {
		localVarQueryParams.Add("lang", parameterToString(localVarOptionals.Lang.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.StatusProgress.IsSet() {
		localVarQueryParams.Add("status.progress", parameterToString(localVarOptionals.StatusProgress.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.StatusNote.IsSet() {
		localVarQueryParams.Add("status.note", parameterToString(localVarOptionals.StatusNote.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.En.IsSet() {
		localVarQueryParams.Add("en", parameterToString(localVarOptionals.En.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.NoProgress.IsSet() {
		localVarQueryParams.Add("no.progress", parameterToString(localVarOptionals.NoProgress.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.NoNote.IsSet() {
		localVarQueryParams.Add("no.note", parameterToString(localVarOptionals.NoNote.Value(), ""))
	}
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json", "application/x-foo-mime"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json", "application/x-foo-mime"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	if ctx != nil {
		// API Key Authentication
		if auth, ok := ctx.Value(ContextAPIKey).(APIKey); ok {
			var key string
			if auth.Prefix != "" {
				key = auth.Prefix + " " + auth.Key
			} else {
				key = auth.Key
			}
			localVarHeaderParams["X-API-Key"] = key
			
		}
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		if err == nil { 
			return localVarReturnValue, localVarHttpResponse, err
		}
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		
		if localVarHttpResponse.StatusCode == 200 {
			var v ExamplepbUnannotatedSimpleMessage
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		if localVarHttpResponse.StatusCode == 403 {
			var v interface{}
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		if localVarHttpResponse.StatusCode == 404 {
			var v int32
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		if localVarHttpResponse.StatusCode == 503 {
			var v interface{}
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		if localVarHttpResponse.StatusCode == 0 {
			var v RpcStatus
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}

/* 
UnannotatedEchoServiceApiService Summary: Echo rpc
Description Echo
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param id Id represents the message identifier.
 * @param num Int value field
 * @param optional nil or *UnannotatedEchoServiceEcho2Opts - Optional Parameters:
     * @param "Duration" (optional.String) - 
     * @param "LineNum" (optional.String) - 
     * @param "Lang" (optional.String) - 
     * @param "StatusProgress" (optional.String) - 
     * @param "StatusNote" (optional.String) - 
     * @param "En" (optional.String) - 
     * @param "NoProgress" (optional.String) - 
     * @param "NoNote" (optional.String) - 

@return ExamplepbUnannotatedSimpleMessage
*/

type UnannotatedEchoServiceEcho2Opts struct { 
	Duration optional.String
	LineNum optional.String
	Lang optional.String
	StatusProgress optional.String
	StatusNote optional.String
	En optional.String
	NoProgress optional.String
	NoNote optional.String
}

func (a *UnannotatedEchoServiceApiService) UnannotatedEchoServiceEcho2(ctx context.Context, id string, num string, localVarOptionals *UnannotatedEchoServiceEcho2Opts) (ExamplepbUnannotatedSimpleMessage, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue ExamplepbUnannotatedSimpleMessage
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/v1/example/echo/{id}/{num}"
	localVarPath = strings.Replace(localVarPath, "{"+"id"+"}", fmt.Sprintf("%v", id), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"num"+"}", fmt.Sprintf("%v", num), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarOptionals != nil && localVarOptionals.Duration.IsSet() {
		localVarQueryParams.Add("duration", parameterToString(localVarOptionals.Duration.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.LineNum.IsSet() {
		localVarQueryParams.Add("lineNum", parameterToString(localVarOptionals.LineNum.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Lang.IsSet() {
		localVarQueryParams.Add("lang", parameterToString(localVarOptionals.Lang.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.StatusProgress.IsSet() {
		localVarQueryParams.Add("status.progress", parameterToString(localVarOptionals.StatusProgress.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.StatusNote.IsSet() {
		localVarQueryParams.Add("status.note", parameterToString(localVarOptionals.StatusNote.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.En.IsSet() {
		localVarQueryParams.Add("en", parameterToString(localVarOptionals.En.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.NoProgress.IsSet() {
		localVarQueryParams.Add("no.progress", parameterToString(localVarOptionals.NoProgress.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.NoNote.IsSet() {
		localVarQueryParams.Add("no.note", parameterToString(localVarOptionals.NoNote.Value(), ""))
	}
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json", "application/x-foo-mime"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json", "application/x-foo-mime"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	if ctx != nil {
		// API Key Authentication
		if auth, ok := ctx.Value(ContextAPIKey).(APIKey); ok {
			var key string
			if auth.Prefix != "" {
				key = auth.Prefix + " " + auth.Key
			} else {
				key = auth.Key
			}
			localVarHeaderParams["X-API-Key"] = key
			
		}
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		if err == nil { 
			return localVarReturnValue, localVarHttpResponse, err
		}
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		
		if localVarHttpResponse.StatusCode == 200 {
			var v ExamplepbUnannotatedSimpleMessage
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		if localVarHttpResponse.StatusCode == 403 {
			var v interface{}
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		if localVarHttpResponse.StatusCode == 404 {
			var v int32
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		if localVarHttpResponse.StatusCode == 503 {
			var v interface{}
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		if localVarHttpResponse.StatusCode == 0 {
			var v RpcStatus
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}

/* 
UnannotatedEchoServiceApiService EchoBody method receives a simple message and returns it.
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param body

@return ExamplepbUnannotatedSimpleMessage
*/
func (a *UnannotatedEchoServiceApiService) UnannotatedEchoServiceEchoBody(ctx context.Context, body ExamplepbUnannotatedSimpleMessage) (ExamplepbUnannotatedSimpleMessage, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Post")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue ExamplepbUnannotatedSimpleMessage
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/v1/example/echo_body"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json", "application/x-foo-mime"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json", "application/x-foo-mime"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	// body params
	localVarPostBody = &body
	if ctx != nil {
		// API Key Authentication
		if auth, ok := ctx.Value(ContextAPIKey).(APIKey); ok {
			var key string
			if auth.Prefix != "" {
				key = auth.Prefix + " " + auth.Key
			} else {
				key = auth.Key
			}
			localVarHeaderParams["X-API-Key"] = key
			
		}
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		if err == nil { 
			return localVarReturnValue, localVarHttpResponse, err
		}
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		
		if localVarHttpResponse.StatusCode == 200 {
			var v ExamplepbUnannotatedSimpleMessage
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		if localVarHttpResponse.StatusCode == 403 {
			var v interface{}
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		if localVarHttpResponse.StatusCode == 404 {
			var v string
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		if localVarHttpResponse.StatusCode == 0 {
			var v RpcStatus
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}

/* 
UnannotatedEchoServiceApiService EchoDelete method receives a simple message and returns it.
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param num Int value field
 * @param optional nil or *UnannotatedEchoServiceEchoDeleteOpts - Optional Parameters:
     * @param "Id" (optional.String) -  Id represents the message identifier.
     * @param "Duration" (optional.String) - 
     * @param "LineNum" (optional.String) - 
     * @param "Lang" (optional.String) - 
     * @param "StatusProgress" (optional.String) - 
     * @param "StatusNote" (optional.String) - 
     * @param "En" (optional.String) - 
     * @param "NoProgress" (optional.String) - 
     * @param "NoNote" (optional.String) - 

@return ExamplepbUnannotatedSimpleMessage
*/

type UnannotatedEchoServiceEchoDeleteOpts struct { 
	Id optional.String
	Duration optional.String
	LineNum optional.String
	Lang optional.String
	StatusProgress optional.String
	StatusNote optional.String
	En optional.String
	NoProgress optional.String
	NoNote optional.String
}

func (a *UnannotatedEchoServiceApiService) UnannotatedEchoServiceEchoDelete(ctx context.Context, num string, localVarOptionals *UnannotatedEchoServiceEchoDeleteOpts) (ExamplepbUnannotatedSimpleMessage, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Delete")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue ExamplepbUnannotatedSimpleMessage
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/v1/example/echo_delete"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarOptionals != nil && localVarOptionals.Id.IsSet() {
		localVarQueryParams.Add("id", parameterToString(localVarOptionals.Id.Value(), ""))
	}
	localVarQueryParams.Add("num", parameterToString(num, ""))
	if localVarOptionals != nil && localVarOptionals.Duration.IsSet() {
		localVarQueryParams.Add("duration", parameterToString(localVarOptionals.Duration.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.LineNum.IsSet() {
		localVarQueryParams.Add("lineNum", parameterToString(localVarOptionals.LineNum.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Lang.IsSet() {
		localVarQueryParams.Add("lang", parameterToString(localVarOptionals.Lang.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.StatusProgress.IsSet() {
		localVarQueryParams.Add("status.progress", parameterToString(localVarOptionals.StatusProgress.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.StatusNote.IsSet() {
		localVarQueryParams.Add("status.note", parameterToString(localVarOptionals.StatusNote.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.En.IsSet() {
		localVarQueryParams.Add("en", parameterToString(localVarOptionals.En.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.NoProgress.IsSet() {
		localVarQueryParams.Add("no.progress", parameterToString(localVarOptionals.NoProgress.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.NoNote.IsSet() {
		localVarQueryParams.Add("no.note", parameterToString(localVarOptionals.NoNote.Value(), ""))
	}
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json", "application/x-foo-mime"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json", "application/x-foo-mime"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	if ctx != nil {
		// API Key Authentication
		if auth, ok := ctx.Value(ContextAPIKey).(APIKey); ok {
			var key string
			if auth.Prefix != "" {
				key = auth.Prefix + " " + auth.Key
			} else {
				key = auth.Key
			}
			localVarHeaderParams["X-API-Key"] = key
			
		}
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		if err == nil { 
			return localVarReturnValue, localVarHttpResponse, err
		}
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		
		if localVarHttpResponse.StatusCode == 200 {
			var v ExamplepbUnannotatedSimpleMessage
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		if localVarHttpResponse.StatusCode == 403 {
			var v interface{}
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		if localVarHttpResponse.StatusCode == 404 {
			var v string
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		if localVarHttpResponse.StatusCode == 0 {
			var v RpcStatus
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}
