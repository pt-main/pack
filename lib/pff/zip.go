package pff

import (
	"bytes"
	"compress/gzip"
	"io"
)

func (f *PF) Compress() error {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	_, err := gz.Write(f.File)
	if err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	f.File = buf.Bytes()
	return nil
}

func (f *PF) Decompress() error {
	r, err := gzip.NewReader(bytes.NewReader(f.File))
	if err != nil {
		return err
	}
	defer r.Close()
	decompressed, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	f.File = decompressed
	return nil
}
