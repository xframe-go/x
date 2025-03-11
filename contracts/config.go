package contracts

type ConfigValue interface {
	Symbol() any
}

type Config interface {
	Register(config ConfigValue)
	Get(key any) (any, bool)
}
