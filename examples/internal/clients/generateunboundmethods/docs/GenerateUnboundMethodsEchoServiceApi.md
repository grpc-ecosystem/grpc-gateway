# \GenerateUnboundMethodsEchoServiceApi

All URIs are relative to *https://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GenerateUnboundMethodsEchoServiceEcho**](GenerateUnboundMethodsEchoServiceApi.md#GenerateUnboundMethodsEchoServiceEcho) | **Post** /grpc.gateway.examples.internal.examplepb.GenerateUnboundMethodsEchoService/Echo | Echo method receives a simple message and returns it.
[**GenerateUnboundMethodsEchoServiceEchoBody**](GenerateUnboundMethodsEchoServiceApi.md#GenerateUnboundMethodsEchoServiceEchoBody) | **Post** /grpc.gateway.examples.internal.examplepb.GenerateUnboundMethodsEchoService/EchoBody | EchoBody method receives a simple message and returns it.
[**GenerateUnboundMethodsEchoServiceEchoDelete**](GenerateUnboundMethodsEchoServiceApi.md#GenerateUnboundMethodsEchoServiceEchoDelete) | **Post** /grpc.gateway.examples.internal.examplepb.GenerateUnboundMethodsEchoService/EchoDelete | EchoDelete method receives a simple message and returns it.


# **GenerateUnboundMethodsEchoServiceEcho**
> ExamplepbGenerateUnboundMethodsSimpleMessage GenerateUnboundMethodsEchoServiceEcho(ctx, body)
Echo method receives a simple message and returns it.

The message posted as the id parameter will also be returned.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**ExamplepbGenerateUnboundMethodsSimpleMessage**](ExamplepbGenerateUnboundMethodsSimpleMessage.md)|  | 

### Return type

[**ExamplepbGenerateUnboundMethodsSimpleMessage**](examplepbGenerateUnboundMethodsSimpleMessage.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GenerateUnboundMethodsEchoServiceEchoBody**
> ExamplepbGenerateUnboundMethodsSimpleMessage GenerateUnboundMethodsEchoServiceEchoBody(ctx, body)
EchoBody method receives a simple message and returns it.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**ExamplepbGenerateUnboundMethodsSimpleMessage**](ExamplepbGenerateUnboundMethodsSimpleMessage.md)|  | 

### Return type

[**ExamplepbGenerateUnboundMethodsSimpleMessage**](examplepbGenerateUnboundMethodsSimpleMessage.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GenerateUnboundMethodsEchoServiceEchoDelete**
> ExamplepbGenerateUnboundMethodsSimpleMessage GenerateUnboundMethodsEchoServiceEchoDelete(ctx, body)
EchoDelete method receives a simple message and returns it.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**ExamplepbGenerateUnboundMethodsSimpleMessage**](ExamplepbGenerateUnboundMethodsSimpleMessage.md)|  | 

### Return type

[**ExamplepbGenerateUnboundMethodsSimpleMessage**](examplepbGenerateUnboundMethodsSimpleMessage.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

