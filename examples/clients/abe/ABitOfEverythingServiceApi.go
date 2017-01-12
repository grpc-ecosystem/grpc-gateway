package abe

import (
    "strings"
    "fmt"
    "encoding/json"
    "errors"
    "github.com/dghubble/sling"
    "time"
)

type ABitOfEverythingServiceApi struct {
    basePath  string
}

func NewABitOfEverythingServiceApi() *ABitOfEverythingServiceApi{
    return &ABitOfEverythingServiceApi {
        basePath:   "http://localhost",
    }
}

func NewABitOfEverythingServiceApiWithBasePath(basePath string) *ABitOfEverythingServiceApi{
    return &ABitOfEverythingServiceApi {
        basePath:   basePath,
    }
}

/**
 * 
 * 
 * @param floatValue 
 * @param doubleValue 
 * @param int64Value 
 * @param uint64Value 
 * @param int32Value 
 * @param fixed64Value 
 * @param fixed32Value 
 * @param boolValue 
 * @param stringValue 
 * @param uint32Value 
 * @param sfixed32Value 
 * @param sfixed64Value 
 * @param sint32Value 
 * @param sint64Value 
 * @param nonConventionalNameValue 
 * @return ExamplepbABitOfEverything
 */
//func (a ABitOfEverythingServiceApi) Create (floatValue float32, doubleValue float64, int64Value string, uint64Value string, int32Value int32, fixed64Value string, fixed32Value int64, boolValue bool, stringValue string, uint32Value int64, sfixed32Value int32, sfixed64Value string, sint32Value int32, sint64Value string, nonConventionalNameValue string) (ExamplepbABitOfEverything, error) {
func (a ABitOfEverythingServiceApi) Create (floatValue float32, doubleValue float64, int64Value string, uint64Value string, int32Value int32, fixed64Value string, fixed32Value int64, boolValue bool, stringValue string, uint32Value int64, sfixed32Value int32, sfixed64Value string, sint32Value int32, sint64Value string, nonConventionalNameValue string) (ExamplepbABitOfEverything, error) {

    _sling := sling.New().Post(a.basePath)

    // create path and map variables
    path := "/v1/example/a_bit_of_everything/{float_value}/{double_value}/{int64_value}/separator/{uint64_value}/{int32_value}/{fixed64_value}/{fixed32_value}/{bool_value}/{string_value}/{uint32_value}/{sfixed32_value}/{sfixed64_value}/{sint32_value}/{sint64_value}/{nonConventionalNameValue}"
    path = strings.Replace(path, "{" + "float_value" + "}", fmt.Sprintf("%v", floatValue), -1)
    path = strings.Replace(path, "{" + "double_value" + "}", fmt.Sprintf("%v", doubleValue), -1)
    path = strings.Replace(path, "{" + "int64_value" + "}", fmt.Sprintf("%v", int64Value), -1)
    path = strings.Replace(path, "{" + "uint64_value" + "}", fmt.Sprintf("%v", uint64Value), -1)
    path = strings.Replace(path, "{" + "int32_value" + "}", fmt.Sprintf("%v", int32Value), -1)
    path = strings.Replace(path, "{" + "fixed64_value" + "}", fmt.Sprintf("%v", fixed64Value), -1)
    path = strings.Replace(path, "{" + "fixed32_value" + "}", fmt.Sprintf("%v", fixed32Value), -1)
    path = strings.Replace(path, "{" + "bool_value" + "}", fmt.Sprintf("%v", boolValue), -1)
    path = strings.Replace(path, "{" + "string_value" + "}", fmt.Sprintf("%v", stringValue), -1)
    path = strings.Replace(path, "{" + "uint32_value" + "}", fmt.Sprintf("%v", uint32Value), -1)
    path = strings.Replace(path, "{" + "sfixed32_value" + "}", fmt.Sprintf("%v", sfixed32Value), -1)
    path = strings.Replace(path, "{" + "sfixed64_value" + "}", fmt.Sprintf("%v", sfixed64Value), -1)
    path = strings.Replace(path, "{" + "sint32_value" + "}", fmt.Sprintf("%v", sint32Value), -1)
    path = strings.Replace(path, "{" + "sint64_value" + "}", fmt.Sprintf("%v", sint64Value), -1)
    path = strings.Replace(path, "{" + "nonConventionalNameValue" + "}", fmt.Sprintf("%v", nonConventionalNameValue), -1)

    _sling = _sling.Path(path)

    // accept header
    accepts := []string { "application/json" }
    for key := range accepts {
        _sling = _sling.Set("Accept", accepts[key])
        break // only use the first Accept
    }


  var successPayload = new(ExamplepbABitOfEverything)

  // We use this map (below) so that any arbitrary error JSON can be handled.
  // FIXME: This is in the absence of this Go generator honoring the non-2xx
  // response (error) models, which needs to be implemented at some point.
  var failurePayload map[string]interface{}

  httpResponse, err := _sling.Receive(successPayload, &failurePayload)

  if err == nil {
    // err == nil only means that there wasn't a sub-application-layer error (e.g. no network error)
    if failurePayload != nil {
      // If the failurePayload is present, there likely was some kind of non-2xx status
      // returned (and a JSON payload error present)
      var str []byte
      str, err = json.Marshal(failurePayload)
      if err == nil { // For safety, check for an error marshalling... probably superfluous
        // This will return the JSON error body as a string
        err = errors.New(string(str))
      }
  } else {
    // So, there was no network-type error, and nothing in the failure payload,
    // but we should still check the status code
    if httpResponse == nil {
      // This should never happen...
      err = errors.New("No HTTP Response received.")
    } else if code := httpResponse.StatusCode; 200 > code || code > 299 {
        err = errors.New("HTTP Error: " + string(httpResponse.StatusCode))
      }
    }
  }

  return *successPayload, err
}
/**
 * 
 * 
 * @param body 
 * @return ExamplepbABitOfEverything
 */
//func (a ABitOfEverythingServiceApi) CreateBody (body ExamplepbABitOfEverything) (ExamplepbABitOfEverything, error) {
func (a ABitOfEverythingServiceApi) CreateBody (body ExamplepbABitOfEverything) (ExamplepbABitOfEverything, error) {

    _sling := sling.New().Post(a.basePath)

    // create path and map variables
    path := "/v1/example/a_bit_of_everything"

    _sling = _sling.Path(path)

    // accept header
    accepts := []string { "application/json" }
    for key := range accepts {
        _sling = _sling.Set("Accept", accepts[key])
        break // only use the first Accept
    }

// body params
    _sling = _sling.BodyJSON(body)

  var successPayload = new(ExamplepbABitOfEverything)

  // We use this map (below) so that any arbitrary error JSON can be handled.
  // FIXME: This is in the absence of this Go generator honoring the non-2xx
  // response (error) models, which needs to be implemented at some point.
  var failurePayload map[string]interface{}

  httpResponse, err := _sling.Receive(successPayload, &failurePayload)

  if err == nil {
    // err == nil only means that there wasn't a sub-application-layer error (e.g. no network error)
    if failurePayload != nil {
      // If the failurePayload is present, there likely was some kind of non-2xx status
      // returned (and a JSON payload error present)
      var str []byte
      str, err = json.Marshal(failurePayload)
      if err == nil { // For safety, check for an error marshalling... probably superfluous
        // This will return the JSON error body as a string
        err = errors.New(string(str))
      }
  } else {
    // So, there was no network-type error, and nothing in the failure payload,
    // but we should still check the status code
    if httpResponse == nil {
      // This should never happen...
      err = errors.New("No HTTP Response received.")
    } else if code := httpResponse.StatusCode; 200 > code || code > 299 {
        err = errors.New("HTTP Error: " + string(httpResponse.StatusCode))
      }
    }
  }

  return *successPayload, err
}
/**
 * 
 * 
 * @param singleNestedName 
 * @param body 
 * @return ExamplepbABitOfEverything
 */
//func (a ABitOfEverythingServiceApi) DeepPathEcho (singleNestedName string, body ExamplepbABitOfEverything) (ExamplepbABitOfEverything, error) {
func (a ABitOfEverythingServiceApi) DeepPathEcho (singleNestedName string, body ExamplepbABitOfEverything) (ExamplepbABitOfEverything, error) {

    _sling := sling.New().Post(a.basePath)

    // create path and map variables
    path := "/v1/example/a_bit_of_everything/{single_nested.name}"
    path = strings.Replace(path, "{" + "single_nested.name" + "}", fmt.Sprintf("%v", singleNestedName), -1)

    _sling = _sling.Path(path)

    // accept header
    accepts := []string { "application/json" }
    for key := range accepts {
        _sling = _sling.Set("Accept", accepts[key])
        break // only use the first Accept
    }

// body params
    _sling = _sling.BodyJSON(body)

  var successPayload = new(ExamplepbABitOfEverything)

  // We use this map (below) so that any arbitrary error JSON can be handled.
  // FIXME: This is in the absence of this Go generator honoring the non-2xx
  // response (error) models, which needs to be implemented at some point.
  var failurePayload map[string]interface{}

  httpResponse, err := _sling.Receive(successPayload, &failurePayload)

  if err == nil {
    // err == nil only means that there wasn't a sub-application-layer error (e.g. no network error)
    if failurePayload != nil {
      // If the failurePayload is present, there likely was some kind of non-2xx status
      // returned (and a JSON payload error present)
      var str []byte
      str, err = json.Marshal(failurePayload)
      if err == nil { // For safety, check for an error marshalling... probably superfluous
        // This will return the JSON error body as a string
        err = errors.New(string(str))
      }
  } else {
    // So, there was no network-type error, and nothing in the failure payload,
    // but we should still check the status code
    if httpResponse == nil {
      // This should never happen...
      err = errors.New("No HTTP Response received.")
    } else if code := httpResponse.StatusCode; 200 > code || code > 299 {
        err = errors.New("HTTP Error: " + string(httpResponse.StatusCode))
      }
    }
  }

  return *successPayload, err
}
/**
 * 
 * 
 * @param uuid 
 * @return ProtobufEmpty
 */
//func (a ABitOfEverythingServiceApi) Delete (uuid string) (ProtobufEmpty, error) {
func (a ABitOfEverythingServiceApi) Delete (uuid string) (ProtobufEmpty, error) {

    _sling := sling.New().Delete(a.basePath)

    // create path and map variables
    path := "/v1/example/a_bit_of_everything/{uuid}"
    path = strings.Replace(path, "{" + "uuid" + "}", fmt.Sprintf("%v", uuid), -1)

    _sling = _sling.Path(path)

    // accept header
    accepts := []string { "application/json" }
    for key := range accepts {
        _sling = _sling.Set("Accept", accepts[key])
        break // only use the first Accept
    }


  var successPayload = new(ProtobufEmpty)

  // We use this map (below) so that any arbitrary error JSON can be handled.
  // FIXME: This is in the absence of this Go generator honoring the non-2xx
  // response (error) models, which needs to be implemented at some point.
  var failurePayload map[string]interface{}

  httpResponse, err := _sling.Receive(successPayload, &failurePayload)

  if err == nil {
    // err == nil only means that there wasn't a sub-application-layer error (e.g. no network error)
    if failurePayload != nil {
      // If the failurePayload is present, there likely was some kind of non-2xx status
      // returned (and a JSON payload error present)
      var str []byte
      str, err = json.Marshal(failurePayload)
      if err == nil { // For safety, check for an error marshalling... probably superfluous
        // This will return the JSON error body as a string
        err = errors.New(string(str))
      }
  } else {
    // So, there was no network-type error, and nothing in the failure payload,
    // but we should still check the status code
    if httpResponse == nil {
      // This should never happen...
      err = errors.New("No HTTP Response received.")
    } else if code := httpResponse.StatusCode; 200 > code || code > 299 {
        err = errors.New("HTTP Error: " + string(httpResponse.StatusCode))
      }
    }
  }

  return *successPayload, err
}
/**
 * 
 * 
 * @param value 
 * @return SubStringMessage
 */
//func (a ABitOfEverythingServiceApi) Echo (value string) (SubStringMessage, error) {
func (a ABitOfEverythingServiceApi) Echo (value string) (SubStringMessage, error) {

    _sling := sling.New().Get(a.basePath)

    // create path and map variables
    path := "/v1/example/a_bit_of_everything/echo/{value}"
    path = strings.Replace(path, "{" + "value" + "}", fmt.Sprintf("%v", value), -1)

    _sling = _sling.Path(path)

    // accept header
    accepts := []string { "application/json" }
    for key := range accepts {
        _sling = _sling.Set("Accept", accepts[key])
        break // only use the first Accept
    }


  var successPayload = new(SubStringMessage)

  // We use this map (below) so that any arbitrary error JSON can be handled.
  // FIXME: This is in the absence of this Go generator honoring the non-2xx
  // response (error) models, which needs to be implemented at some point.
  var failurePayload map[string]interface{}

  httpResponse, err := _sling.Receive(successPayload, &failurePayload)

  if err == nil {
    // err == nil only means that there wasn't a sub-application-layer error (e.g. no network error)
    if failurePayload != nil {
      // If the failurePayload is present, there likely was some kind of non-2xx status
      // returned (and a JSON payload error present)
      var str []byte
      str, err = json.Marshal(failurePayload)
      if err == nil { // For safety, check for an error marshalling... probably superfluous
        // This will return the JSON error body as a string
        err = errors.New(string(str))
      }
  } else {
    // So, there was no network-type error, and nothing in the failure payload,
    // but we should still check the status code
    if httpResponse == nil {
      // This should never happen...
      err = errors.New("No HTTP Response received.")
    } else if code := httpResponse.StatusCode; 200 > code || code > 299 {
        err = errors.New("HTTP Error: " + string(httpResponse.StatusCode))
      }
    }
  }

  return *successPayload, err
}
/**
 * 
 * 
 * @param value 
 * @return SubStringMessage
 */
//func (a ABitOfEverythingServiceApi) Echo_1 (value string) (SubStringMessage, error) {
func (a ABitOfEverythingServiceApi) Echo_1 (value string) (SubStringMessage, error) {

    _sling := sling.New().Get(a.basePath)

    // create path and map variables
    path := "/v2/example/echo"

    _sling = _sling.Path(path)

    type QueryParams struct {
        value    string `url:"value,omitempty"`
        
}
    _sling = _sling.QueryStruct(&QueryParams{ value: value })
    // accept header
    accepts := []string { "application/json" }
    for key := range accepts {
        _sling = _sling.Set("Accept", accepts[key])
        break // only use the first Accept
    }


  var successPayload = new(SubStringMessage)

  // We use this map (below) so that any arbitrary error JSON can be handled.
  // FIXME: This is in the absence of this Go generator honoring the non-2xx
  // response (error) models, which needs to be implemented at some point.
  var failurePayload map[string]interface{}

  httpResponse, err := _sling.Receive(successPayload, &failurePayload)

  if err == nil {
    // err == nil only means that there wasn't a sub-application-layer error (e.g. no network error)
    if failurePayload != nil {
      // If the failurePayload is present, there likely was some kind of non-2xx status
      // returned (and a JSON payload error present)
      var str []byte
      str, err = json.Marshal(failurePayload)
      if err == nil { // For safety, check for an error marshalling... probably superfluous
        // This will return the JSON error body as a string
        err = errors.New(string(str))
      }
  } else {
    // So, there was no network-type error, and nothing in the failure payload,
    // but we should still check the status code
    if httpResponse == nil {
      // This should never happen...
      err = errors.New("No HTTP Response received.")
    } else if code := httpResponse.StatusCode; 200 > code || code > 299 {
        err = errors.New("HTTP Error: " + string(httpResponse.StatusCode))
      }
    }
  }

  return *successPayload, err
}
/**
 * 
 * 
 * @param body 
 * @return SubStringMessage
 */
//func (a ABitOfEverythingServiceApi) Echo_2 (body string) (SubStringMessage, error) {
func (a ABitOfEverythingServiceApi) Echo_2 (body string) (SubStringMessage, error) {

    _sling := sling.New().Post(a.basePath)

    // create path and map variables
    path := "/v2/example/echo"

    _sling = _sling.Path(path)

    // accept header
    accepts := []string { "application/json" }
    for key := range accepts {
        _sling = _sling.Set("Accept", accepts[key])
        break // only use the first Accept
    }

// body params
    _sling = _sling.BodyJSON(body)

  var successPayload = new(SubStringMessage)

  // We use this map (below) so that any arbitrary error JSON can be handled.
  // FIXME: This is in the absence of this Go generator honoring the non-2xx
  // response (error) models, which needs to be implemented at some point.
  var failurePayload map[string]interface{}

  httpResponse, err := _sling.Receive(successPayload, &failurePayload)

  if err == nil {
    // err == nil only means that there wasn't a sub-application-layer error (e.g. no network error)
    if failurePayload != nil {
      // If the failurePayload is present, there likely was some kind of non-2xx status
      // returned (and a JSON payload error present)
      var str []byte
      str, err = json.Marshal(failurePayload)
      if err == nil { // For safety, check for an error marshalling... probably superfluous
        // This will return the JSON error body as a string
        err = errors.New(string(str))
      }
  } else {
    // So, there was no network-type error, and nothing in the failure payload,
    // but we should still check the status code
    if httpResponse == nil {
      // This should never happen...
      err = errors.New("No HTTP Response received.")
    } else if code := httpResponse.StatusCode; 200 > code || code > 299 {
        err = errors.New("HTTP Error: " + string(httpResponse.StatusCode))
      }
    }
  }

  return *successPayload, err
}
/**
 * 
 * 
 * @param uuid 
 * @param singleNestedName name is nested field.
 * @param singleNestedAmount 
 * @param singleNestedOk  - FALSE: FALSE is false.\n - TRUE: TRUE is true.
 * @param floatValue 
 * @param doubleValue 
 * @param int64Value 
 * @param uint64Value 
 * @param int32Value 
 * @param fixed64Value 
 * @param fixed32Value 
 * @param boolValue 
 * @param stringValue 
 * @param uint32Value TODO(yugui) add bytes_value.
 * @param enumValue  - ZERO: ZERO means 0\n - ONE: ONE means 1
 * @param sfixed32Value 
 * @param sfixed64Value 
 * @param sint32Value 
 * @param sint64Value 
 * @param repeatedStringValue 
 * @param oneofString 
 * @param nonConventionalNameValue 
 * @param timestampValue 
 * @param repeatedEnumValue repeated enum value. it is comma-separated in query.\n\n - ZERO: ZERO means 0\n - ONE: ONE means 1
 * @return ProtobufEmpty
 */
//func (a ABitOfEverythingServiceApi) GetQuery (uuid string, singleNestedName string, singleNestedAmount int64, singleNestedOk string, floatValue float32, doubleValue float64, int64Value string, uint64Value string, int32Value int32, fixed64Value string, fixed32Value int64, boolValue bool, stringValue string, uint32Value int64, enumValue string, sfixed32Value int32, sfixed64Value string, sint32Value int32, sint64Value string, repeatedStringValue []string, oneofString string, nonConventionalNameValue string, timestampValue time.Time, repeatedEnumValue []string) (ProtobufEmpty, error) {
func (a ABitOfEverythingServiceApi) GetQuery (uuid string, singleNestedName string, singleNestedAmount int64, singleNestedOk string, floatValue float32, doubleValue float64, int64Value string, uint64Value string, int32Value int32, fixed64Value string, fixed32Value int64, boolValue bool, stringValue string, uint32Value int64, enumValue string, sfixed32Value int32, sfixed64Value string, sint32Value int32, sint64Value string, repeatedStringValue []string, oneofString string, nonConventionalNameValue string, timestampValue time.Time, repeatedEnumValue []string) (ProtobufEmpty, error) {

    _sling := sling.New().Get(a.basePath)

    // create path and map variables
    path := "/v1/example/a_bit_of_everything/query/{uuid}"
    path = strings.Replace(path, "{" + "uuid" + "}", fmt.Sprintf("%v", uuid), -1)

    _sling = _sling.Path(path)

    type QueryParams struct {
        singleNestedName    string `url:"single_nested.name,omitempty"`
        singleNestedAmount    int64 `url:"single_nested.amount,omitempty"`
        singleNestedOk    string `url:"single_nested.ok,omitempty"`
        floatValue    float32 `url:"float_value,omitempty"`
        doubleValue    float64 `url:"double_value,omitempty"`
        int64Value    string `url:"int64_value,omitempty"`
        uint64Value    string `url:"uint64_value,omitempty"`
        int32Value    int32 `url:"int32_value,omitempty"`
        fixed64Value    string `url:"fixed64_value,omitempty"`
        fixed32Value    int64 `url:"fixed32_value,omitempty"`
        boolValue    bool `url:"bool_value,omitempty"`
        stringValue    string `url:"string_value,omitempty"`
        uint32Value    int64 `url:"uint32_value,omitempty"`
        enumValue    string `url:"enum_value,omitempty"`
        sfixed32Value    int32 `url:"sfixed32_value,omitempty"`
        sfixed64Value    string `url:"sfixed64_value,omitempty"`
        sint32Value    int32 `url:"sint32_value,omitempty"`
        sint64Value    string `url:"sint64_value,omitempty"`
        repeatedStringValue    []string `url:"repeated_string_value,omitempty"`
        oneofString    string `url:"oneof_string,omitempty"`
        nonConventionalNameValue    string `url:"nonConventionalNameValue,omitempty"`
        timestampValue    time.Time `url:"timestamp_value,omitempty"`
        repeatedEnumValue    []string `url:"repeated_enum_value,omitempty"`
        
}
    _sling = _sling.QueryStruct(&QueryParams{ singleNestedName: singleNestedName,singleNestedAmount: singleNestedAmount,singleNestedOk: singleNestedOk,floatValue: floatValue,doubleValue: doubleValue,int64Value: int64Value,uint64Value: uint64Value,int32Value: int32Value,fixed64Value: fixed64Value,fixed32Value: fixed32Value,boolValue: boolValue,stringValue: stringValue,uint32Value: uint32Value,enumValue: enumValue,sfixed32Value: sfixed32Value,sfixed64Value: sfixed64Value,sint32Value: sint32Value,sint64Value: sint64Value,repeatedStringValue: repeatedStringValue,oneofString: oneofString,nonConventionalNameValue: nonConventionalNameValue,timestampValue: timestampValue,repeatedEnumValue: repeatedEnumValue })
    // accept header
    accepts := []string { "application/json" }
    for key := range accepts {
        _sling = _sling.Set("Accept", accepts[key])
        break // only use the first Accept
    }


  var successPayload = new(ProtobufEmpty)

  // We use this map (below) so that any arbitrary error JSON can be handled.
  // FIXME: This is in the absence of this Go generator honoring the non-2xx
  // response (error) models, which needs to be implemented at some point.
  var failurePayload map[string]interface{}

  httpResponse, err := _sling.Receive(successPayload, &failurePayload)

  if err == nil {
    // err == nil only means that there wasn't a sub-application-layer error (e.g. no network error)
    if failurePayload != nil {
      // If the failurePayload is present, there likely was some kind of non-2xx status
      // returned (and a JSON payload error present)
      var str []byte
      str, err = json.Marshal(failurePayload)
      if err == nil { // For safety, check for an error marshalling... probably superfluous
        // This will return the JSON error body as a string
        err = errors.New(string(str))
      }
  } else {
    // So, there was no network-type error, and nothing in the failure payload,
    // but we should still check the status code
    if httpResponse == nil {
      // This should never happen...
      err = errors.New("No HTTP Response received.")
    } else if code := httpResponse.StatusCode; 200 > code || code > 299 {
        err = errors.New("HTTP Error: " + string(httpResponse.StatusCode))
      }
    }
  }

  return *successPayload, err
}
/**
 * 
 * 
 * @param uuid 
 * @return ExamplepbABitOfEverything
 */
//func (a ABitOfEverythingServiceApi) Lookup (uuid string) (ExamplepbABitOfEverything, error) {
func (a ABitOfEverythingServiceApi) Lookup (uuid string) (ExamplepbABitOfEverything, error) {

    _sling := sling.New().Get(a.basePath)

    // create path and map variables
    path := "/v1/example/a_bit_of_everything/{uuid}"
    path = strings.Replace(path, "{" + "uuid" + "}", fmt.Sprintf("%v", uuid), -1)

    _sling = _sling.Path(path)

    // accept header
    accepts := []string { "application/json" }
    for key := range accepts {
        _sling = _sling.Set("Accept", accepts[key])
        break // only use the first Accept
    }


  var successPayload = new(ExamplepbABitOfEverything)

  // We use this map (below) so that any arbitrary error JSON can be handled.
  // FIXME: This is in the absence of this Go generator honoring the non-2xx
  // response (error) models, which needs to be implemented at some point.
  var failurePayload map[string]interface{}

  httpResponse, err := _sling.Receive(successPayload, &failurePayload)

  if err == nil {
    // err == nil only means that there wasn't a sub-application-layer error (e.g. no network error)
    if failurePayload != nil {
      // If the failurePayload is present, there likely was some kind of non-2xx status
      // returned (and a JSON payload error present)
      var str []byte
      str, err = json.Marshal(failurePayload)
      if err == nil { // For safety, check for an error marshalling... probably superfluous
        // This will return the JSON error body as a string
        err = errors.New(string(str))
      }
  } else {
    // So, there was no network-type error, and nothing in the failure payload,
    // but we should still check the status code
    if httpResponse == nil {
      // This should never happen...
      err = errors.New("No HTTP Response received.")
    } else if code := httpResponse.StatusCode; 200 > code || code > 299 {
        err = errors.New("HTTP Error: " + string(httpResponse.StatusCode))
      }
    }
  }

  return *successPayload, err
}
/**
 * 
 * 
 * @return ProtobufEmpty
 */
//func (a ABitOfEverythingServiceApi) Timeout () (ProtobufEmpty, error) {
func (a ABitOfEverythingServiceApi) Timeout () (ProtobufEmpty, error) {

    _sling := sling.New().Get(a.basePath)

    // create path and map variables
    path := "/v2/example/timeout"

    _sling = _sling.Path(path)

    // accept header
    accepts := []string { "application/json" }
    for key := range accepts {
        _sling = _sling.Set("Accept", accepts[key])
        break // only use the first Accept
    }


  var successPayload = new(ProtobufEmpty)

  // We use this map (below) so that any arbitrary error JSON can be handled.
  // FIXME: This is in the absence of this Go generator honoring the non-2xx
  // response (error) models, which needs to be implemented at some point.
  var failurePayload map[string]interface{}

  httpResponse, err := _sling.Receive(successPayload, &failurePayload)

  if err == nil {
    // err == nil only means that there wasn't a sub-application-layer error (e.g. no network error)
    if failurePayload != nil {
      // If the failurePayload is present, there likely was some kind of non-2xx status
      // returned (and a JSON payload error present)
      var str []byte
      str, err = json.Marshal(failurePayload)
      if err == nil { // For safety, check for an error marshalling... probably superfluous
        // This will return the JSON error body as a string
        err = errors.New(string(str))
      }
  } else {
    // So, there was no network-type error, and nothing in the failure payload,
    // but we should still check the status code
    if httpResponse == nil {
      // This should never happen...
      err = errors.New("No HTTP Response received.")
    } else if code := httpResponse.StatusCode; 200 > code || code > 299 {
        err = errors.New("HTTP Error: " + string(httpResponse.StatusCode))
      }
    }
  }

  return *successPayload, err
}
/**
 * 
 * 
 * @param uuid 
 * @param body 
 * @return ProtobufEmpty
 */
//func (a ABitOfEverythingServiceApi) Update (uuid string, body ExamplepbABitOfEverything) (ProtobufEmpty, error) {
func (a ABitOfEverythingServiceApi) Update (uuid string, body ExamplepbABitOfEverything) (ProtobufEmpty, error) {

    _sling := sling.New().Put(a.basePath)

    // create path and map variables
    path := "/v1/example/a_bit_of_everything/{uuid}"
    path = strings.Replace(path, "{" + "uuid" + "}", fmt.Sprintf("%v", uuid), -1)

    _sling = _sling.Path(path)

    // accept header
    accepts := []string { "application/json" }
    for key := range accepts {
        _sling = _sling.Set("Accept", accepts[key])
        break // only use the first Accept
    }

// body params
    _sling = _sling.BodyJSON(body)

  var successPayload = new(ProtobufEmpty)

  // We use this map (below) so that any arbitrary error JSON can be handled.
  // FIXME: This is in the absence of this Go generator honoring the non-2xx
  // response (error) models, which needs to be implemented at some point.
  var failurePayload map[string]interface{}

  httpResponse, err := _sling.Receive(successPayload, &failurePayload)

  if err == nil {
    // err == nil only means that there wasn't a sub-application-layer error (e.g. no network error)
    if failurePayload != nil {
      // If the failurePayload is present, there likely was some kind of non-2xx status
      // returned (and a JSON payload error present)
      var str []byte
      str, err = json.Marshal(failurePayload)
      if err == nil { // For safety, check for an error marshalling... probably superfluous
        // This will return the JSON error body as a string
        err = errors.New(string(str))
      }
  } else {
    // So, there was no network-type error, and nothing in the failure payload,
    // but we should still check the status code
    if httpResponse == nil {
      // This should never happen...
      err = errors.New("No HTTP Response received.")
    } else if code := httpResponse.StatusCode; 200 > code || code > 299 {
        err = errors.New("HTTP Error: " + string(httpResponse.StatusCode))
      }
    }
  }

  return *successPayload, err
}
