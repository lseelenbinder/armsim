// Filename: instructions.go
// Contents: The Instruction base struct and the other structs and methods
//	related to instructions.

package armsim

import (
	"fmt"
	"log"
	"strings"
)

// Implements a Go interface allowing Instruction polymorphism, Go style.
type Instruction interface {
	Execute() (status bool)
	Disassemble() string
	decode(base *baseInstruction)
}

// Holds values typical to all ARM instructions.
type baseInstruction struct {
	InstructionBits uint32 // The original bits of the instruction.
	Type            uint32 // Type bits
	CondCode        uint32 // Condition bits
	Rd              uint32 // Destination register
	Rn              uint32 // First register operand

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

	base.cpu = cpu // Set instruction's CPU

	base.InstructionBits = instructionBits // Set instruction bits
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

// Decodes a specific instruction from a baseInstruction.
//
// Returns an instruction interface.
func (bi *baseInstruction) BuildFromBase() (instruction Instruction) {
	// Check type of instruction and call proper decode method
	switch bi.Type {
	case 0x0, 0x1:
		// Check for BX
		if !(ExtractShiftBits(bi.InstructionBits, 21, 25) == 0x9) {
			bi.log.Printf("Data Processing")
			instruction = new(dataInstruction)
		} else {
			bi.log.Printf("Branch (BX)")
			instruction = new(branchInstruction)
		}
	case 0x2:
		bi.log.Printf("Load/Store: Immediate Offset")
		instruction = new(loadStoreInstruction)
	case 0x3:
		bi.log.Printf("Load/Store: Register Offset")
		instruction = new(loadStoreInstruction)
	case 0x4:
		bi.log.Printf("Load/Store: Multiple")
		instruction = new(loadStoreMultipleInstruction)
	case 0x5:
		bi.log.Printf("Branch")
		instruction = new(branchInstruction)
	case 0x7:
		bi.log.Printf("Software Interrupt")
		instruction = new(swiInstruction)
	default:
		bi.log.Printf("Unknown")
		instruction = new(unimplementedInstruction)
	}

	bi.log.SetPrefix("Instruction Decoding: ")
	instruction.decode(bi)

	return
}

// Holds values typical to a DataInstruction.
type dataInstruction struct {
	*baseInstruction // Embed a general instruction

	Opcode   byte   // Opcode
	S        bool   // S bit
	I        bool   // I bit
	Operand2 uint32 // Second operand
}

const (
	AND byte = 0x0  // 0000
	EOR      = 0x1  // 0001
	SUB      = 0x2  // 0010
	RSB      = 0x3  // 0011
	ADD      = 0x4  // 0100
	CMP      = 0xA  // 1010
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

	// Check for MUL
	if di.I == false && ExtractShiftBits(di.InstructionBits, 4, 5) == 1 && ExtractShiftBits(di.InstructionBits, 7, 8) == 1 {
		di.Opcode = MUL
	}

	// Get Operand2
	di.Operand2 = ExtractShiftBits(di.InstructionBits, 0, 12)
	di.log.Printf("Op2 bits: %012b", di.Operand2)

	// Parse the Operand2
	di.shifter = NewFromOperand2(di.Operand2, di.I, di.cpu)

	// Get S bit
	di.S = ExtractShiftBits(di.InstructionBits, 20, 21) == 1
	di.log.Printf("S bit: %01t", di.S)

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

	if !ConditionPassed(di.baseInstruction) {
		return true
	}

	shifter_operand := di.shifter.Shift()
	rn, _ := di.cpu.FetchRegisterFromInstruction(di.Rn)

	// Ascertain specific instruction
	switch di.Opcode {
	case MOV, MNV:
		// Negate for MNV
		if di.Opcode == MNV {
			shifter_operand ^= 0xFFFFFFFF
		}
		di.cpu.WriteRegisterFromInstruction(di.Rd, shifter_operand)
	case ADD:
		// Rd = Rn + shifter_operand
		di.cpu.WriteRegisterFromInstruction(di.Rd, rn+shifter_operand)
	case SUB:
		// Rd = Rn - shifter_operand
		di.cpu.WriteRegisterFromInstruction(di.Rd, rn-shifter_operand)
	case RSB:
		// Rd = shifter_operand - Rn
		di.cpu.WriteRegisterFromInstruction(di.Rd, shifter_operand-rn)
	case AND:
		// Rd = Rn AND shifter_operand
		di.cpu.WriteRegisterFromInstruction(di.Rd, rn&shifter_operand)
	case EOR:
		// Rd = Rn XOR shifter_operand
		di.cpu.WriteRegisterFromInstruction(di.Rd, rn^shifter_operand)
	case ORR:
		// Rd = Rn OR shifter_operand
		di.cpu.WriteRegisterFromInstruction(di.Rd, rn|shifter_operand)
	case BIC:
		// Rd = Rn AND NOT shifter_operand
		di.cpu.WriteRegisterFromInstruction(di.Rd, rn&^shifter_operand)
	case MUL:
		// Rd = Rm * Rs
		// This instruction is highly irregular, so the actual calculation is:
		// Rn = Rm * Rs
		di.cpu.WriteRegisterFromInstruction(di.Rn, di.shifter.GetRm()*di.shifter.GetRs())
	case CMP:
		// alu_out = Rn - shifter_operand
		alu_out := rn - shifter_operand

		// N Flag = alu_out[31]
		di.cpu.registers.SetFlag(CPSR, N, alu_out>>31 == 1)

		// Z Flag = if alu_out == 0 then 1 else 0
		di.cpu.registers.SetFlag(CPSR, Z, alu_out == 0)

		// C Flag = NOT BorrowFrom(Rn - shifter_operand)
		di.cpu.registers.SetFlag(CPSR, C, shifter_operand <= rn)

		// V Flag = OverflowFrom(Rn - shifter_operand)
		// Shift to sign bits and compare
		v := (shifter_operand>>31 != rn>>31) && (shifter_operand>>31 == alu_out>>31)
		di.cpu.registers.SetFlag(CPSR, V, v)

	default:
		di.log.Printf("Unknown Opcode: %04b", di.Opcode)
	}
	return true
}

// Builds an assembly string representing the instruction.
//
// Returns a string containing the mnemonic and related arguments.
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
	case CMP:
		assembly += "cmp"
	default:
		assembly += "unk"
	}

	assembly += ConditionMnemonic(di.CondCode)

	if di.Opcode == MUL {
		// Handle this special case
		assembly += fmt.Sprintf(" r%d, r%d, r%d", di.Rn, di.shifter.Rn, di.shifter.Rs)
	} else {
		assembly += fmt.Sprintf(" r%d, ", di.Rd)
		assembly += di.shifter.Disassemble()
	}

	return
}

// Holds values typical to Load / Store instructions.
type loadStoreInstruction struct {
	*baseInstruction // Embed a general instruction

	I bool // I bit
	P bool // P bit
	U bool // U bit
	B bool // B bit
	W bool // W bit
	L bool // L bit

	offset12 uint32         // Offset
	shifter  *BarrelShifter // Embed a shifter to handle operand2
}

// Decodes a load/store instruction
//
// Parameters:
//  base - a generic instruction containing most information
//
// Returns: None
func (lsi *loadStoreInstruction) decode(base *baseInstruction) {
	lsi.baseInstruction = base
	lsi.log.SetPrefix("Load/Store Decoder: ")
	// I bit
	lsi.I = ExtractShiftBits(base.InstructionBits, 25, 26) == 1
	lsi.log.Printf("I bit: %t", lsi.I)
	// P bit
	lsi.P = ExtractShiftBits(base.InstructionBits, 24, 25) == 1
	lsi.log.Printf("P bit: %t", lsi.P)
	// U bit
	lsi.U = ExtractShiftBits(base.InstructionBits, 23, 24) == 1
	lsi.log.Printf("U bit: %t", lsi.U)
	// B bit
	lsi.B = ExtractShiftBits(base.InstructionBits, 22, 23) == 1
	lsi.log.Printf("B bit: %t", lsi.B)
	// W bit
	lsi.W = ExtractShiftBits(base.InstructionBits, 21, 22) == 1
	lsi.log.Printf("W bit: %t", lsi.W)
	// L bit
	lsi.L = ExtractShiftBits(base.InstructionBits, 20, 21) == 1
	lsi.log.Printf("L bit: %t", lsi.L)

	// Offset
	op2 := ExtractShiftBits(base.InstructionBits, 0, 12)
	lsi.log.Printf("op2 bits: %#012b", op2)

	if !lsi.I {
		// Immediate
		lsi.offset12 = op2
		lsi.log.Printf("Immediate offset: %#012b", op2)
	} else {
		// I can take advantage of the BarrelShifter's logic
		lsi.shifter = NewFromOperand2(op2, false, lsi.cpu)
		lsi.log.Printf("Scaled offset: %#012b", lsi.shifter.Shift())
	}

	return
}

// Executes a load/store instruction
//
// Parameters: None
//
// Returns:
//  status - a boolean that determines if the CPU continues after this
//  instruction
func (lsi *loadStoreInstruction) Execute() (status bool) {
	if !ConditionPassed(lsi.baseInstruction) {
		return true
	}

	var address, base, offset, data uint32
	var data8 byte

	// Get base and offset
	base, _ = lsi.cpu.FetchRegisterFromInstruction(lsi.Rn)
	if !lsi.I {
		offset = lsi.offset12
	} else {
		offset = lsi.shifter.Shift()
	}

	// Pre-Index
	if lsi.P {
		address = lsi.calculateAddress(base, offset)
		lsi.log.Printf("Pre-Address: %#x", address)
	}

	// Load or Store
	if lsi.L {
		// Load
		if lsi.B {
			// Byte
			data8, _ = lsi.cpu.ram.ReadByte(address)
			data = uint32(data8)
		} else {
			// Word
			data, _ = lsi.cpu.ram.ReadWord(address)
		}

		// Write to register
		lsi.cpu.WriteRegisterFromInstruction(lsi.Rd, data)
	} else {
		// Store
		data, _ = lsi.cpu.FetchRegisterFromInstruction(lsi.Rd)

		if lsi.B {
			// Byte
			data8 = byte(data)
			// Write to memory
			lsi.cpu.ram.WriteByte(address, data8)
		} else {
			// Write to memory
			lsi.cpu.ram.WriteWord(address, data)
		}
	}

	// Post-Index
	if !lsi.P {
		address = lsi.calculateAddress(base, offset)
		lsi.log.Printf("Post-Address: %#x", address)
	}

	// Writeback
	if lsi.W {
		lsi.cpu.WriteRegisterFromInstruction(lsi.Rn, address)
		lsi.log.Printf("Write-back: %#d = %#x", lsi.Rn, address)
	}

	return true
}

// Builds an assembly string representing the instruction.
//
// Returns a string containing the mnemonic and related arguments.
func (lsi *loadStoreInstruction) Disassemble() (assembly string) {
	var mnemonic, arguments, writeback, shift string

	if lsi.L {
		mnemonic += "ldr"
	} else {
		mnemonic += "str"
	}

	if lsi.B {
		mnemonic += "b"
	}

	mnemonic += ConditionMnemonic(lsi.CondCode)

	if !lsi.I {
		shift = fmt.Sprintf("#%d", lsi.offset12)
	} else {
		shift = lsi.shifter.Disassemble()
	}

	arguments = fmt.Sprintf("r%d", lsi.Rn)
	if shift != "#0" {
		arguments += ", "
		if !lsi.U {
			arguments += "-"
		}
		arguments += fmt.Sprintf("%s", shift)
	}

	if lsi.W {
		writeback = "!"
	}

	return fmt.Sprintf("%s r%d, [%s] %s", mnemonic, lsi.Rd, arguments, writeback)
}

// Calulates an effective address based on a base address, an offset, and the U
// bit
func (lsi *loadStoreInstruction) calculateAddress(base, offset uint32) (address uint32) {
	if lsi.U {
		address = base + offset
	} else {
		address = base - offset
	}

	return
}

// Holds values typical to Load / Store Multiple instructions.
type loadStoreMultipleInstruction struct {
	*baseInstruction // Embed a general instruction

	P bool // P bit
	U bool // U bit
	S bool // S bit
	W bool // W bit
	L bool // L bit

	registerList [16]bool // Registers
}

// Decodes a load/store multiple instruction
//
// Parameters:
//  base - a generic instruction containing most information
//
// Returns: None
func (lsi *loadStoreMultipleInstruction) decode(base *baseInstruction) {
	lsi.baseInstruction = base
	lsi.log.SetPrefix("Load/Store Decoder: ")
	// P bit
	lsi.P = ExtractShiftBits(base.InstructionBits, 24, 25) == 1
	lsi.log.Printf("P bit: %t", lsi.P)
	// U bit
	lsi.U = ExtractShiftBits(base.InstructionBits, 23, 24) == 1
	lsi.log.Printf("U bit: %t", lsi.U)
	// S bit
	lsi.S = ExtractShiftBits(base.InstructionBits, 22, 23) == 1
	lsi.log.Printf("S bit: %t", lsi.S)
	// W bit
	lsi.W = ExtractShiftBits(base.InstructionBits, 21, 22) == 1
	lsi.log.Printf("W bit: %t", lsi.W)
	// L bit
	lsi.L = ExtractShiftBits(base.InstructionBits, 20, 21) == 1
	lsi.log.Printf("L bit: %t", lsi.L)

	for i := 0; i < 16; i++ {
		lsi.registerList[i] = ExtractShiftBits(base.InstructionBits, uint32(i), uint32(i+1)) == 1
	}
	lsi.log.Println("Registers: ", lsi.registerList)

	return
}

// Executes a load/store multiple instruction
//
// Parameters: None
//
// Returns:
//  status - a boolean that determines if the CPU continues after this
//  instruction
func (lsi *loadStoreMultipleInstruction) Execute() (status bool) {
	if !ConditionPassed(lsi.baseInstruction) {
		return true
	}
	var address, start_address, end_address, data uint32
	Rn, _ := lsi.cpu.FetchRegisterFromInstruction(lsi.Rn)

	if lsi.P {
		if lsi.U { // Increment before
			start_address = Rn + 4
			end_address = Rn + (lsi.CountSetBits() * 4)
			Rn += lsi.CountSetBits() * 4
		} else { // Decrement before
			start_address = Rn - (lsi.CountSetBits() * 4)
			end_address = Rn - 4
			Rn -= lsi.CountSetBits() * 4
		}
	} else {
		if lsi.U { // Increment after
			start_address = Rn
			end_address = Rn + (lsi.CountSetBits() * 4) - 4
			Rn += lsi.CountSetBits() * 4
		} else { // Decrement after
			start_address = Rn - (lsi.CountSetBits() * 4) + 4
			end_address = Rn
			Rn -= lsi.CountSetBits() * 4
		}
	}
	lsi.log.Printf("start_address: %#x; end_address: %#x", start_address, end_address)

	address = start_address
	for i := 0; i < 16; i++ {
		if lsi.registerList[i] {
			if lsi.L { // Load
				data, _ = lsi.cpu.ram.ReadWord(address)
				lsi.cpu.WriteRegisterFromInstruction(uint32(i), data)
			} else { // Store
				data, _ = lsi.cpu.FetchRegisterFromInstruction(uint32(i))
				lsi.cpu.ram.WriteWord(address, data)
			}
			address += 4
		}
	}

	if lsi.W { // Writeback
		lsi.cpu.WriteRegister(SP, Rn)
	}

	return true
}

// Builds an assembly string representing the instruction.
//
// Returns a string containing the mnemonic and related arguments.
func (lsi *loadStoreMultipleInstruction) Disassemble() (assembly string) {
	var mnemonic, registers, rn string

	if lsi.L {
		mnemonic += "ldm"
	} else {
		mnemonic += "stm"
	}

	if lsi.P {
		if lsi.U {
			mnemonic += "ib"
		} else {
			mnemonic += "db"
		}
	} else {
		if lsi.U {
			mnemonic += "ia"
		} else {
			mnemonic += "da"
		}
	}

	mnemonic += ConditionMnemonic(lsi.CondCode)

	for i := 0; i < 16; i++ {
		if lsi.registerList[i] {
			registers += fmt.Sprintf("r%d ", i)
		}
	}
	registers = strings.TrimSpace(registers)

	rn = fmt.Sprintf("r%d", lsi.Rn)
	if lsi.W {
		rn += "!"
	}

	assembly = fmt.Sprintf("%s %s, {%s}", mnemonic, rn, registers)

	return
}

func (lsi *loadStoreMultipleInstruction) CountSetBits() (count uint32) {
	for i := 0; i < 16; i++ {
		if lsi.registerList[i] {
			count++
		}
	}
	return
}

// Holds values typical to Branch instructions.
type branchInstruction struct {
	// Embedding a general instruction
	*baseInstruction

	// BX
	bx bool

	// L bit
	L bool

	// Rm (for BX)
	Rm uint32

	// Offset bits
	Offset int32
}

// Decodes a branch instruction
//
// Parameters:
//  base - a generic instruction containing most information
//
// Returns: None
func (bi *branchInstruction) decode(base *baseInstruction) {
	bi.baseInstruction = base
	bi.log.SetPrefix("Branch Instruction (Decode): ")

	// Check for BX
	if bi.Type == 0x5 {
		// B or BL
		bi.bx = false

		// Link bit
		bi.L = ExtractShiftBits(bi.InstructionBits, 24, 25) == 1
		bi.log.Printf("L bit: %01t", bi.L)

		// Offset is a 24-bit signed number, I need to sign extend to 32 and then
		// shift right 6 places (to account for multiplication by 4)
		bi.Offset = int32(ExtractBits(bi.InstructionBits, 0, 24)<<8) >> 6
		bi.log.Printf("Offset: %d", bi.Offset)
	} else {
		// BX
		bi.bx = true

		// Set L bit
		bi.L = false

		// Rm
		bi.Rm = ExtractShiftBits(bi.InstructionBits, 0, 4)
		bi.log.Printf("Rm: %d", bi.Rm)
	}

	return
}

// Executes a branch instruction
//
// Parameters: None
//
// Returns:
//  status - a boolean that determins if the CPU continues after this
//  instruction
func (bi *branchInstruction) Execute() (status bool) {
	bi.log.SetPrefix("Branch Instruction (Execute): ")

	// Check condition
	if !ConditionPassed(bi.baseInstruction) {
		// Don't continue execution
		return true
	}

	if bi.L {
		pc, _ := bi.cpu.FetchRegister(PC)
		bi.cpu.WriteRegister(LR, pc-4)
	}

	// Check for BX
	var newPC uint32
	if !bi.bx {
		// B or BL
		pc, _ := bi.cpu.FetchRegister(PC)
		newPC = uint32(int32(pc) + bi.Offset)
	} else {
		newPC, _ = bi.cpu.FetchRegisterFromInstruction(bi.Rm)
		newPC &= 0xFFFFFFFE
	}

	bi.log.Printf("Branching to %X...", newPC)
	bi.cpu.WriteRegister(PC, newPC)

	return true
}

// Builds an assembly string representing the instruction.
//
// Returns a string containing the mnemonic and related arguments.
func (bi *branchInstruction) Disassemble() (assembly string) {
	assembly = "b"

	if bi.bx {
		assembly += "x"
	}
	if bi.L {
		assembly += "l"
	}

	assembly += ConditionMnemonic(bi.CondCode)

	if bi.bx {
		assembly += fmt.Sprintf(" r%d", bi.Rm)
	} else {
		// TODO: This will not be correct out-of-context
		pc, _ := bi.cpu.FetchRegister(PC)
		newPC := uint32(int32(pc) + bi.Offset)
		assembly += fmt.Sprintf(" #%X", newPC)
	}

	return
}

// Holds values typical to Software Interrupt instructions.
type swiInstruction struct {
	*baseInstruction // Embed a general instruction
	Data             uint32
}

// Decodes an interrupt instruction
//
// Parameters:
//  base - a generic instruction containing most information
//
// Returns: None
func (swi *swiInstruction) decode(base *baseInstruction) {
	swi.baseInstruction = base

	swi.Data = ExtractBits(base.InstructionBits, 0, 24)
	swi.log.Printf("Immediate 24 bits: %0#x", swi.Data)
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
	if !ConditionPassed(si.baseInstruction) {
		return true
	}

	// All SWI instructions immediately stop execution at this point
	return false
}

// Builds an assembly string representing the instruction.
//
// Returns a string containing the mnemonic and related arguments.
func (swi *swiInstruction) Disassemble() (assembly string) {
	return fmt.Sprintf("swi #%d", swi.Data)
}

// Stub class to fake unknown or unimplemented instructions
type unimplementedInstruction struct {
	*baseInstruction
}

// Stub method to fake execution of unimplemented instructions
func (ui *unimplementedInstruction) Execute() (status bool) {
	if !ConditionPassed(ui.baseInstruction) {
		return true
	}

	return true
}

// Stub method to fake decoding of unimplemented instructions
func (ui *unimplementedInstruction) decode(base *baseInstruction) {
	ui.baseInstruction = base
	return
}

// Stub method to fake disassemble of unimplemented instructions
func (ui *unimplementedInstruction) Disassemble() (assembly string) { return "unk" }

// Conditions
const (
	EQ  = iota // Equal
	NE         // Not equal
	CS         // Carry set / unsigned higher or same
	CC         // Carry clear/unsigned lower
	MI         // Minus/negative
	PL         // Plus/positive or zero
	VS         // Overflow
	VC         // No overflow
	HI         // Unsigned higher
	LS         // Unsigned lower or same
	GE         // Signed greater than or equal
	LT         // Signed less than
	GT         // Signed greater than
	LE         // Signed less than or equal
	AL         // Always (unconditional)
	UNP        // Unpredictable
)

// Condition Passed implements the logic to test the flags for conditional execution
//
// Parameters:
//  bi - the baseInstruction that contains condition codes and other necessary
//  information
//
// Returns:
//  passed - a bool, true if the condition passed, otherwise false
func ConditionPassed(bi *baseInstruction) (passed bool) {
	// Fetch flags
	z, _ := bi.cpu.registers.TestFlag(CPSR, Z)
	c, _ := bi.cpu.registers.TestFlag(CPSR, C)
	n, _ := bi.cpu.registers.TestFlag(CPSR, N)
	v, _ := bi.cpu.registers.TestFlag(CPSR, V)

	switch bi.CondCode {
	case EQ:
		// Z == 1
		passed = z
	case NE:
		// Z == 0
		passed = !z
	case CS:
		// C == 1
		passed = c
	case CC:
		// C == 0
		passed = !c
	case MI:
		// N == 1
		passed = n
	case PL:
		// N == 0
		passed = !n
	case VS:
		// V == 1
		passed = v
	case VC:
		// V == 0
		passed = !v
	case HI:
		// C set and Z clear
		passed = c && !z
	case LS:
		// C clear or Z set
		passed = !c || z
	case GE:
		// N == V
		passed = n == v
	case LT:
		// N != V
		passed = n != v
	case GT:
		// Z == 0 and N == V
		passed = !z && n == v
	case LE:
		// Z == 1 or N != V
		passed = z || n != v
	case AL:
		// Always
		passed = true
	case UNP:
		bi.log.Printf("Unpredictable condition code...ignoring.")
		passed = true
	default:
		passed = true
	}

	return
}

// Condition Mnemonic returns the approriate mnemonic for a condition code
func ConditionMnemonic(cond uint32) (m string) {
	switch cond {
	case EQ:
		m = "eq"
	case NE:
		m = "ne"
	case CS:
		m = "cs"
	case CC:
		m = "cc"
	case MI:
		m = "mi"
	case PL:
		m = "pl"
	case VS:
		m = "vs"
	case VC:
		m = "vc"
	case HI:
		m = "hi"
	case LS:
		m = "ls"
	case GE:
		m = "ge"
	case LT:
		m = "lt"
	case LE:
		m = "le"
	case AL, UNP:
		m = ""
	}

	return
}
