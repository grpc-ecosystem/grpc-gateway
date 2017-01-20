package runtime

type Middleware func(HandlerFunc) HandlerFunc
