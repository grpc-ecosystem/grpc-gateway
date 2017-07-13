package runtime

type MethodDesciption interface {
	ServiceName() string
	MethodName() string
	InputType() string
	OutputType() string
	HTTPMethod() string
	Pattern() string
}
