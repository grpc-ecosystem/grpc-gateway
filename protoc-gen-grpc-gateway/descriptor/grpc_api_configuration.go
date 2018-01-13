package descriptor

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/api/annotations"
)

// GrpcAPIService represents a stripped down version of google.api.Service .
// Compare to https://github.com/googleapis/googleapis/blob/master/google/api/service.proto
//
// Note that for the purposes of the gateway generator we only consider a subset of all
// available features google supports in their service descriptions and hence we do not
// bother having the full service protobuf which has a lot of other dependencies. Thanks
// to backwards compatibility this should be relatively safe.
type GrpcAPIService struct {
	// Http Rule. Named Http in the actual proto. Changed to suppress linter warning.
	HTTP *annotations.Http `protobuf:"bytes,9,opt,name=http" json:"http,omitempty"`
}

// ProtoMessage returns an empty GrpcAPIService element
func (*GrpcAPIService) ProtoMessage() {}

// Reset resets the GrpcAPIService
func (m *GrpcAPIService) Reset() { *m = GrpcAPIService{} }

// String returns the string representation of the GrpcAPIService
func (m *GrpcAPIService) String() string { return proto.CompactTextString(m) }

func loadGrpcAPIServiceFromYAML(yamlFileContents []byte, yamlSourceLogName string) (*GrpcAPIService, error) {
	jsonContents, err := yaml.YAMLToJSON(yamlFileContents)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert gRPC API Configuration from YAML in '%v' to JSON: %v", yamlSourceLogName, err)
	}

	// As our GrpcAPIService is incomplete accept unkown fields.
	unmarshaler := jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}

	serviceConfiguration := GrpcAPIService{}
	if err := unmarshaler.Unmarshal(bytes.NewReader(jsonContents), &serviceConfiguration); err != nil {
		return nil, fmt.Errorf("Failed to parse gRPC API Configuration from YAML in '%v': %v", yamlSourceLogName, err)
	}

	return &serviceConfiguration, nil
}

func registerGrpcAPIService(registry *Registry, service *GrpcAPIService, sourceLogName string) error {
	if service.HTTP == nil {
		// Nothing to do
		return nil
	}

	for _, rule := range service.HTTP.GetRules() {
		selector := "." + strings.Trim(rule.GetSelector(), " ")
		if strings.ContainsAny(selector, "*, ") {
			return fmt.Errorf("Selector '%v' in %v must specify a single service method without wildcards", rule.GetSelector(), sourceLogName)
		}

		registry.externalHTTPRules[selector] = append(registry.externalHTTPRules[selector], rule)
	}

	return nil
}

// LoadGrpcAPIServiceFromYAML loads a gRPC API Configuration from the given YAML file
// and registers the HttpRule descriptions contained in it  as externalHTTPRules in
// the given registry. This must be done before loading the proto file.
//
// You can learn more about gRPC API Service descriptions from google's documentation
// at https://cloud.google.com/endpoints/docs/grpc/grpc-service-config
//
// Note that for the purposes of the gateway generator we only consider a subset of all
// available features google supports in their service descriptions.
func (r *Registry) LoadGrpcAPIServiceFromYAML(yamlFile string) error {
	yamlFileContents, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return fmt.Errorf("Failed to read gRPC API Configuration description from '%v': %v", yamlFile, err)
	}

	service, err := loadGrpcAPIServiceFromYAML(yamlFileContents, yamlFile)
	if err != nil {
		return err
	}

	return registerGrpcAPIService(r, service, yamlFile)
}
