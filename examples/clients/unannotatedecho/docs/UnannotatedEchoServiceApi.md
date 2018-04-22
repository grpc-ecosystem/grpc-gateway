# \UnannotatedEchoServiceApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Echo**](UnannotatedEchoServiceApi.md#Echo) | **Post** /v1/example/echo/{id} | Echo method receives a simple message and returns it.
[**Echo2**](UnannotatedEchoServiceApi.md#Echo2) | **Get** /v1/example/echo/{id}/{num} | Echo method receives a simple message and returns it.
[**EchoBody**](UnannotatedEchoServiceApi.md#EchoBody) | **Post** /v1/example/echo_body | EchoBody method receives a simple message and returns it.


# **Echo**
> ExamplepbUnannotatedSimpleMessage Echo($id)

Echo method receives a simple message and returns it.

The message posted as the id parameter will also be returned.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**|  | 

### Return type

[**ExamplepbUnannotatedSimpleMessage**](examplepbUnannotatedSimpleMessage.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **Echo2**
> ExamplepbUnannotatedSimpleMessage Echo2($id, $num, $duration)

Echo method receives a simple message and returns it.

The message posted as the id parameter will also be returned.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**|  | 
 **num** | **string**|  | 
 **duration** | **string**|  | [optional] 

### Return type

[**ExamplepbUnannotatedSimpleMessage**](examplepbUnannotatedSimpleMessage.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **EchoBody**
> ExamplepbUnannotatedSimpleMessage EchoBody($body)

EchoBody method receives a simple message and returns it.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**ExamplepbUnannotatedSimpleMessage**](ExamplepbUnannotatedSimpleMessage.md)|  | 

### Return type

[**ExamplepbUnannotatedSimpleMessage**](examplepbUnannotatedSimpleMessage.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

