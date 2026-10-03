package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/iancoleman/orderedmap"
	"github.com/pt-main/pack/lib/core"
	"github.com/pt-main/pack/lib/pff"
	"github.com/pt-main/tap/go"
	"github.com/pt-main/tap/go/color"
)

func open(file string) ([]byte, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("Open: %v", err)
	}
	return data, nil
}

func write(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Write: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		return fmt.Errorf("Write: %v", err)
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("Write: %v", err)
	}
	return nil
}

func CreateHandler(p *tap.Parser, args []string) error {
	om := orderedmap.New()
	var output = args[0]
	if len(args) == 3 {
		input, err := open(args[0])
		if err != nil {
			return fmt.Errorf("Create: %v", err)
		}
		output = args[1]
		name := args[2]
		om.Set(name, input)
	}
	c := core.NewCore(om)
	pff := pff.NewPFfromCore(c, p.Flags)
	f, _ := pff.GetFromCore()
	err := write(output, f)
	if err != nil {
		return fmt.Errorf("Create: %v", err)
	}
	return nil
}

func PushHandler(p *tap.Parser, args []string) error {
	input, err := open(args[0])
	if err != nil {
		return fmt.Errorf("Push: %v", err)
	}
	pff := pff.NewPFfromFile(input, p.Flags)
	for key, val := range p.Flags {
		if slices.Contains([]string{"secure", "zip", "key"}, key) {
			continue
		}
		if val == "" {
			val = key
		}
		what, err := open(key)
		if err != nil {
			return fmt.Errorf("Push: %v", err)
		}
		pff.Core.Containers.Set(val, what)

	}
	f, err := pff.GetFromCore()
	if err != nil {
		return fmt.Errorf("Push: %v", err)
	}
	err = write(args[0], f)
	if err != nil {
		return fmt.Errorf("Push: %v", err)
	}
	return nil
}

func GetHandler(p *tap.Parser, args []string) error {
	input, err := open(args[0])
	if err != nil {
		return fmt.Errorf("Get: %v", err)
	}
	pff := pff.NewPFfromFile(input, p.Flags)
	file, ok := pff.Core.GetContainer(args[1])
	if !ok {
		return fmt.Errorf("Get: Can't get '%s' container", args[1])
	}
	err = write(args[2], file)
	if err != nil {
		return fmt.Errorf("Get: %v", err)
	}
	return nil
}

func StructHandler(p *tap.Parser, args []string) error {
	data, err := open(args[0])
	if err != nil {
		return fmt.Errorf("Struct: open: %v", err)
	}

	pff := pff.NewPFfromFile(data, p.Flags)
	c, err := pff.SetToCore()
	if err != nil {
		return fmt.Errorf("Struct: read: %v", err)
	}

	color.PrintlnColored("[?GN]╭───────[?RT] [?YW]Struct[?RT]")
	color.PrintColored("[?GN]│[?RT]    ")
	fmt.Printf("Magic: %s\n", c.Magic)
	color.PrintColored("[?GN]│[?RT]    ")
	fmt.Printf("Name: %s\n", c.Name)
	color.PrintColored("[?GN]│[?RT]    ")
	fmt.Printf("Meta: %s\n", string(c.Meta))
	color.PrintColored("[?GN]│[?RT]    ")
	fmt.Println("Containers:")
	for _, key := range c.Containers.Keys() {
		val, _ := c.GetContainer(key)
		color.PrintColored("[?GN]│[?RT]    ")
		fmt.Printf(" - %s: %d bytes\n", key, len(val))
	}
	color.PrintlnColored("[?GN]╰───────[?RT]")
	return nil
}

func ShowHandler(p *tap.Parser, args []string) error {
	data, err := open(args[0])
	if err != nil {
		return fmt.Errorf("Show: open: %v", err)
	}

	pff := pff.NewPFfromFile(data, p.Flags)
	c, err := pff.SetToCore()
	if err != nil {
		return fmt.Errorf("Show: read: %v", err)
	}

	for idx, file := range args[1:] {
		filedata, ok := c.GetContainer(file)
		if !ok {
			return fmt.Errorf("Invalid container name: %v", file)
		} else {
			color.PrintlnColored("[?GN]%v[?RT]: [?YW]%v", idx, file)
			fmt.Println(string(filedata))
		}
	}
	return nil
}

func FilesZipHandler(p *tap.Parser, args []string) error {
	packFile := args[0]
	dir := args[1]

	om := orderedmap.New()
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		om.Set(rel, data)
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk error: %w", err)
	}

	c := core.NewCore(om)
	pf := pff.NewPFfromCore(c, p.Flags)
	pfData, err := pf.GetFromCore()
	if err != nil {
		return fmt.Errorf("pack error: %w", err)
	}
	err = write(packFile, pfData)
	if err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	color.PrintlnColored("[?GN]✅ Packed %d files into %s[?RT]", len(om.Keys()), packFile)
	return nil
}

func FilesUnzipHandler(p *tap.Parser, args []string) error {
	packFile := args[0]
	dir := args[1]

	data, err := open(packFile)
	if err != nil {
		return fmt.Errorf("open pack: %w", err)
	}
	pf := pff.NewPFfromFile(data, p.Flags)
	c, err := pf.SetToCore()
	if err != nil {
		return fmt.Errorf("parse pack: %w", err)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	extracted := 0
	for _, key := range c.Containers.Keys() {
		val, _ := c.GetContainer(key)
		fullPath := filepath.Join(dir, filepath.FromSlash(key))

		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("mkdir for %s: %w", fullPath, err)
		}
		if err := write(fullPath, val); err != nil {
			return fmt.Errorf("write %s: %w", fullPath, err)
		}
		extracted++
	}
	color.PrintlnColored("[?GN]✅ Extracted %d files to %s[?RT]", extracted, dir)
	return nil
}
