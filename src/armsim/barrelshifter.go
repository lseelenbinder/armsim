// Filename: barrelshifter.go
// Contents: The BarrelShifter class

package armsim

import (
	"fmt"
	"log"
)

const (
	LSL uint32 = iota
	LSR
	ASR
	ROR
)

type BarrelShifter struct {
	Type        uint32
	ShiftAmount uint32
	Data        uint32
	Rs          uint32
	Rn          uint32
	i           bool
	log         *log.Logger
}

func NewFromOperand2(operand2 uint32, i bool, cpu *CPU) (b *BarrelShifter) {
	var shift, shift_amount, data uint32
	var rs, rn uint32 = 17, 17
	if i {
		shift = ROR
		shift_amount = ExtractShiftBits(operand2, 8, 12) * 2
		data = ExtractShiftBits(operand2, 0, 8)
	} else {
		rn = ExtractShiftBits(operand2, 0, 4)
		data, _ = cpu.FetchRegisterFromInstruction(rn)
		shift = ExtractShiftBits(operand2, 5, 7)
		if (operand2 & 0x10) == 0 {
			// Immediate shift
			shift_amount = ExtractShiftBits(operand2, 7, 12)
		} else {
			// Register shift
			rs = ExtractShiftBits(operand2, 8, 12)
			shift_amount, _ = cpu.FetchRegisterFromInstruction(rs)
		}
	}

	b = &BarrelShifter{shift, shift_amount, data, rs, rn, i, log.New(cpu.logOut, "BarrelShifter: ", 0)}

	return
}

func (b *BarrelShifter) Shift() (result uint32) {
	switch b.Type {
	case ROR:
		result = ror(b.Data, b.ShiftAmount)
	case LSL:
		result = b.Data << b.ShiftAmount
	case LSR:
		result = b.Data >> b.ShiftAmount
	case ASR:
		result = asr(b.Data, b.ShiftAmount)
	}
	return
}

func (b *BarrelShifter) GetRs() (rs uint32) {
	return b.ShiftAmount
}

func (b *BarrelShifter) GetRm() (rm uint32) {
	return b.Data
}

func (b *BarrelShifter) Disassembly() (operands string) {
	var mnemonic, data string
	if b.i {
		return fmt.Sprintf("#%d", b.Shift())
	} else {
		switch b.Type {
		case ROR:
			mnemonic = "ror"
		case LSL:
			mnemonic = "lsl"
		case LSR:
			mnemonic = "lsr"
		case ASR:
			mnemonic = "asr"
		}
		if b.Rs < 16 {
			// Register shift
			data = fmt.Sprintf("r%d", b.GetRs())
		} else {
			// Immediate shift
			data = fmt.Sprintf("#%d", b.ShiftAmount)
		}
	}
	operands = fmt.Sprintf("r%d, %s %s", b.Rn, mnemonic, data)
	b.log.Println(operands)
	return
}

func ror(value, nBits uint32) (result uint32) {
	return (value >> nBits) | (value<<(32-nBits))&0xFFFFFFFF
}

func asr(value, nBits uint32) (result uint32) {
	var mask uint32 = 0x0
	if (value & 0x80000000) > 0 {
		for i := 0; uint32(i) < nBits; i++ {
			mask >>= 1
			mask += 0x80000000
		}
	}
	result = (value >> nBits) | mask
	return
}
