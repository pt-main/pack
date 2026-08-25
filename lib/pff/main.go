// pff - Pack file format
package pff

import (
	"errors"
	"fmt"
	"slices"

	"github.com/iancoleman/orderedmap"
	"github.com/pt-main/pack/lib/core"
)

type PF struct { // pack file
	Core   *core.Core
	Flags  []string
	File   []byte
	Values map[string]string
}

func NewPFfromCore(core *core.Core, processFlags map[string]string) *PF {
	flags := make([]string, 0)
	values := make(map[string]string)
	for key, val := range processFlags {
		if val == "" {
			flags = append(flags, key)
		} else {
			values[key] = val
		}
	}
	return &PF{
		Core:   core,
		Flags:  flags,
		Values: values,
		File:   make([]byte, 0),
	}
}

func NewPFfromFile(file []byte, processFlags map[string]string) *PF {
	flags := make([]string, 0)
	values := make(map[string]string)
	for key, val := range processFlags {
		if val == "" {
			flags = append(flags, key)
		} else {
			values[key] = val
		}
	}
	return &PF{
		Core:   core.NewCore(orderedmap.New()),
		Flags:  flags,
		Values: values,
		File:   file,
	}
}

func (f *PF) Convert() error {
	var err error
	if slices.Contains(f.Flags, "zip") {
		err = f.Compress()
		if err != nil {
			return fmt.Errorf("Convert (1): %v", err)
		}
	}
	if slices.Contains(f.Flags, "secure") {
		key, ok := f.Values["key"]
		if !ok {
			return errors.New("Convert (2): Can't find 'key' key in values" +
				" (to encrypt file)")
		}
		f.File, err = f.Encrypt(f.File, []byte(key))
		if err != nil {
			return fmt.Errorf("Convert (3): %v", err)
		}
	}
	return nil
}

func (f *PF) Deconvert() error {
	var err error
	if slices.Contains(f.Flags, "secure") {
		key, ok := f.Values["key"]
		if !ok {
			return errors.New("Deconvert (1): Can't find 'key' key in values" +
				" (to decrypt file)")
		}
		f.File, err = f.Decrypt(f.File, []byte(key))
		if err != nil {
			return fmt.Errorf("Deconvert (2): %v", err)
		}
	}
	if slices.Contains(f.Flags, "zip") {
		err = f.Decompress()
		if err != nil {
			return fmt.Errorf("Deconvert (3): %v", err)
		}
	}
	return nil
}

func (f *PF) Apply() error {
	var err error
	if f.File != nil && len(f.File) > 0 {
		err = f.Convert()
	} else {
		err = f.Deconvert()
	}
	if err != nil {
		return fmt.Errorf("Apply: %v", err)
	}
	return nil
}

func (f *PF) GetFromCore() ([]byte, error) {
	var err error
	f.File, err = f.Core.CreateFile()
	if err != nil {
		return nil, fmt.Errorf("GetFromCore (1): %v", err)
	}
	err = f.Apply()
	if err != nil {
		return nil, fmt.Errorf("GetFromCore (2): %v", err)
	}
	file := f.File
	f.File = nil
	return file, nil
}

func (f *PF) SetToCore() (*core.Core, error) {
	var err error
	err = f.Apply()
	if err != nil {
		return nil, fmt.Errorf("SetToCore (1): %v", err)
	}
	err = f.Core.ReadFile(f.File)
	if err != nil {
		return nil, fmt.Errorf("SetToCore (2): %v", err)
	}
	return f.Core, nil
}
