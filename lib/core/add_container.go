package core

import (
	"github.com/pt-main/lc/v2/engine"
	"github.com/pt-main/lc/v2/engine/core"
	"github.com/pt-main/lc/v2/parsing/byteParsing"
)

func (c *Core) addContainer(e engine.ByteEngineInterface, a *byteParsing.ParsedBytes) core.ErrorInterface {
	if len(a.Args) != 2 {
		return core.Err("addContainer", "Invalid argument length: %d"+
			"(must be equal to 2: 'container_name' string and 'container_value' bytes)", len(a.Args))
	}
	c.Containers.Set(string(a.Args[0]), a.Args[1])
	return nil
}
