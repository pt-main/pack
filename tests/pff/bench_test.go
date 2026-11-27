package pff

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/iancoleman/orderedmap"
	"github.com/pt-main/pack/lib/core"
	"github.com/pt-main/pack/lib/pff"
)

func generateDummyData(fileCount, fileSize int) *orderedmap.OrderedMap {
	om := orderedmap.New()
	chunk := bytes.Repeat([]byte("pack-benchmark-payload-data-1234567890\n"), fileSize/39+1)[:fileSize]
	for i := 0; i < fileCount; i++ {
		om.Set(fmt.Sprintf("dir/file_%d.txt", i), chunk)
	}
	return om
}

func BenchmarkPack(b *testing.B) {
	om := generateDummyData(100, 10*1024)
	testKey := "12345678901234567890123456789012"

	modes := []struct {
		name  string
		flags map[string]string
	}{
		{"Raw", map[string]string{}},
		{"ZipOnly", map[string]string{"zip": ""}},
		{"SecureOnly", map[string]string{"secure": "", "key": testKey}},
		{"ZipAndSecure", map[string]string{"zip": "", "secure": "", "key": testKey}},
	}

	for _, mode := range modes {
		b.Run("Encode_"+mode.name, func(b *testing.B) {
			b.ReportAllocs()
			c := core.NewCore(om)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				pf := pff.NewPFfromCore(c, mode.flags)
				_, err := pf.GetFromCore()
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		c := core.NewCore(om)
		pfInit := pff.NewPFfromCore(c, mode.flags)
		packedBytes, err := pfInit.GetFromCore()
		if err != nil {
			b.Fatal(err)
		}

		b.Run("Decode_"+mode.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				pf := pff.NewPFfromFile(packedBytes, mode.flags)
				_, err := pf.SetToCore()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
