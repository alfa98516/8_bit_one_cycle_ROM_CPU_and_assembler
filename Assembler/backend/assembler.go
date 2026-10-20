package backend

import (
	"fmt"
	"os"

	. "Assembler/common"
	. "Assembler/frontend"
)

type Assembler struct {
	out    *os.File
	parser Parser
}

// Encoding:
//
//	imm : 00_vvvvvv   (6-bit value)
//	alu : 01_000_ooo
//	mov : 10_ddd_sss  (ddd = dest, 0-5 = A-F, 6 = out; sss = source A-F)
//	nop/jumps : 11_000_ooo
var OpBits = map[TokenType]byte{
	TtKwNand: 0b01_000_000,
	TtKwOr:   0b01_000_001,
	TtKwAnd:  0b01_000_010,
	TtKwNor:  0b01_000_011,
	TtKwAdd:  0b01_000_100,
	TtKwSub:  0b01_000_101,
	TtKwXor:  0b01_000_110,
	TtKwNot:  0b01_000_111,
	TtKwNop:  0b11_000_000,
	TtKwJmp:  0b11_000_001,
	TtKwJz:   0b11_000_010,
	TtKwJnz:  0b11_000_011,
	TtKwJl:   0b11_000_100,
	TtKwJge:  0b11_000_101,
	TtKwJle:  0b11_000_110,
	TtKwJg:   0b11_000_111,
}

const (
	immBase = 0b00_000000
	movBase = 0b10_000_000
	outDest = 0b110
)

func regBits(t TokenType) byte {
	if t == TtKwOut {
		return outDest
	}
	return byte(t - TtKwA) // A..F are consecutive tokens
}

func encode(in Instr) byte {
	switch in.Op {
	case TtKwImm:
		return immBase | byte(in.Imm)
	case TtKwMov:
		return movBase | regBits(in.Dst)<<3 | regBits(in.Src)
	default:
		b, ok := OpBits[in.Op]
		if !ok {
			panic(fmt.Sprintf("no encoding for %s at %s", in.Op, in.Pos))
		}
		return b
	}
}

func outputString(file *os.File, line string) {
	fmt.Print(line)
	if _, err := file.WriteString(line); err != nil {
		panic(err)
	}
}

func (assembler *Assembler) Init(ofile *os.File, src []byte) {
	assembler.out = ofile
	errorHandler := func(pos Position, msg string) {
		outputString(ofile, fmt.Sprintf("Error at (%d %d): %s\n", pos.Line, pos.Col, msg))
	}
	assembler.parser.ParserSrc(src, errorHandler)
	assembler.Assemble()
}

const romWords = 256 // size of the ROM image, matches your example file

// Assemble writes a Logisim-evolution ROM image in the
// "v3.0 hex words addressed" format: 16 words per row, unused words zero.
func (assembler *Assembler) Assemble() {
	instrs := assembler.parser.Instrs()
	if len(instrs) > romWords {
		panic(fmt.Sprintf("program has %d instructions, ROM holds %d", len(instrs), romWords))
	}

	image := make([]byte, romWords)
	for i, in := range instrs {
		image[i] = encode(in)
	}

	outputString(assembler.out, "v3.0 hex words addressed\n")
	for row := 0; row < romWords; row += 16 {
		line := fmt.Sprintf("%02x:", row)
		for _, b := range image[row : row+16] {
			line += fmt.Sprintf(" %02x", b)
		}
		outputString(assembler.out, line+"\n")
	}
}
