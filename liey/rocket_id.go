package liey

import (
	"cnb.cool/liey/liey-go/contracts"
	"cnb.cool/liey/liey-go/snowflake"
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
