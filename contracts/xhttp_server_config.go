package contracts

import "embed"

type ServerConfig struct {
	Port int

	Host string

	PublishDir string

	PublicFS *embed.FS

	JsonMarshal func(v interface{}) ([]byte, error)

	JsonUnmarshal func(data []byte, v interface{}) error

	RoutingProviders []RoutingProvider

	Middlewares []Middleware
}
