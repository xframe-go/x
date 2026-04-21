package x

import (
	"github.com/xframe-go/x/contracts"
	"github.com/xframe-go/x/snowflake"
)

func (r *Rocket) getIDGenerator() contracts.IdGenerator {
	if r.id != nil {
		return r.id
	}

	r.id = snowflake.New()
	return r.id
}

func ID() contracts.IdGenerator {
	return rocket.getIDGenerator()
}
