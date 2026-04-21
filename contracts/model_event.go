package contracts

type ModelEvent[T any] interface {
	GetEventType() string
	GetModel() T
}

type ModelEventListener[T any] interface {
	OnCreated(model T)
	OnUpdated(oldModel, newModel T)
	OnDeleted(model T)
}
