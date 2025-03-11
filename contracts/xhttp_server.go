package contracts

type XServer interface {
	Router

	Start() error
}

type RoutingProvider func(router Router)
