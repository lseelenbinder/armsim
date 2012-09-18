// Filename: cpu.go
// Contents: The CPU struct and related methods

package armsim

import (
	"log"
	"os"
	"time"
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
	log *log.Logger
}

// Initializes a CPU
//
// Parameters:
//  ram - a pointer to an initialized Memory struct
//  registers - a pointer to an initialized Memory struct with size of 64 bytes
//
// Returns:
//  a pointer to the newly created CPU
func NewCPU(ram *Memory, registers *Memory) (cpu *CPU) {
	cpu = new(CPU)

	cpu.log = log.New(os.Stdout, "CPU: ", 0)
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

// Decodes an instruction. (Currently does nothing.)
func (cpu *CPU) Decode() {
	// Does nothing; this is a stub.
	cpu.log.Println("Decoding...")
}

// Executes an instruction. (Currently pauses execution 0.25 seconds.)
func (cpu *CPU) Execute() {
	cpu.log.Println("Executing...waiting...")
	time.Sleep(time.Duration(250) * time.Millisecond)
}
