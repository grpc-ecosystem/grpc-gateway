# \EchoServiceApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Echo**](EchoServiceApi.md#Echo) | **Post** /v1/example/echo/{id} | Echo method receives a simple message and returns it.
[**EchoBody**](EchoServiceApi.md#EchoBody) | **Post** /v1/example/echo_body | EchoBody method receives a simple message and returns it.


# **Echo**
> ExamplepbSimpleMessage Echo($id)

Echo method receives a simple message and returns it.

The message posted as the id parameter will also be returned.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**|  | 

### Return type

[**ExamplepbSimpleMessage**](examplepbSimpleMessage.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **EchoBody**
> ExamplepbSimpleMessage EchoBody($body)

EchoBody method receives a simple message and returns it.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**ExamplepbSimpleMessage**](ExamplepbSimpleMessage.md)|  | 

### Return type

[**ExamplepbSimpleMessage**](examplepbSimpleMessage.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

