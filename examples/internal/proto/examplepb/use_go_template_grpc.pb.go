// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package examplepb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// LoginServiceClient is the client API for LoginService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LoginServiceClient interface {
	// Login
	//
	// {{.MethodDescriptorProto.Name}} is a call with the method(s) {{$first := true}}{{range .Bindings}}{{if $first}}{{$first = false}}{{else}}, {{end}}{{.HTTPMethod}}{{end}} within the "{{.Service.Name}}" service.
	// It takes in "{{.RequestType.Name}}" and returns a "{{.ResponseType.Name}}".
	//
	// ## {{.RequestType.Name}}
	// | Field ID    | Name      | Type                                                       | Description                  |
	// | ----------- | --------- | ---------------------------------------------------------  | ---------------------------- | {{range .RequestType.Fields}}
	// | {{.Number}} | {{.Name}} | {{if eq .Label.String "LABEL_REPEATED"}}[]{{end}}{{.Type}} | {{fieldcomments .Message .}} | {{end}}
	//
	// ## {{.ResponseType.Name}}
	// | Field ID    | Name      | Type                                                       | Description                  |
	// | ----------- | --------- | ---------------------------------------------------------- | ---------------------------- | {{range .ResponseType.Fields}}
	// | {{.Number}} | {{.Name}} | {{if eq .Label.String "LABEL_REPEATED"}}[]{{end}}{{.Type}} | {{fieldcomments .Message .}} | {{end}}
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error)
	// Logout
	//
	// {{.MethodDescriptorProto.Name}} is a call with the method(s) {{$first := true}}{{range .Bindings}}{{if $first}}{{$first = false}}{{else}}, {{end}}{{.HTTPMethod}}{{end}} within the "{{.Service.Name}}" service.
	// It takes in "{{.RequestType.Name}}" and returns a "{{.ResponseType.Name}}".
	//
	// ## {{.RequestType.Name}}
	// | Field ID    | Name      | Type                                                       | Description                  |
	// | ----------- | --------- | ---------------------------------------------------------  | ---------------------------- | {{range .RequestType.Fields}}
	// | {{.Number}} | {{.Name}} | {{if eq .Label.String "LABEL_REPEATED"}}[]{{end}}{{.Type}} | {{fieldcomments .Message .}} | {{end}}
	//
	// ## {{.ResponseType.Name}}
	// | Field ID    | Name      | Type                                                       | Description                  |
	// | ----------- | --------- | ---------------------------------------------------------- | ---------------------------- | {{range .ResponseType.Fields}}
	// | {{.Number}} | {{.Name}} | {{if eq .Label.String "LABEL_REPEATED"}}[]{{end}}{{.Type}} | {{fieldcomments .Message .}} | {{end}}
	Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutReply, error)
}

type loginServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLoginServiceClient(cc grpc.ClientConnInterface) LoginServiceClient {
	return &loginServiceClient{cc}
}

var loginServiceLoginStreamDesc = &grpc.StreamDesc{
	StreamName: "Login",
}

func (c *loginServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error) {
	out := new(LoginReply)
	err := c.cc.Invoke(ctx, "/grpc.gateway.examples.internal.proto.examplepb.LoginService/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var loginServiceLogoutStreamDesc = &grpc.StreamDesc{
	StreamName: "Logout",
}

func (c *loginServiceClient) Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutReply, error) {
	out := new(LogoutReply)
	err := c.cc.Invoke(ctx, "/grpc.gateway.examples.internal.proto.examplepb.LoginService/Logout", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoginServiceService is the service API for LoginService service.
// Fields should be assigned to their respective handler implementations only before
// RegisterLoginServiceService is called.  Any unassigned fields will result in the
// handler for that method returning an Unimplemented error.
type LoginServiceService struct {
	// Login
	//
	// {{.MethodDescriptorProto.Name}} is a call with the method(s) {{$first := true}}{{range .Bindings}}{{if $first}}{{$first = false}}{{else}}, {{end}}{{.HTTPMethod}}{{end}} within the "{{.Service.Name}}" service.
	// It takes in "{{.RequestType.Name}}" and returns a "{{.ResponseType.Name}}".
	//
	// ## {{.RequestType.Name}}
	// | Field ID    | Name      | Type                                                       | Description                  |
	// | ----------- | --------- | ---------------------------------------------------------  | ---------------------------- | {{range .RequestType.Fields}}
	// | {{.Number}} | {{.Name}} | {{if eq .Label.String "LABEL_REPEATED"}}[]{{end}}{{.Type}} | {{fieldcomments .Message .}} | {{end}}
	//
	// ## {{.ResponseType.Name}}
	// | Field ID    | Name      | Type                                                       | Description                  |
	// | ----------- | --------- | ---------------------------------------------------------- | ---------------------------- | {{range .ResponseType.Fields}}
	// | {{.Number}} | {{.Name}} | {{if eq .Label.String "LABEL_REPEATED"}}[]{{end}}{{.Type}} | {{fieldcomments .Message .}} | {{end}}
	Login func(context.Context, *LoginRequest) (*LoginReply, error)
	// Logout
	//
	// {{.MethodDescriptorProto.Name}} is a call with the method(s) {{$first := true}}{{range .Bindings}}{{if $first}}{{$first = false}}{{else}}, {{end}}{{.HTTPMethod}}{{end}} within the "{{.Service.Name}}" service.
	// It takes in "{{.RequestType.Name}}" and returns a "{{.ResponseType.Name}}".
	//
	// ## {{.RequestType.Name}}
	// | Field ID    | Name      | Type                                                       | Description                  |
	// | ----------- | --------- | ---------------------------------------------------------  | ---------------------------- | {{range .RequestType.Fields}}
	// | {{.Number}} | {{.Name}} | {{if eq .Label.String "LABEL_REPEATED"}}[]{{end}}{{.Type}} | {{fieldcomments .Message .}} | {{end}}
	//
	// ## {{.ResponseType.Name}}
	// | Field ID    | Name      | Type                                                       | Description                  |
	// | ----------- | --------- | ---------------------------------------------------------- | ---------------------------- | {{range .ResponseType.Fields}}
	// | {{.Number}} | {{.Name}} | {{if eq .Label.String "LABEL_REPEATED"}}[]{{end}}{{.Type}} | {{fieldcomments .Message .}} | {{end}}
	Logout func(context.Context, *LogoutRequest) (*LogoutReply, error)
}

func (s *LoginServiceService) login(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/grpc.gateway.examples.internal.proto.examplepb.LoginService/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}
func (s *LoginServiceService) logout(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/grpc.gateway.examples.internal.proto.examplepb.LoginService/Logout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.Logout(ctx, req.(*LogoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterLoginServiceService registers a service implementation with a gRPC server.
func RegisterLoginServiceService(s grpc.ServiceRegistrar, srv *LoginServiceService) {
	srvCopy := *srv
	if srvCopy.Login == nil {
		srvCopy.Login = func(context.Context, *LoginRequest) (*LoginReply, error) {
			return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
		}
	}
	if srvCopy.Logout == nil {
		srvCopy.Logout = func(context.Context, *LogoutRequest) (*LogoutReply, error) {
			return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
		}
	}
	sd := grpc.ServiceDesc{
		ServiceName: "grpc.gateway.examples.internal.proto.examplepb.LoginService",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Login",
				Handler:    srvCopy.login,
			},
			{
				MethodName: "Logout",
				Handler:    srvCopy.logout,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "examples/internal/proto/examplepb/use_go_template.proto",
	}

	s.RegisterService(&sd, nil)
}

// NewLoginServiceService creates a new LoginServiceService containing the
// implemented methods of the LoginService service in s.  Any unimplemented
// methods will result in the gRPC server returning an UNIMPLEMENTED status to the client.
// This includes situations where the method handler is misspelled or has the wrong
// signature.  For this reason, this function should be used with great care and
// is not recommended to be used by most users.
func NewLoginServiceService(s interface{}) *LoginServiceService {
	ns := &LoginServiceService{}
	if h, ok := s.(interface {
		Login(context.Context, *LoginRequest) (*LoginReply, error)
	}); ok {
		ns.Login = h.Login
	}
	if h, ok := s.(interface {
		Logout(context.Context, *LogoutRequest) (*LogoutReply, error)
	}); ok {
		ns.Logout = h.Logout
	}
	return ns
}

// UnstableLoginServiceService is the service API for LoginService service.
// New methods may be added to this interface if they are added to the service
// definition, which is not a backward-compatible change.  For this reason,
// use of this type is not recommended.
type UnstableLoginServiceService interface {
	// Login
	//
	// {{.MethodDescriptorProto.Name}} is a call with the method(s) {{$first := true}}{{range .Bindings}}{{if $first}}{{$first = false}}{{else}}, {{end}}{{.HTTPMethod}}{{end}} within the "{{.Service.Name}}" service.
	// It takes in "{{.RequestType.Name}}" and returns a "{{.ResponseType.Name}}".
	//
	// ## {{.RequestType.Name}}
	// | Field ID    | Name      | Type                                                       | Description                  |
	// | ----------- | --------- | ---------------------------------------------------------  | ---------------------------- | {{range .RequestType.Fields}}
	// | {{.Number}} | {{.Name}} | {{if eq .Label.String "LABEL_REPEATED"}}[]{{end}}{{.Type}} | {{fieldcomments .Message .}} | {{end}}
	//
	// ## {{.ResponseType.Name}}
	// | Field ID    | Name      | Type                                                       | Description                  |
	// | ----------- | --------- | ---------------------------------------------------------- | ---------------------------- | {{range .ResponseType.Fields}}
	// | {{.Number}} | {{.Name}} | {{if eq .Label.String "LABEL_REPEATED"}}[]{{end}}{{.Type}} | {{fieldcomments .Message .}} | {{end}}
	Login(context.Context, *LoginRequest) (*LoginReply, error)
	// Logout
	//
	// {{.MethodDescriptorProto.Name}} is a call with the method(s) {{$first := true}}{{range .Bindings}}{{if $first}}{{$first = false}}{{else}}, {{end}}{{.HTTPMethod}}{{end}} within the "{{.Service.Name}}" service.
	// It takes in "{{.RequestType.Name}}" and returns a "{{.ResponseType.Name}}".
	//
	// ## {{.RequestType.Name}}
	// | Field ID    | Name      | Type                                                       | Description                  |
	// | ----------- | --------- | ---------------------------------------------------------  | ---------------------------- | {{range .RequestType.Fields}}
	// | {{.Number}} | {{.Name}} | {{if eq .Label.String "LABEL_REPEATED"}}[]{{end}}{{.Type}} | {{fieldcomments .Message .}} | {{end}}
	//
	// ## {{.ResponseType.Name}}
	// | Field ID    | Name      | Type                                                       | Description                  |
	// | ----------- | --------- | ---------------------------------------------------------- | ---------------------------- | {{range .ResponseType.Fields}}
	// | {{.Number}} | {{.Name}} | {{if eq .Label.String "LABEL_REPEATED"}}[]{{end}}{{.Type}} | {{fieldcomments .Message .}} | {{end}}
	Logout(context.Context, *LogoutRequest) (*LogoutReply, error)
}
