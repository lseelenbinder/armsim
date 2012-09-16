// Filename: cpu.go
// Contents: The CPU struct and related methods

package armsim

import (
	"log"
	"time"
)

// Registers
const (
	// iota is an enumerator that returns 0, 1, 2, ... on successive calls.
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
	// Stack Pointer
	SP = r13
	// Link Register
	LR = r14
	// Program Counter
	PC = r15
)

// A CPU holds references for RAM and registers a CPU needs to function.
type CPU struct {
	// A reference to the assigned memory bank
	ram *Memory

	// A reference to the assigned registers bank
	registers *Memory
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
	log.SetPrefix("CPU: ")

	cpu = new(CPU)
	log.Println("Created new CPU.")

	cpu.ram = ram
	log.Println("Assigned RAM @", &ram)

	cpu.registers = registers
	log.Println("Assigned registers @", &registers)

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
		log.Panic("Unable to read PC.")
	}
	log.Printf("Current PC: %#x", address)

	// Read instruction stored at address
	instruction, err = cpu.ram.ReadWord(address)
	if err != nil {
		log.Panic("Unable to read next instruction.")
	}

	// Increment PC
	cpu.registers.WriteWord(PC, address+4)
	return
}

// Decodes an instruction. (Currently does nothing.)
func (cpu *CPU) Decode() {
	// Does nothing; this is a stub.
}

// Executes an instruction. (Currently pauses execution 0.25 seconds.)
func (cpu *CPU) Execute() {
	time.Sleep(time.Duration(250) * time.Millisecond)
}
