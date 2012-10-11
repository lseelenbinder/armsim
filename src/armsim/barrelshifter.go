// Filename: barrelshifter.go
// Contents: The BarrelShifter class

package armsim

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
}

func NewFromOperand2(operand2 uint32, i bool, cpu *CPU) (b *BarrelShifter) {
	if i {
		b = &BarrelShifter{ROR, ExtractShiftBits(operand2, 8, 12) * 2, ExtractShiftBits(operand2, 0, 8)}
	} else {
		rm, _ := cpu.FetchRegisterFromInstruction(ExtractShiftBits(operand2, 0, 4))
		shift := ExtractShiftBits(operand2, 5, 7)
		if (operand2 & 0x10) == 0 {
			// Immediate shift
			b = &BarrelShifter{shift, ExtractShiftBits(operand2, 7, 12), rm}
		} else {
			// Register shift
			rs, _ := cpu.FetchRegisterFromInstruction(ExtractShiftBits(operand2, 8, 12))
			b = &BarrelShifter{shift, rs, rm}
		}
	}
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
	}
	return
}

func (b *BarrelShifter) Rs() (rs uint32) {
	return b.ShiftAmount
}

func (b *BarrelShifter) Rm() (rs uint32) {
	return b.Data
}

func ror(value uint32, nBits uint32) (result uint32) {
	return (value >> nBits) | (value<<(32-nBits))&0xFFFFFFFF
}
