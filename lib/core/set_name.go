package core

import (
	"github.com/pt-main/lc/v2/engine"
	"github.com/pt-main/lc/v2/engine/core"
	"github.com/pt-main/lc/v2/parsing/byteParsing"
)

func (c *Core) setName(e engine.ByteEngineInterface, a *byteParsing.ParsedBytes) core.ErrorInterface {
	_has_name, ok := c.Scope["name"]
	if !ok {
		return core.Err("setName", "Invalid scope: has no 'name' field")
	}
	has_name, ok := _has_name.(bool)
	if !ok {
		return core.Err("setName", "Invalid scope: invalid 'name' field")
	}
	if has_name {
		return core.Err("setName", "Can't set name: name already set")
	}
	if len(a.Args) != 1 {
		return core.Err("setName", "Invalid argument length: %d"+
			"(must be equal to 1)", len(a.Args))
	}
	c.Name = string(a.Args[0])
	c.Scope["name"] = true
	return nil
}
