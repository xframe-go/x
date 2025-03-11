package contracts

type Handler func(req Request) Response

type RestHandler interface {
	Index(req Request) Response
	Store(req Request) Response
	Show(req Request) Response
	Edit(req Request) Response
	Destroy(req Request) Response
}
