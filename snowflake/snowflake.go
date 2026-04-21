package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"github.com/xframe-go/x/contracts"
)

type Generator struct {
	*snowflake.Node
}

func (g *Generator) Generate() string {
	return g.Node.Generate().Base36()
}

func New() contracts.IdGenerator {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return &Generator{node}
}
