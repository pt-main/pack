package core

import (
	"strconv"

	"github.com/pt-main/lc/engine"
	"github.com/pt-main/lc/engine/core"
	"github.com/pt-main/lc/parsing/byteParsing"
)

func (c *Core) addMeta(e engine.ByteEngineInterface, a *byteParsing.ParsedBytes) core.ErrorInterface {
	if len(a.Args) < 1 {
		return core.Err("addMeta", "Invalid argument length: "+strconv.Itoa(len(a.Args))+
			"(must be more or equal to 1)")
	}
	for _, arg := range a.Args {
		c.Meta = append(c.Meta, arg...)
	}
	return nil
}
