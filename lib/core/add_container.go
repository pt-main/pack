package core

import (
	"strconv"

	"github.com/pt-main/lc/engine"
	"github.com/pt-main/lc/engine/core"
	"github.com/pt-main/lc/parsing/byteParsing"
)

func (c *Core) addContainer(e engine.ByteEngineInterface, a *byteParsing.ParsedBytes) core.ErrorInterface {
	if len(a.Args) != 2 {
		return core.Err("addContainer", "Invalid argument length: "+strconv.Itoa(len(a.Args))+
			"(must be equal to 2: 'container_name' string and 'container_value' bytes)")
	}
	c.Containers.Set(string(a.Args[0]), a.Args[1])
	return nil
}
