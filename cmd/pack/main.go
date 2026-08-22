package main

import (
	"fmt"

	"github.com/pt-main/tap"
)

func NewParser() *tap.Parser {
	parser := tap.NewParser("pack", `[?MA]╭───────[?RT] [?BE]Pack cli[?RT]
[?MA]│[?RT]    Pack - file orchestrator. 
[?MA]│[?RT]    Only [?YW]humanmade[?RT], By [?YW]Pt[?RT].
[?MA]╰───────[?RT]

[?GN]Flags usage[?RT]
	[?BBK]--secure[?RT] to enable file encrypting/decrypting
	[?BBK]--zip[?RT] to compress file
	[?BBK]--key='...'[?RT] to set secure key

All flags used while pack creating must be used on all working with pack commands.`,
		[]string{"-h", "help"}, tap.DefaultParserConfig())
	parser.AddCommand("create", CreateHandler, `[?GN]Create file.[?GN]`,
		[]string{"output-file"}, []string{"input-file", "container-name"}, false)
	parser.AddCommand("push", PushHandler, `[?GN]Push file to pack file.[?GN]
Using like:
	[?BBK]pack push --file1="container1Name" --file1="container2Name"...[?RT]`,
		[]string{"pack-file"}, nil, false)
	parser.AddCommand("get", GetHandler, `[?GN]Get file from pack file.[?GN]`,
		[]string{"pack-file", "output-file", "container-name"}, nil, false)
	parser.AddCommand("struct", StructHandler, `[?GN]Show pack structure.[?GN]`,
		[]string{"pack-file"}, nil, false)
	parser.AddCommand("show", ShowHandler, `[?GN]Show file data.[?GN]`,
		[]string{"pack-file"}, nil, true)
	return parser
}

func main() {
	p := NewParser()
	if err := p.Main(); err != nil {
		fmt.Println(err)
	}
}
