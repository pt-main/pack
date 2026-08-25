package core

import (
	"github.com/pt-main/lc/engine"
	"github.com/pt-main/lc/engine/core"
	"github.com/pt-main/lc/parsing/byteParsing"
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
