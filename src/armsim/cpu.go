// Filename: cpu.go
// Contents: The CPU struct and related methods

package armsim

import (
	"io"
	"log"
	"os"
)

// Registers
const (
	// Note: In a normal ARM architecture, there would be 32 registers. Two
	// 15-place banks, a CPSR, and a SPSR. In this implementation, the dual
	// "normal" register banks are not implemented, nor is the SPSR implemented.

	// iota is an enumerator that returns 0, 1, 2, ..., on successive calls.
	// Each consecutive constant is assumed to be equal to the previous constant;
	// however, because of iota this means they are created correctly.
	r0 uint32 = 4 * iota
	r1
	r2
	r3
	r4
	r5
	r6
	r7
	r8
	r9
	r10
	r11
	r12
	r13
	r14
	r15
	r16
	// Stack Pointer
	SP = r13
	// Link Register
	LR = r14
	// Program Counter
	PC = r15
	// Current Program Status (N, Z, C, V or F flags)
	CPSR = r16
)

// Flags
const (
	// Ignore first result
	_ = iota
	// Negative Flag
	N uint32 = 32 - iota
	// Zero Flag
	Z
	// Carry Flag
	C
	// Overflow Flag
	V
	// Overflow Flag (alternative spelling)
	F = V
)

// A CPU holds references for RAM and registers a CPU needs to function.
type CPU struct {
	// A reference to the assigned memory bank
	ram *Memory

	// A reference to the assigned registers bank
	registers *Memory

	// Logging class
	log    *log.Logger
	logOut io.Writer
}

// Initializes a CPU
//
// Parameters:
//  ram - a pointer to an initialized Memory struct
//  registers - a pointer to an initialized Memory struct with size of 64 bytes
//  logOut - an io.Writer out stream for the logger to use (or nil to use StdErr)
//
// Returns:
//  a pointer to the newly created CPU
func NewCPU(ram *Memory, registers *Memory, logOut io.Writer) (cpu *CPU) {
	cpu = new(CPU)

	if logOut == nil {
		logOut = os.Stderr
	}
	cpu.log = log.New(logOut, "CPU: ", 0)
	cpu.logOut = logOut
	cpu.log.Println("Created new CPU.")

	// Assign RAM
	cpu.ram = ram
	cpu.log.Println("Assigned RAM @", &ram)

	// Assign Registers
	cpu.registers = registers
	cpu.log.Println("Assigned registers @", &registers)

	return
}

// Fetches the next instruction and increments the program counter
//
// Returns:
//  encoded instruction - 32-bit unsigned integer (i.e., a word)
func (cpu *CPU) Fetch() (instruction uint32) {
	// Read address stored in the PC
	address, err := cpu.registers.ReadWord(PC)
	if err != nil {
		cpu.log.Panic("Unable to read PC.")
	}
	cpu.log.Printf("Current PC: %#x", address)

	// Read instruction stored at address
	instruction, err = cpu.ram.ReadWord(address)
	if err != nil {
		cpu.log.Panic("Unable to read next instruction.")
	}
	cpu.log.Printf("Instruction fetched: %#x", instruction)

	// Increment PC
	cpu.registers.WriteWord(PC, address+4)
	return
}

// Decodes an instruction.
//
// Parameters:
//	instructionBits - word of data representing the next instruction
//
// Returns:
//	instruction - a decoded instruction of type Instruction
func (cpu *CPU) Decode(instructionBits uint32) (instruction Instruction) {
	cpu.log.Println("Decoding...")
	instruction = Decode(cpu, instructionBits)
	return
}

// Executes an instruction. (Currently pauses execution 0.25 seconds.)
func (cpu *CPU) Execute(i Instruction) (status bool) {
	cpu.log.Println("Executing...", i.Disassemble())
	return i.Execute()
}

// Fetches a register's value. This function accounts for the fact PC should be R[PC] + 8
//
// Parameters:
//  r - register (equal to one of the constants defined above)
//
// Returns:
//  value - the register's content
func (cpu *CPU) FetchRegister(r uint32) (value uint32, err error) {
	value, err = cpu.registers.ReadWord(r)

	// Because of pipelining, any PC access will need to be +8. However, the PC is already
	// incremented, so it will only be +4.
	if r == PC {
		value += 4
	}
	return
}

// Wraps FetchRegister to allow a register index value obtained from an instruction.
// Parameters and return value are the same as FetchRegister.
func (cpu *CPU) FetchRegisterFromInstruction(r uint32) (value uint32, err error) {
	return cpu.FetchRegister(r << 2)
}

// Writes a register's value.
//
// Parameters:
//  r - register (equal to one of the constants defined above)
//  data - a word of data to write to a register
//
// Returns:
//  err - any error that may have occured
func (cpu *CPU) WriteRegister(r, data uint32) (err error) {
	return cpu.registers.WriteWord(r, data)
}

// Wraps WriteRegister to allow it to use register index value obtained from an instruction.
// Parameters and return value are the same as WriteRegister.
func (cpu *CPU) WriteRegisterFromInstruction(r, data uint32) (err error) {
	return cpu.WriteRegister(r<<2, data)
}
