// Filename: instructions.go
// Contents: The Instruction base struct and the other structs and methods
//	related to instructions.

package armsim

import (
	"io"
	"log"
	"os"
)

// Implements a Go interface allowing polymorphism, Go style.
type Instruction interface {
	Execute(cpu *CPU) error
	decode(base *baseInstruction) error
}

// Decodes an instruction.
//
// Parameters:
//	instructionBits - word of data representing the next instruction
//
// Returns:
//	instruction - a decoded instruction of type Instruction
func Decode(instructionBits uint32, logOut io.Writer) (instruction Instruction) {
	base := new(baseInstruction)
	if logOut == nil {
		logOut = os.Stderr
	}
	base.log = log.New(logOut, "Instruction Factory: ", 0)

	base.log.Printf("Decoding instruction: 0x%08x", instructionBits)

	// Set instruction bits
	base.InstructionBits = instructionBits
	base.shifter = new(BarrelShifter)

	// Get condition
	base.CondCode = ExtractShiftBits(instructionBits, 28, 32)
	base.log.Printf("Condition bits: %04b", base.CondCode)

	// Get instruction type
	base.Type = ExtractShiftBits(instructionBits, 25, 28)
	base.log.Printf("Type bits: %03b", base.Type)

	// Get Rn
	base.Rn = ExtractShiftBits(instructionBits, 16, 20)
	base.log.Printf("Rn: %d", base.Rn)

	// Get Rd
	base.Rd = ExtractShiftBits(instructionBits, 12, 16)
	base.log.Printf("Rd: %d", base.Rd)

	instruction = base.BuildFromBase()

	return
}

type baseInstruction struct {
	// The original bits of the instruction.
	InstructionBits uint32

	// Type bits
	Type uint32
	// Condition bits
	CondCode uint32
	// Destination register
	Rd uint32
	// First register operand
	Rn uint32

	log     *log.Logger
	shifter *BarrelShifter
}

func (bi *baseInstruction) BuildFromBase() (instruction Instruction) {
	// Check edge cases

	// Check type of instruction and call proper decode method
	switch bi.Type {
	case 0x0, 0x1:
		bi.log.Printf("Data Processing Instruction")
		instruction = new(dataInstruction)
	case 0x2:
		bi.log.Printf("Load/Store: Immediate Offset")
		instruction = new(loadStoreInstruction)
	case 0x3:
		bi.log.Printf("Load/Store: Register Offset")
		instruction = new(unimplementedInstruction)
	case 0x4:
		bi.log.Printf("Load/Store: Multiple")
		instruction = new(unimplementedInstruction)
	case 0x5:
		bi.log.Printf("Branch")
		instruction = new(branchInstruction)
	case 0xF:
		bi.log.Printf("Software Interrupt")
		instruction = new(unimplementedInstruction)
	default:
		bi.log.Printf("Unknown Instruction")
		instruction = new(unimplementedInstruction)
	}

	bi.log.SetPrefix("Instruction Decoding: ")
	instruction.decode(bi)

	return
}

type dataInstruction struct {
	// Embedding a general instruction
	*baseInstruction

	// Opcode
	Opcode byte
	// S bit
	S bool
	// I bit
	I bool
	// Second operand
	Operand2 uint32
}

const (
	AND = 0x0 // 0000
	SUB = 0x2 // 0010
	RSB = 0x3 // 0011
	ADD = 0x4 // 0100
	BIC = 0xE // 1110
	MOV byte = 0xD // 1101
	MNV = 0xF // 1111
)

// Executes a data instruction
//
// Parameters:
//  ram - a pointer to a block of memory, presumably ram
//  registers - a pointer to a block of memory, presumably a register bank
//
// Returns:
//  err - an error
func (di *dataInstruction) Execute(cpu *CPU) (err error) {
	di.log.SetPrefix("Data Instruction (Execute): ")

	// Parse the Operand2
	di.shifter = NewFromOperand2(di.Operand2, di.I, cpu)
	result := di.shifter.Shift()
	rn, _ := cpu.FetchRegisterFromInstruction(di.Rn)

	// Assertain specific instruction
	switch di.Opcode {
	case MOV, MNV:
		// Negate for MNV
		if di.Opcode == MNV {
			result ^= 0xFFFFFFFF
		}
		cpu.WriteRegisterFromInstruction(di.Rd, result)
	case ADD:
		// Rd = Rn + shifter_operand
		cpu.WriteRegisterFromInstruction(di.Rd, rn + result)
	case SUB:
		// Rd = Rn - shifter_operand
		cpu.WriteRegisterFromInstruction(di.Rd, rn - result)
	case RSB:
		// Rd = shifter_operand - Rn
		cpu.WriteRegisterFromInstruction(di.Rd, result - rn)
	case AND:
		// Rd = Rn AND shifter_operand
		cpu.WriteRegisterFromInstruction(di.Rd, rn & result)
	case BIC:
		// Rd = Rn AND NOT shifter_operand
		cpu.WriteRegisterFromInstruction(di.Rd, rn &^ result)
	}
	return
}

// Decodes a data instruction
//
// Parameters:
//  base - a generic instruction containing most information
//
// Returns:
//  err - an error
func (di *dataInstruction) decode(base *baseInstruction) (err error) {
	di.baseInstruction = base
	di.log.SetPrefix("Data Instruction (Decode): ")

	// Get I bit
	di.I = ExtractShiftBits(di.InstructionBits, 25, 26) == 1
	di.log.Printf("I bit: %01t", di.I)

	// Get opcode
	di.Opcode = byte(ExtractShiftBits(di.InstructionBits, 21, 25))
	di.log.Printf("Opcode bits: %04b", di.Opcode)

	// Get Operand2
	di.Operand2 = ExtractShiftBits(di.InstructionBits, 0, 12)
	di.log.Printf("Op2 bits: %012b", di.Operand2)

	// Get S bit
	di.S = ExtractShiftBits(di.InstructionBits, 20, 21) == 1
	di.log.Printf("S bit: %01t", di.S)

	return
}

type loadStoreInstruction struct {
	// Embedding a general instruction
	*baseInstruction

	// B bit
	B bool
	// L bit
	L bool
	// P bit
	P bool
	// U bit
	U bool
	// W bit
	W bool
}

// Executes a load/store instruction
//
// Parameters:
//  ram - a pointer to a block of memory, presumably ram
//  registers - a pointer to a block of memory, presumably a register bank
//
// Returns:
//  err - an error
func (lsi *loadStoreInstruction) Execute(cpu *CPU) (err error) {
	// Stub
	return
}

// Decodes a load/store instruction
//
// Parameters:
//  base - a generic instruction containing most information
//
// Returns:
//  err - an error
func (lsi *loadStoreInstruction) decode(base *baseInstruction) (err error) {
	// Stub
	return
}

type branchInstruction struct {
	// Embedding a general instruction
	*baseInstruction

	// L bit
	L bool

	// Offset bits
	Offset uint32
}

// Executes a branch instruction
//
// Parameters:
//  ram - a pointer to a block of memory, presumably ram
//  registers - a pointer to a block of memory, presumably a register bank
//
// Returns:
//  err - an error
func (bi *branchInstruction) Execute(cpu *CPU) (err error) {
	// Stub
	return
}

// Decodes a branch instruction
//
// Parameters:
//  base - a generic instruction containing most information
//
// Returns:
//  err - an error
func (bi *branchInstruction) decode(base *baseInstruction) (err error) {
	// Stub
	return
}

type unimplementedInstruction struct {
}

// Stub method to fake execution of unimplemented instructions
func (ui *unimplementedInstruction) Execute(cpu *CPU) (err error) {
	return
}

// Stub method to fake decoding of unimplemented instructions
func (ui *unimplementedInstruction) decode(base *baseInstruction) (err error) {
	return
}
