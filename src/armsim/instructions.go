// Filename: instructions.go
// Contents: The Instruction base struct and the other structs and methods
//	related to instructions.

package armsim

import (
	"fmt"
	"log"
)

// Implements a Go interface allowing polymorphism, Go style.
type Instruction interface {
	Execute() (status bool)
	Disassemble() string
	decode(base *baseInstruction)
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
	cpu     *CPU
}

// Decodes an instruction.
//
// Parameters:
//	instructionBits - word of data representing the next instruction
//
// Returns:
//	instruction - a decoded instruction of type Instruction
func Decode(cpu *CPU, instructionBits uint32) (instruction Instruction) {
	base := new(baseInstruction)
	base.log = log.New(cpu.logOut, "Instruction Factory: ", 0)

	base.log.Printf("Decoding instruction: 0x%08x", instructionBits)

	// Set instruction's CPU
	base.cpu = cpu

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
	case 0x7:
		bi.log.Printf("Software Interrupt")
		instruction = new(swiInstruction)
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
	AND byte = 0x0  // 0000
	EOR      = 0x1  // 0001
	SUB      = 0x2  // 0010
	RSB      = 0x3  // 0011
	ADD      = 0x4  // 0100
	ORR      = 0xC  // 1100
	BIC      = 0xE  // 1110
	MOV      = 0xD  // 1101
	MNV      = 0xF  // 1111
	MUL      = 0x10 // 1 0000 (custom opcode)
)

// Decodes a data instruction
//
// Parameters:
//  base - a generic instruction containing most information
//
// Returns: None
func (di *dataInstruction) decode(base *baseInstruction) {
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

	// Parse the Operand2
	di.shifter = NewFromOperand2(di.Operand2, di.I, di.cpu)

	// Get S bit
	di.S = ExtractShiftBits(di.InstructionBits, 20, 21) == 1
	di.log.Printf("S bit: %01t", di.S)

	// Check for MUL
	if di.I == false && ExtractShiftBits(di.InstructionBits, 4, 5) == 1 && ExtractShiftBits(di.InstructionBits, 7, 8) == 1 {
		di.Opcode = MUL
	}

	di.log.Printf("Decoded: %s", di.Disassemble())

	return
}

// Executes a data instruction
//
// Parameters:
//  ram - a pointer to a block of memory, presumably ram
//  registers - a pointer to a block of memory, presumably a register bank
//
// Returns:
//  err - an error
func (di *dataInstruction) Execute() (status bool) {
	di.log.SetPrefix("Data Instruction (Execute): ")

	result := di.shifter.Shift()
	rn, _ := di.cpu.FetchRegisterFromInstruction(di.Rn)

	// Assertain specific instruction

	switch di.Opcode {
	case MOV, MNV:
		// Negate for MNV
		if di.Opcode == MNV {
			result ^= 0xFFFFFFFF
		}
		di.cpu.WriteRegisterFromInstruction(di.Rd, result)
	case ADD:
		// Rd = Rn + shifter_operand
		di.cpu.WriteRegisterFromInstruction(di.Rd, rn+result)
	case SUB:
		// Rd = Rn - shifter_operand
		di.cpu.WriteRegisterFromInstruction(di.Rd, rn-result)
	case RSB:
		// Rd = shifter_operand - Rn
		di.cpu.WriteRegisterFromInstruction(di.Rd, result-rn)
	case AND:
		// Rd = Rn AND shifter_operand
		di.cpu.WriteRegisterFromInstruction(di.Rd, rn&result)
	case EOR:
		// Rd = Rn XOR shifter_operand
		di.cpu.WriteRegisterFromInstruction(di.Rd, rn^result)
	case ORR:
		// Rd = Rn OR shifter_operand
		di.cpu.WriteRegisterFromInstruction(di.Rd, rn|result)
	case BIC:
		// Rd = Rn AND NOT shifter_operand
		di.cpu.WriteRegisterFromInstruction(di.Rd, rn&^result)
	case MUL:
		// Rd = Rm * Rs
		// This instruction is highly irregular, so the actual calculation is:
		// Rn = Rm * Rs
		di.cpu.WriteRegisterFromInstruction(di.Rn, di.shifter.GetRm()*di.shifter.GetRs())
	default:
		di.log.Printf("Unknown. Opcode: %04b", di.Opcode)
	}
	return true
}

func (di *dataInstruction) Disassemble() (assembly string) {
	// Get the Opcode
	switch di.Opcode {
	case MOV:
		assembly += "mov"
	case MNV:
		assembly += "mnv"
	case ADD:
		assembly += "add"
	case SUB:
		assembly += "sub"
	case RSB:
		assembly += "rsb"
	case AND:
		assembly += "and"
	case EOR:
		assembly += "eor"
	case ORR:
		assembly += "orr"
	case BIC:
		assembly += "bic"
	case MUL:
		assembly += "mul"
	default:
		assembly += "unk"
	}

	if di.Opcode == MUL {
		// Handle this special case
		assembly += fmt.Sprintf(" r%d, r%d, r%d", di.Rn, di.shifter.Rn, di.shifter.Rs)
	} else {
		assembly += fmt.Sprintf(" r%d, ", di.Rd)
		assembly += di.shifter.Disassemble()
	}

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
// Parameters: None
//
// Returns:
//  status - a boolean that determines in the CPU statuss after this
//  instruction
func (lsi *loadStoreInstruction) Execute() (status bool) {
	// Stub
	return
}

// Decodes a load/store instruction
//
// Parameters:
//  base - a generic instruction containing most information
//
// Returns: None
func (lsi *loadStoreInstruction) decode(base *baseInstruction) {
	// Stub
	return
}

func (lsi *loadStoreInstruction) Disassemble() (assembly string) { return }

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
func (bi *branchInstruction) Execute() (status bool) {
	// Stub
	return
}

// Decodes a branch instruction
//
// Parameters:
//  base - a generic instruction containing most information
//
// Returns: None
func (bi *branchInstruction) decode(base *baseInstruction) {
	// Stub
	return
}

func (bi *branchInstruction) Disassemble() (assembly string) { return }

type swiInstruction struct {
	Data uint32
}

func (swi *swiInstruction) decode(base *baseInstruction) {
	// I don't even really need this data, but the OS might.
	swi.Data = ExtractBits(base.InstructionBits, 0, 25)
	return
}

// Executes a software interrupt instruction
//
// Parameters: None
//
// Returns:
//  status - a boolean that determines in the CPU statuss after this
//  instruction
func (si *swiInstruction) Execute() (status bool) {
	// All SWI instructions immediately stop execution at this point
	return false
}

func (swi *swiInstruction) Disassemble() (assembly string) { return }

type unimplementedInstruction struct {
}

// Stub method to fake execution of unimplemented instructions
func (ui *unimplementedInstruction) Execute() (status bool) {
	return
}

// Stub method to fake decoding of unimplemented instructions
func (ui *unimplementedInstruction) decode(base *baseInstruction) {
	return
}

func (ui *unimplementedInstruction) Disassemble() (assembly string) { return }
