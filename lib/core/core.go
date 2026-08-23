package core

import (
	"context"
	"errors"
	"slices"

	"github.com/iancoleman/orderedmap"
	"github.com/pt-main/lc"
	"github.com/pt-main/lc/engine"
	"github.com/pt-main/lc/engine/core"
	"github.com/pt-main/lc/parsing/byteParsing"
	"github.com/pt-main/lc/public"
	"github.com/pt-main/lc/tooling/bytecode"
	"github.com/pt-main/tap/color"
)

var gc = bytecode.GenerationConfig{
	CommandBytelen:   1,
	ArgscountBytelen: 2,
	ArglenBytelen:    32,
	Endianess:        public.BigEndian,
}

const API_VER = "0"

type Core struct {
	Magic      []byte
	Name       string
	Meta       []byte
	Containers *orderedmap.OrderedMap
	Scope      core.ScopeType
}

func (c *Core) GetContainer(name string) ([]byte, bool) {
	_cont, ok := c.Containers.Get(name)
	if !ok {
		return nil, ok
	}
	cont, ok := _cont.([]byte)
	return cont, ok
}

func (c *Core) CreateFile() ([]byte, error) {
	file := c.Magic
	ig := bytecode.InstructionsGenerator{
		Config: gc,
	}
	name := []byte(c.Name)
	if len(name) == 0 {
		return nil, errors.New("Name can't be empty")
	}
	file = append(file, ig.Generate(0, [][]byte{name})...)
	meta := [][]byte{}
	if c.Meta != nil && len(c.Meta) != 0 {
		meta = append(meta, c.Meta)
	}
	file = append(file, ig.Generate(1, meta)...)
	for _, container := range c.Containers.Keys() {
		cont, ok := c.GetContainer(container)
		if !ok {
			return nil, errors.New("CreateFile: Invalid container: Invlid value type (all values must has []byte type)")
		}
		file = append(file, ig.Generate(2, [][]byte{[]byte(container), cont})...)
	}
	return file, nil
}

func (c *Core) ReadFile(file []byte) error {
	magicLen := len(c.Magic)
	if len(file) < magicLen || !slices.Equal(file[:magicLen], c.Magic) {
		return errors.New("ReadFile: Invalid file: invalid magic")
	}
	file = file[magicLen:]
	p, err := NewParser(c)
	if err != nil {
		return err
	}
	return p.Process(file)
}

func NewCore(input *orderedmap.OrderedMap) *Core {
	if input == nil {
		input = orderedmap.New()
	}
	s := make(core.ScopeType)
	s["name"] = false
	return &Core{
		Magic:      []byte("PACK@" + API_VER),
		Name:       "pack-file",
		Meta:       make([]byte, 0),
		Scope:      s,
		Containers: input,
	}
}

func NewParser(c *Core) (*engine.ByteEngine, error) {
	idx := 0
	e := lc.NewByteEngine(
		0, make([]string, 0), true,
		&byteParsing.Parser1{
			Config: byteParsing.Parser1Config{
				GConfig: gc, Shifter: *bytecode.NewShift(make([]byte, 0), &idx),
			},
		}, public.BigEndian,
		color.ColorEnabled,
		context.Background(),
	)
	e.NewCommandFull(0, c.setName, "set_name", true)
	e.NewCommandFull(1, c.addMeta, "add_meta", true)
	e.NewCommandFull(2, c.addContainer, "add_container", true)
	// e.UEP.Logger.Logging[public.LogParsing] = true
	// e.UEP.Logger.Logging[public.LogEvents] = true
	// e.UEP.Logger.Logging[public.LogVerbose] = true
	return e, nil
}
