package contracts

type ServerConfig struct {
	Port int

	Host string

	JsonMarshal func(v interface{}) ([]byte, error)

	JsonUnmarshal func(data []byte, v interface{}) error

	RoutingProviders []RoutingProvider

	Middlewares []Middleware
}
