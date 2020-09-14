package descriptor

import (
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
)

func TestLoadGrpcAPIServiceFromYAMLInvalidType(t *testing.T) {
	// Ideally this would fail but for now this test documents that it doesn't
	service, err := loadGrpcAPIServiceFromYAML([]byte(`type: not.the.right.type`), "invalidtype")
	if err != nil {
		t.Fatal(err)
	}

	if service == nil {
		t.Fatal("No service returned")
	}
}

func TestLoadGrpcAPIServiceFromYAMLSingleRule(t *testing.T) {
	service, err := loadGrpcAPIServiceFromYAML([]byte(`
type: google.api.Service
config_version: 3

http:
 rules:
 - selector: grpctest.YourService.Echo
   post: /v1/myecho
   body: "*"
`), "example")
	if err != nil {
		t.Fatal(err)
	}

	if service.Http == nil {
		t.Fatal("HTTP is empty")
	}

	if len(service.Http.GetRules()) != 1 {
		t.Fatalf("Have %v rules instead of one. Got: %v", len(service.Http.GetRules()), service.Http.GetRules())
	}

	rule := service.Http.GetRules()[0]
	if rule.GetSelector() != "grpctest.YourService.Echo" {
		t.Errorf("Rule has unexpected selector '%v'", rule.GetSelector())
	}
	if rule.GetPost() != "/v1/myecho" {
		t.Errorf("Rule has unexpected post '%v'", rule.GetPost())
	}
	if rule.GetBody() != "*" {
		t.Errorf("Rule has unexpected body '%v'", rule.GetBody())
	}
}

func TestLoadGrpcAPIServiceFromYAMLRejectInvalidYAML(t *testing.T) {
	service, err := loadGrpcAPIServiceFromYAML([]byte(`
type: google.api.Service
config_version: 3

http:
 rules:
 - selector: grpctest.YourService.Echo
   - post: thislinebreakstheselectorblockabovewiththeleadingdash
   body: "*"
`), "invalidyaml")
	if err == nil {
		t.Fatal(err)
	}

	if !strings.Contains(err.Error(), "line 7") {
		t.Errorf("Expected yaml error to be detected in line 7. Got other error: %v", err)
	}

	if service != nil {
		t.Fatal("Service returned")
	}
}

func TestLoadGrpcAPIServiceFromYAMLMultipleWithAdditionalBindings(t *testing.T) {
	service, err := loadGrpcAPIServiceFromYAML([]byte(`
type: google.api.Service
config_version: 3

http:
 rules:
 - selector: first.selector
   post: /my/post/path
   body: "*"
   additional_bindings:
   - post: /additional/post/path
   - put: /additional/put/{value}/path
   - delete: "{value}"
   - patch: "/additional/patch/{value}"
 - selector: some.other.service
   delete: foo
`), "example")
	if err != nil {
		t.Fatalf("Failed to load service description from YAML: %v", err)
	}

	if service == nil {
		t.Fatal("No service returned")
	}

	if service.Http == nil {
		t.Fatal("HTTP is empty")
	}

	if len(service.Http.GetRules()) != 2 {
		t.Fatalf("%v service(s) returned when two were expected. Got: %v", len(service.Http.GetRules()), service.Http)
	}

	first := service.Http.GetRules()[0]
	if first.GetSelector() != "first.selector" {
		t.Errorf("first.selector has unexpected selector '%v'", first.GetSelector())
	}
	if first.GetBody() != "*" {
		t.Errorf("first.selector has unexpected body '%v'", first.GetBody())
	}
	if first.GetPost() != "/my/post/path" {
		t.Errorf("first.selector has unexpected post '%v'", first.GetPost())
	}
	if len(first.GetAdditionalBindings()) != 4 {
		t.Fatalf("first.selector has unexpected number of bindings %v instead of four. Got: %v", len(first.GetAdditionalBindings()), first.GetAdditionalBindings())
	}
	if first.GetAdditionalBindings()[0].GetPost() != "/additional/post/path" {
		t.Errorf("first.selector additional binding 0 has unexpected post '%v'", first.GetAdditionalBindings()[0].GetPost())
	}
	if first.GetAdditionalBindings()[1].GetPut() != "/additional/put/{value}/path" {
		t.Errorf("first.selector additional binding 1 has unexpected put '%v'", first.GetAdditionalBindings()[0].GetPost())
	}
	if first.GetAdditionalBindings()[2].GetDelete() != "{value}" {
		t.Errorf("first.selector additional binding 2 has unexpected delete '%v'", first.GetAdditionalBindings()[0].GetPost())
	}
	if first.GetAdditionalBindings()[3].GetPatch() != "/additional/patch/{value}" {
		t.Errorf("first.selector additional binding 3 has unexpected patch '%v'", first.GetAdditionalBindings()[0].GetPost())
	}

	second := service.Http.GetRules()[1]
	if second.GetSelector() != "some.other.service" {
		t.Errorf("some.other.service has unexpected selector '%v'", second.GetSelector())
	}
	if second.GetDelete() != "foo" {
		t.Errorf("some.other.service has unexpected delete '%v'", second.GetDelete())
	}
	if len(second.GetAdditionalBindings()) != 0 {
		t.Errorf("some.other.service has %v additional bindings when it should not have any. Got: %v", len(second.GetAdditionalBindings()), second.GetAdditionalBindings())
	}
}

func TestLoadOpenAPIConfigFromYAMLRejectInvalidYAML(t *testing.T) {
	config, err := loadOpenAPIConfigFromYAML([]byte(`
openapiOptions:
file:
- file: test.proto
  - option:
      schemes:
        - HTTP
        - HTTPS
        - WSS
      securityDefinitions:
        security:
          ApiKeyAuth:
            type: TYPE_API_KEY
            in: IN_HEADER
            name: "X-API-Key"
`), "invalidyaml")
	if err == nil {
		t.Fatal(err)
	}

	if !strings.Contains(err.Error(), "line 4") {
		t.Errorf("Expected yaml error to be detected in line 4. Got other error: %v", err)
	}

	if config != nil {
		t.Fatal("Config returned")
	}
}

func TestLoadOpenAPIConfigFromYAML(t *testing.T) {
	config, err := loadOpenAPIConfigFromYAML([]byte(`
openapiOptions:
  file:
  - file: test.proto
    option:
      schemes:
      - HTTP
      - HTTPS
      - WSS
      securityDefinitions:
        security:
          ApiKeyAuth:
            type: TYPE_API_KEY
            in: IN_HEADER
            name: "X-API-Key"
`), "openapi_options")
	if err != nil {
		t.Fatal(err)
	}

	if config.OpenapiOptions == nil {
		t.Fatal("OpenAPIOptions is empty")
	}

	opts := config.OpenapiOptions
	if numFileOpts := len(opts.File); numFileOpts != 1 {
		t.Fatalf("expected 1 file option but got %d", numFileOpts)
	}

	fileOpt := opts.File[0]

	if fileOpt.File != "test.proto" {
		t.Fatalf("file option has unexpected binding %s", fileOpt.File)
	}

	swaggerOpt := fileOpt.Option

	if swaggerOpt == nil {
		t.Fatal("expected option to be set")
	}

	if numSchemes := len(swaggerOpt.Schemes); numSchemes != 3 {
		t.Fatalf("expected 3 schemes but got %d", numSchemes)
	}
	if swaggerOpt.Schemes[0] != options.Scheme_HTTP {
		t.Fatalf("expected first scheme to be HTTP but got %s", swaggerOpt.Schemes[0])
	}
	if swaggerOpt.Schemes[1] != options.Scheme_HTTPS {
		t.Fatalf("expected second scheme to be HTTPS but got %s", swaggerOpt.Schemes[1])
	}
	if swaggerOpt.Schemes[2] != options.Scheme_WSS {
		t.Fatalf("expected third scheme to be WSS but got %s", swaggerOpt.Schemes[2])
	}

	if swaggerOpt.SecurityDefinitions == nil {
		t.Fatal("expected securityDefinitions to be set")
	}
	if numSecOpts := len(swaggerOpt.SecurityDefinitions.Security); numSecOpts != 1 {
		t.Fatalf("expected 1 security option but got %d", numSecOpts)
	}
	secOpt, ok := swaggerOpt.SecurityDefinitions.Security["ApiKeyAuth"]
	if !ok {
		t.Fatal("no SecurityScheme for key \"ApiKeyAuth\"")
	}
	if secOpt.Type != options.SecurityScheme_TYPE_API_KEY {
		t.Fatalf("expected scheme type to be TYPE_API_KEY but got %s", secOpt.Type)
	}
	if secOpt.In != options.SecurityScheme_IN_HEADER {
		t.Fatalf("expected scheme  in to be IN_HEADER but got %s", secOpt.In)
	}
	if secOpt.Name != "X-API-Key" {
		t.Fatalf("expected name to be X-API-Key but got %s", secOpt.Name)
	}
}
