package descriptor

import options "github.com/shilkin/grpc-gateway/third_party/googleapis/google/api"

type apiOptions struct {
	httpRule   *options.HttpRule
	middleware []string
}
