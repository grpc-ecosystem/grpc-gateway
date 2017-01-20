package runtime

// Middleware is an interface which allows to create chains of handler functions
type Middleware func(HandlerFunc) HandlerFunc
