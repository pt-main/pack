package core

import (
	"github.com/pt-main/lc/v2/engine"
	"github.com/pt-main/lc/v2/engine/core"
	"github.com/pt-main/lc/v2/parsing/byteParsing"
)

func (c *Core) addMeta(e engine.ByteEngineInterface, a *byteParsing.ParsedBytes) core.ErrorInterface {
	if len(a.Args) < 1 {
		return nil
	}
	for _, arg := range a.Args {
		c.Meta = append(c.Meta, arg...)
	}
	return nil
}
