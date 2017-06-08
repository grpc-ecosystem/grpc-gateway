# \ABitOfEverythingServiceApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Create**](ABitOfEverythingServiceApi.md#Create) | **Post** /v1/example/a_bit_of_everything/{float_value}/{double_value}/{int64_value}/separator/{uint64_value}/{int32_value}/{fixed64_value}/{fixed32_value}/{bool_value}/{string_value}/{uint32_value}/{sfixed32_value}/{sfixed64_value}/{sint32_value}/{sint64_value}/{nonConventionalNameValue} | 
[**CreateBody**](ABitOfEverythingServiceApi.md#CreateBody) | **Post** /v1/example/a_bit_of_everything | 
[**DeepPathEcho**](ABitOfEverythingServiceApi.md#DeepPathEcho) | **Post** /v1/example/a_bit_of_everything/{single_nested.name} | 
[**Delete**](ABitOfEverythingServiceApi.md#Delete) | **Delete** /v1/example/a_bit_of_everything/{uuid} | 
[**Echo**](ABitOfEverythingServiceApi.md#Echo) | **Get** /v1/example/a_bit_of_everything/echo/{value} | 
[**Echo_0**](ABitOfEverythingServiceApi.md#Echo_0) | **Get** /v2/example/echo | 
[**Echo_1**](ABitOfEverythingServiceApi.md#Echo_1) | **Post** /v2/example/echo | 
[**Lookup**](ABitOfEverythingServiceApi.md#Lookup) | **Get** /v1/example/a_bit_of_everything/{uuid} | 
[**Timeout**](ABitOfEverythingServiceApi.md#Timeout) | **Get** /v2/example/timeout | 
[**Update**](ABitOfEverythingServiceApi.md#Update) | **Put** /v1/example/a_bit_of_everything/{uuid} | 


# **Create**
> ExamplepbABitOfEverything Create($floatValue, $doubleValue, $int64Value, $uint64Value, $int32Value, $fixed64Value, $fixed32Value, $boolValue, $stringValue, $uint32Value, $sfixed32Value, $sfixed64Value, $sint32Value, $sint64Value, $nonConventionalNameValue)




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **floatValue** | **float32**|  | 
 **doubleValue** | **float64**|  | 
 **int64Value** | **string**|  | 
 **uint64Value** | **string**|  | 
 **int32Value** | **int32**|  | 
 **fixed64Value** | **string**|  | 
 **fixed32Value** | **int64**|  | 
 **boolValue** | **bool**|  | 
 **stringValue** | **string**|  | 
 **uint32Value** | **int64**|  | 
 **sfixed32Value** | **int32**|  | 
 **sfixed64Value** | **string**|  | 
 **sint32Value** | **int32**|  | 
 **sint64Value** | **string**|  | 
 **nonConventionalNameValue** | **string**|  | 

### Return type

[**ExamplepbABitOfEverything**](examplepbABitOfEverything.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateBody**
> ExamplepbABitOfEverything CreateBody($body)




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**ExamplepbABitOfEverything**](ExamplepbABitOfEverything.md)|  | 

### Return type

[**ExamplepbABitOfEverything**](examplepbABitOfEverything.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeepPathEcho**
> ExamplepbABitOfEverything DeepPathEcho($singleNestedName, $body)




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **singleNestedName** | **string**|  | 
 **body** | [**ExamplepbABitOfEverything**](ExamplepbABitOfEverything.md)|  | 

### Return type

[**ExamplepbABitOfEverything**](examplepbABitOfEverything.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **Delete**
> ProtobufEmpty Delete($uuid)




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **uuid** | **string**|  | 

### Return type

[**ProtobufEmpty**](protobufEmpty.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **Echo**
> SubStringMessage Echo($value)




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **value** | **string**|  | 

### Return type

[**SubStringMessage**](subStringMessage.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **Echo_0**
> SubStringMessage Echo_0()




### Parameters
This endpoint does not need any parameter.

### Return type

[**SubStringMessage**](subStringMessage.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **Echo_1**
> SubStringMessage Echo_1($body)




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | **string**|  | 

### Return type

[**SubStringMessage**](subStringMessage.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **Lookup**
> ExamplepbABitOfEverything Lookup($uuid)




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **uuid** | **string**|  | 

### Return type

[**ExamplepbABitOfEverything**](examplepbABitOfEverything.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **Timeout**
> ProtobufEmpty Timeout()




### Parameters
This endpoint does not need any parameter.

### Return type

[**ProtobufEmpty**](protobufEmpty.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **Update**
> ProtobufEmpty Update($uuid, $body)




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **uuid** | **string**|  | 
 **body** | [**ExamplepbABitOfEverything**](ExamplepbABitOfEverything.md)|  | 

### Return type

[**ProtobufEmpty**](protobufEmpty.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

