# \ResponseBodyServiceApi

All URIs are relative to *https://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ResponseBodyServiceGetResponseBodyStream**](ResponseBodyServiceApi.md#ResponseBodyServiceGetResponseBodyStream) | **Get** /responsebody/{data} | 
[**ResponseBodyServiceListResponseBodiesStream**](ResponseBodyServiceApi.md#ResponseBodyServiceListResponseBodiesStream) | **Get** /responsebodies/{data} | 
[**ResponseBodyServiceListResponseStringsStream**](ResponseBodyServiceApi.md#ResponseBodyServiceListResponseStringsStream) | **Get** /responsestrings/{data} | 


# **ResponseBodyServiceGetResponseBodyStream**
> StreamResultOfExamplepbResponseBodyOut ResponseBodyServiceGetResponseBodyStream(ctx, data)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **data** | **string**|  | 

### Return type

[**StreamResultOfExamplepbResponseBodyOut**](Stream result of examplepbResponseBodyOut.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ResponseBodyServiceListResponseBodiesStream**
> StreamResultOfExamplepbRepeatedResponseBodyOut ResponseBodyServiceListResponseBodiesStream(ctx, data)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **data** | **string**|  | 

### Return type

[**StreamResultOfExamplepbRepeatedResponseBodyOut**](Stream result of examplepbRepeatedResponseBodyOut.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ResponseBodyServiceListResponseStringsStream**
> StreamResultOfExamplepbRepeatedResponseBodyOut ResponseBodyServiceListResponseStringsStream(ctx, data)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **data** | **string**|  | 

### Return type

[**StreamResultOfExamplepbRepeatedResponseBodyOut**](Stream result of examplepbRepeatedResponseBodyOut.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

