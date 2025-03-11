package contracts

type Router interface {
	Use(m ...Middleware)
	Get(uri string, handler Handler) Router
	Post(uri string, handler Handler) Router
	Put(uri string, handler Handler) Router
	Patch(uri string, handler Handler) Router
	Delete(uri string, handler Handler) Router
	Any(uri string, handler Handler) Router
	Prefix(uri string, m ...Middleware) Router
}

type Middleware func(request Request, next func() error) error
