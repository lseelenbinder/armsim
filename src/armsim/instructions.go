// Filename: instructions.go
// Contents: The Instruction base struct and the other structs and methods
//	related to instructions.

package armsim

import "log"

// Implements a Go interface allowing polymorphism, Go style.
type Instruction interface {
	Execute(ram *Memory, registers *Memory) error
}

// Decodes an instruction.
//
// Parameters:
//	instructionBits - word of data representing the next instruction
//
// Returns:
//	instruction - a decoded instruction of type Instruction
func Decode(instructionBits uint32) (instruction Instruction) {
	log.Println("Decoding instruction: ", instructionBits)

	// Get condition
	condBits := ExtractBits(instructionBits, 28, 32) >> 28
	log.Printf("Condition bits: %04b", condBits)

	// Get instruction type
	typeBits := ExtractBits(instructionBits, 25, 28) >> 25
	log.Printf("Type bits: %03b", typeBits)

	// Get Rn
	rnBits := ExtractBits(instructionBits, 16, 20) >> 16
	log.Printf("Rn bits: %04b", rnBits)

	// Get Rd
	rdBits := ExtractBits(instructionBits, 12, 16) >> 12
	log.Printf("Rd bits: %04b", rdBits)

	// Initalize generic instruction
	base := baseInstruction{instructionBits, byte(condBits), rdBits,
		rnBits}

	// Check edge cases

	// Check type of instruction and call proper decode method
	// switch(typeBits)

	// Assume a data instruction for now
	di := new(dataInstruction)
	di.decode(&base)
	instruction = di

	return
}

type baseInstruction struct {
	// The original bits of the instruction.
	InstructionBits uint32

	// Condition bits
	CondCode byte
	// Destination register
	Rd uint32
	// First register operand
	Rn uint32
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
	log.Printf("Opcode bits: %04b", di.Opcode)

	// Get Operand2
	di.Operand2 = ExtractBits(di.InstructionBits, 0, 12)
	log.Printf("Op2 bits: %012b", di.Operand2)

	// Get S bit
	di.S = (ExtractBits(di.InstructionBits, 20, 21) >> 20) == 1
	log.Printf("S bit: %01b", di.S)

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
