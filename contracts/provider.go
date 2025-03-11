package contracts

type Provider interface {
	Register(app Application)
}
