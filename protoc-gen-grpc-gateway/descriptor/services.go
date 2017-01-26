package descriptor

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	gateway_options "github.com/shilkin/grpc-gateway/options"
	"github.com/shilkin/grpc-gateway/protoc-gen-grpc-gateway/httprule"
	google_options "github.com/shilkin/grpc-gateway/third_party/googleapis/google/api"
)

// loadServices registers services and their methods from "targetFile" to "r".
// It must be called after loadFile is called for all files so that loadServices
// can resolve names of message types and their fields.
func (r *Registry) loadServices(file *File) error {
	glog.V(1).Infof("Loading services from %s", file.GetName())
	var svcs []*Service
	for _, sd := range file.GetService() {
		glog.V(2).Infof("Registering %s", sd.GetName())
		svc := &Service{
			File: file,
			ServiceDescriptorProto: sd,
		}
		for _, md := range sd.GetMethod() {
			glog.V(2).Infof("Processing %s.%s", sd.GetName(), md.GetName())
			opts, err := extractAPIOptions(md)
			if err != nil {
				glog.Errorf("Failed to extract ApiMethodOptions from %s.%s: %v", svc.GetName(), md.GetName(), err)
				return err
			}
			if opts == nil {
				glog.V(1).Infof("Found non-target method: %s.%s", svc.GetName(), md.GetName())
			}
			glog.V(2).Infof("API options for %s.%s: %#v", svc.GetName(), md.GetName(), opts)
			meth, err := r.newMethod(svc, md, opts)
			if err != nil {
				return err
			}
			svc.Methods = append(svc.Methods, meth)
		}
		if len(svc.Methods) == 0 {
			continue
		}
		glog.V(2).Infof("Registered %s with %d method(s)", svc.GetName(), len(svc.Methods))
		svcs = append(svcs, svc)
	}
	file.Services = svcs
	return nil
}

func (r *Registry) newMethod(svc *Service, md *descriptor.MethodDescriptorProto, opts *apiOptions) (*Method, error) {
	requestType, err := r.LookupMsg(svc.File.GetPackage(), md.GetInputType())
	if err != nil {
		return nil, err
	}
	responseType, err := r.LookupMsg(svc.File.GetPackage(), md.GetOutputType())
	if err != nil {
		return nil, err
	}
	meth := &Method{
		Service:               svc,
		MethodDescriptorProto: md,
		RequestType:           requestType,
		ResponseType:          responseType,
	}

	newBinding := func(opts *apiOptions, idx int) (*Binding, error) {
		var (
			httpMethod   string
			pathTemplate string
		)
		switch {
		case opts.httpRule.GetGet() != "":
			httpMethod = "GET"
			pathTemplate = opts.httpRule.GetGet()
			if opts.httpRule.Body != "" {
				return nil, fmt.Errorf("needs request body even though http method is GET: %s", md.GetName())
			}

		case opts.httpRule.GetPut() != "":
			httpMethod = "PUT"
			pathTemplate = opts.httpRule.GetPut()

		case opts.httpRule.GetPost() != "":
			httpMethod = "POST"
			pathTemplate = opts.httpRule.GetPost()

		case opts.httpRule.GetDelete() != "":
			httpMethod = "DELETE"
			pathTemplate = opts.httpRule.GetDelete()
			if opts.httpRule.Body != "" && !r.allowDeleteBody {
				return nil, fmt.Errorf("needs request body even though http method is DELETE: %s", md.GetName())
			}

		case opts.httpRule.GetPatch() != "":
			httpMethod = "PATCH"
			pathTemplate = opts.httpRule.GetPatch()

		case opts.httpRule.GetCustom() != nil:
			custom := opts.httpRule.GetCustom()
			httpMethod = custom.Kind
			pathTemplate = custom.Path

		default:
			glog.V(1).Infof("No pattern specified in google.api.HttpRule: %s", md.GetName())
			return nil, nil
		}

		parsed, err := httprule.Parse(pathTemplate)
		if err != nil {
			return nil, err
		}
		tmpl := parsed.Compile()

		if md.GetClientStreaming() && len(tmpl.Fields) > 0 {
			return nil, fmt.Errorf("cannot use path parameter in client streaming")
		}

		b := &Binding{
			Method:     meth,
			Index:      idx,
			PathTmpl:   tmpl,
			HTTPMethod: httpMethod,
			Middleware: opts.getMiddleware(),
		}

		for _, f := range tmpl.Fields {
			param, err := r.newParam(meth, f)
			if err != nil {
				return nil, err
			}
			b.PathParams = append(b.PathParams, param)
		}

		// TODO(yugui) Handle query params

		b.Body, err = r.newBody(meth, opts.httpRule.Body)
		if err != nil {
			return nil, err
		}

		return b, nil
	}
	b, err := newBinding(opts, 0)
	if err != nil {
		return nil, err
	}

	if b != nil {
		meth.Bindings = append(meth.Bindings, b)
	}
	for i, additional := range opts.httpRule.GetAdditionalBindings() {
		if len(additional.AdditionalBindings) > 0 {
			return nil, fmt.Errorf("additional_binding in additional_binding not allowed: %s.%s", svc.GetName(), meth.GetName())
		}
		apiOpts := &apiOptions{httpRule: additional, methodOpts: opts.methodOpts}
		b, err := newBinding(apiOpts, i+1)
		if err != nil {
			return nil, err
		}
		meth.Bindings = append(meth.Bindings, b)
	}

	return meth, nil
}

func extractAPIOptions(meth *descriptor.MethodDescriptorProto) (*apiOptions, error) { // (*options.HttpRule, error) {
	var opts apiOptions

	if meth.Options == nil {
		return nil, nil
	}
	// google api extension
	if proto.HasExtension(meth.Options, google_options.E_Http) {
		ext, err := proto.GetExtension(meth.Options, google_options.E_Http)
		if err != nil {
			return nil, err
		}
		httpRule, ok := ext.(*google_options.HttpRule)
		if !ok {
			return nil, fmt.Errorf("extension is %T; want an HttpRule", ext)
		}
		opts.httpRule = httpRule
	}
	// grpc gateway middleware extension
	if proto.HasExtension(meth.Options, gateway_options.E_MethodOptions) {
		ext, err := proto.GetExtension(meth.Options, gateway_options.E_MethodOptions)
		if err != nil {
			return nil, err
		}
		methodOpts, ok := ext.(*gateway_options.MethodOptions)
		if !ok {
			return nil, fmt.Errorf("extension is %T; want an MethodOptions", ext)
		}
		opts.methodOpts = methodOpts
	}

	return &opts, nil
}

func (r *Registry) newParam(meth *Method, path string) (Parameter, error) {
	msg := meth.RequestType
	fields, err := r.resolveFiledPath(msg, path)
	if err != nil {
		return Parameter{}, err
	}
	l := len(fields)
	if l == 0 {
		return Parameter{}, fmt.Errorf("invalid field access list for %s", path)
	}
	target := fields[l-1].Target
	switch target.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE, descriptor.FieldDescriptorProto_TYPE_GROUP:
		return Parameter{}, fmt.Errorf("aggregate type %s in parameter of %s.%s: %s", target.Type, meth.Service.GetName(), meth.GetName(), path)
	}
	return Parameter{
		FieldPath: FieldPath(fields),
		Method:    meth,
		Target:    fields[l-1].Target,
	}, nil
}

func (r *Registry) newBody(meth *Method, path string) (*Body, error) {
	msg := meth.RequestType
	switch path {
	case "":
		return nil, nil
	case "*":
		return &Body{FieldPath: nil}, nil
	}
	fields, err := r.resolveFiledPath(msg, path)
	if err != nil {
		return nil, err
	}
	return &Body{FieldPath: FieldPath(fields)}, nil
}

// lookupField looks up a field named "name" within "msg".
// It returns nil if no such field found.
func lookupField(msg *Message, name string) *Field {
	for _, f := range msg.Fields {
		if f.GetName() == name {
			return f
		}
	}
	return nil
}

// resolveFieldPath resolves "path" into a list of fieldDescriptor, starting from "msg".
func (r *Registry) resolveFiledPath(msg *Message, path string) ([]FieldPathComponent, error) {
	if path == "" {
		return nil, nil
	}

	root := msg
	var result []FieldPathComponent
	for i, c := range strings.Split(path, ".") {
		if i > 0 {
			f := result[i-1].Target
			switch f.GetType() {
			case descriptor.FieldDescriptorProto_TYPE_MESSAGE, descriptor.FieldDescriptorProto_TYPE_GROUP:
				var err error
				msg, err = r.LookupMsg(msg.FQMN(), f.GetTypeName())
				if err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("not an aggregate type: %s in %s", f.GetName(), path)
			}
		}

		glog.V(2).Infof("Lookup %s in %s", c, msg.FQMN())
		f := lookupField(msg, c)
		if f == nil {
			return nil, fmt.Errorf("no field %q found in %s", path, root.GetName())
		}
		if f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return nil, fmt.Errorf("repeated field not allowed in field path: %s in %s", f.GetName(), path)
		}
		result = append(result, FieldPathComponent{Name: c, Target: f})
	}
	return result, nil
}
