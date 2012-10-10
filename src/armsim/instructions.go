// Filename: instructions.go
// Contents: The Instruction base struct and the other structs and methods
//	related to instructions.

package armsim

import (
	"log"
	"io"
)

// Implements a Go interface allowing polymorphism, Go style.
type Instruction interface {
	Execute(ram *Memory, registers *Memory) error
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
	base.log = log.New(logOut, "Instruction Factory: ", 0)

	base.log.Printf("Decoding instruction: 0x%08x", instructionBits)

	// Set instruction bits
	base.InstructionBits = instructionBits

	// Get condition
	base.CondCode = ExtractBits(instructionBits, 28, 32) >> 28
	base.log.Printf("Condition bits: %04b", base.CondCode)

	// Get instruction type
	base.Type = ExtractBits(instructionBits, 25, 28) >> 25
	base.log.Printf("Type bits: %03b", base.Type)

	// Get Rn
	base.Rn = ExtractBits(instructionBits, 16, 20) >> 16
	base.log.Printf("Rn bits: %04b", base.Rn)

	// Get Rd
	base.Rd = ExtractBits(instructionBits, 12, 16) >> 12
	base.log.Printf("Rd bits: %04b", base.Rd)

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

	log *log.Logger
}

func (bi *baseInstruction) BuildFromBase() (instruction Instruction) {
	// Check edge cases

	// Check type of instruction and call proper decode method
	switch (bi.Type) {
	case 0x0, 0x1:
		log.Printf("Data Processing Instruction")
		instruction = new(dataInstruction)
	case 0x2:
		log.Printf("Load/Store: Immediate Offset")
		instruction = new(loadStoreInstruction)
	case 0x3:
		log.Printf("Load/Store: Register Offset")
		instruction = new(unimplementedInstruction)
	case 0x4:
		log.Printf("Load/Store: Multiple")
		instruction = new(unimplementedInstruction)
	case 0x5:
		log.Printf("Branch")
		instruction = new(branchInstruction)
	case 0xF:
		log.Printf("Software Interrupt")
		instruction = new(unimplementedInstruction)
	default:
		log.Printf("Unknown Instruction")
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
	// Second operand
	Operand2 uint32
}

// Executes a data instruction
//
// Parameters:
//  ram - a pointer to a block of memory, presumably ram
//  registers - a pointer to a block of memory, presumably a register bank
//
// Returns:
//  err - an error
func (di *dataInstruction) Execute(ram *Memory, registers *Memory) (err error) {
	di.log.SetPrefix("Instruction Executing: ")
	// Eventually this will contain logic to differentiate between dataInstruction
	// types
	// switch (opcode)
	//  case (. . .)

	// Assume MOV immediate
	registers.WriteWord(di.Rd << 2, di.Operand2)
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
	// Get opcode
	di.Opcode = byte(ExtractBits(di.InstructionBits, 21, 25) >> 21)
	di.log.Printf("Opcode bits: %04b", di.Opcode)

	// Get Operand2
	di.Operand2 = ExtractBits(di.InstructionBits, 0, 12)
	di.log.Printf("Op2 bits: %012b", di.Operand2)

	// Get S bit
	di.S = (ExtractBits(di.InstructionBits, 20, 21) >> 20) == 1
	di.log.Printf("S bit: %01b", di.S)

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
func (lsi *loadStoreInstruction) Execute(ram *Memory, registers *Memory) (err error) {
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
func (bi *branchInstruction) Execute(ram *Memory, registers *Memory) (err error) {
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
func (ui *unimplementedInstruction) Execute(ram *Memory, registers *Memory) (err error) {
	return
}

// Stub method to fake decoding of unimplemented instructions
func (ui *unimplementedInstruction) decode(base *baseInstruction) (err error) {
	return
}
