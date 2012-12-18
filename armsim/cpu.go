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
	CPSR // Current Program Status (N, Z, C, V or F flags)
	SPSR // Saved Program Status

	r13_svc  // Banked r13 for Supervisor
	r14_svc  // Banked r14 for Supervisor
	r13_irq  // Banked r13 for IRQ
	r14_irq  // Banked r14 for IRQ
	SPSR_svc // SPSR for Supervisor
	_
	SPSR_irq // SPSR for IRQ

	SP     = r13     // Stack Pointer
	SP_irq = r13_irq // Stack Pointer (IRQ mode)
	SP_svc = r13_svc // Stack Pointer (Supervisor mode)
	LR     = r14     // Link Register
	PC     = r15     // Program Counter

)

// Flags
const (
	_        = iota      // Ignore first result
	N uint32 = 32 - iota // Negative Flag
	Z                    // Zero Flag
	C                    // Carry Flag
	V                    // Overflow Flag
	F = V                // Overflow Flag (alternative spelling)
	I = 7                // Interrupt Bit
)

// Modes
const (
	User       = iota + 0x10 // PC, R14 to R0, CPSR
	_                        // FIQ (Not implemented)
	IRQ                      // PC, R14_irq, R13_irq, R12 to R0, CPSR, SPSR_irq
	Supervisor               // PC, R14_svc, R13_svc, R12 to R0, CPSR, SPSR_svc
	System     = 0x1F        // PC, R14 to R0, CPSR
)

// A CPU holds references for RAM and registers a CPU needs to function.
type CPU struct {
	// A reference to the assigned memory bank
	ram *Memory

	// A reference to the assigned registers bank
	registers *Memory

	// A channel for the Keyboard (since we don't have a true bus)
	keyboard chan byte

	// A channel for the Console (since we don't have a true bus)
	console chan byte

	// The IRQ pin
	irq chan bool

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
func NewCPU(ram *Memory, registers *Memory, keyboard chan byte,
	console chan byte, logOut io.Writer) (cpu *CPU) {
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

	// Assign Keyboard & Console
	cpu.keyboard = keyboard
	cpu.console = console

	// Setup IRQ
	cpu.irq = make(chan bool, 1)

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
	address, _ := cpu.registers.ReadWord(PC)
	instruction = Decode(cpu, address, instructionBits)
	return
}

// Executes an instruction.
//
// Parameter: i - Instruction interface
//
// Returns: status - bool determining if the CPU should continue executing
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
	value, err = cpu.registers.ReadWord(cpu.bankedRegister(r))

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
	return cpu.registers.WriteWord(cpu.bankedRegister(r), data)
}

// Wraps WriteRegister to allow it to use register index value obtained from an instruction.
// Parameters and return value are the same as WriteRegister.
func (cpu *CPU) WriteRegisterFromInstruction(r, data uint32) (err error) {
	return cpu.WriteRegister(r<<2, data)
}

// Wraps Memory.WriteByte to allow for memory-mapped IO
//
// Parameters:
//  address - 32-bit address of write location in memory
//  data - byte of data to write
//
// Returns:
//  err - any error that may have occurred
func (c *CPU) WriteOutByte(address uint32, data byte) (err error) {
	if address == 0x100001 {
		c.log.Printf("ERROR: Attempted to write to keyboard...")
	} else if address == 0x100000 {
		// add byte to console buffer
		c.console <- data
	} else {
		err = c.ram.WriteByte(address, data)
	}
	return
}

// Wraps Memory.ReadByte to allow for memory-mapped IO
//
// Parameters:
//  address - 32-bit address of read location in memory
//
// Returns:
//  data - byte of data at address
//  err - any error that may have occurred
func (c *CPU) ReadInByte(address uint32) (data byte, err error) {
	if address == 0x100000 {
		c.log.Printf("ERROR: Attempted to read from console...")
		return
	} else if address == 0x100001 {
		// read char from keyboard
		if len(c.keyboard) > 0 {
			data = byte(<-c.keyboard)
		} else {
			data = 0
		}
	} else {
		data, err = c.ram.ReadByte(address)
	}

	return
}

// Banked register locations
//
// Parameters:
//  r - a register
//
// Returns:
//	actualR - proper address for the register based on mode
func (cpu *CPU) bankedRegister(r uint32) (actualR uint32) {
	// Check for banked register
	if r == r13 || r == r14 || r == SPSR {
		cpsr, _ := cpu.FetchRegister(CPSR)
		mode := ExtractBits(cpsr, 0, 5)

		switch mode {
		case Supervisor:
			r += 5 << 2
			cpu.log.Printf("Using banked supervisor register %d...", r)
		case IRQ:
			cpu.log.Printf("Using banked IRQ register %d...", r)
			r += 7 << 2
		}
	}

	return r
}
