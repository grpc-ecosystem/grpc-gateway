package descriptor

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor/openapiconfig"
	"google.golang.org/protobuf/encoding/protojson"
)

func loadOpenAPIConfigFromYAML(yamlFileContents []byte, yamlSourceLogName string) (*openapiconfig.OpenAPIConfig, error) {
	jsonContents, err := yaml.YAMLToJSON(yamlFileContents)
	if err != nil {
		return nil, fmt.Errorf("failed to convert OpenAPI Configuration from YAML in '%v' to JSON: %v", yamlSourceLogName, err)
	}

	// Reject unknown fields because OpenAPIConfig is only used here
	unmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}

	openapiConfiguration := openapiconfig.OpenAPIConfig{}
	if err := unmarshaler.Unmarshal(jsonContents, &openapiConfiguration); err != nil {
		return nil, fmt.Errorf("failed to parse gRPC API Configuration from YAML in '%v': %v", yamlSourceLogName, err)
	}

	return &openapiConfiguration, nil
}

func registerOpenAPIOptions(registry *Registry, openAPIConfig *openapiconfig.OpenAPIConfig, yamlSourceLogName string) error {
	if openAPIConfig.OpenapiOptions == nil {
		// Nothing to do
		return nil
	}

	if err := registry.RegisterOpenAPIOptions(openAPIConfig.OpenapiOptions); err != nil {
		return fmt.Errorf("failed to register option in %s: %s", yamlSourceLogName, err)
	}
	return nil
}

// LoadOpenAPIConfigFromYAML loads an  OpenAPI Configuration from the given YAML file
// and registers the OpenAPI options the given registry.
// This must be done after loading the proto file.
func (r *Registry) LoadOpenAPIConfigFromYAML(yamlFile string) error {
	yamlFileContents, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return fmt.Errorf("failed to read gRPC API Configuration description from '%v': %v", yamlFile, err)
	}

	config, err := loadOpenAPIConfigFromYAML(yamlFileContents, yamlFile)
	if err != nil {
		return err
	}

	return registerOpenAPIOptions(r, config, yamlFile)
}
