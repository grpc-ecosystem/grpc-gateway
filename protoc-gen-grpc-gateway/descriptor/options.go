package descriptor

import (
	gateway_options "github.com/shilkin/grpc-gateway/options"
	google_options "github.com/shilkin/grpc-gateway/third_party/googleapis/google/api"
)

type apiOptions struct {
	httpRule   *google_options.HttpRule
	methodOpts *gateway_options.MethodOptions
}

func (opts *apiOptions) getMiddleware() []string {
	if opts.methodOpts == nil {
		return []string{}
	}
	return opts.methodOpts.Middleware
}
