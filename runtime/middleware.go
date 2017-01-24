package runtime

// Middleware is an interface that declares the function creating HandlerFunc decorators.
// It's designed to create request handling pipelines.
type Middleware func(HandlerFunc) HandlerFunc
