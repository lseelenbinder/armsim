// Filename: computer.go
// Contents: The Computer struct and related methods.

package armsim

import (
	"debug/elf"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

// A Computer holds the RAM, registers, and CPU of the simulated ARM
// architecture. It has methods to allow the loading and execution of a program
// on the simulator.
type Computer struct {
	// A reference to RAM for the simulator
	ram *Memory

	// A reference to the bank of CPU registers,
	// implemented using a standard memory container
	memSize   uint32
	registers *Memory

	// A reference to the CPU for the simulator
	cpu *CPU

	// A simple counter to track number of execution cycles
	step_counter uint64

	// Logger class
	log *log.Logger

	// Trace Log File
	traceFile *os.File

	// Keyboard buffer
	Keyboard chan byte
	// Console buffer
	Console chan byte
}

// A ComputerStatus is an individual module designed to make it easy to pass
// status information to external code.
type ComputerStatus struct {
	Flags       [4]bool    // CPSR Flags
	Disassembly []string   // The 2 previous instructions, current instruction, and next 7 instructions
	Registers   [16]uint32 // A representation of the registers
	Stack       []uint32   // A representation of the top of the stack
	Memory      []string   // A string representation of the RAM
	Steps       uint64     // The number of steps executed so far (step_counter)
	Checksum    int32      // Current RAM Checksum
}

// Initializes a Computer
//
// Parameters:
//  memSize - a unsigned 32-bit integer specifying the size of the RAM in the
//  computer
//  logOut - an io.Writer out stream for the logger to use (or nil to use StdErr)
//
// Returns:
//  a pointer to the newly created Computer
func NewComputer(memSize uint32, logOut io.Writer) (c *Computer) {
	c = new(Computer)

	// Setup logging
	if logOut == nil {
		logOut = os.Stderr
	}
	c.log = log.New(logOut, "Computer: ", 0)

	// Initialize RAM of memSize
	c.memSize = memSize
	c.ram = NewMemory(memSize, logOut)

	// Initialize a register bank to contain all 16 registers + CPSR + Banked
	// registers
	c.registers = NewMemory(SPSR_irq+4, logOut)

	// Initialize buffers
	c.Keyboard = make(chan byte, 100)
	c.Console = make(chan byte, 100)

	// Initialize CPU with RAM and registers
	c.cpu = NewCPU(c.ram, c.registers, c.Keyboard, c.Console, logOut)

	// Trace Log File
	if err := c.EnableTracing(); err != nil {
		c.log.Println("Unable to open trace file -", err)
	}

	// Step Counter
	c.step_counter = 1

	return
}

// Simulates the running of the a computer. It executes the fetch, execute,
// decode cycle until fetch returns false (signifying an instruction of 0x0).
//
// Parameters:
//  haltng - channel to enable midstream halting of running (for Stop/Break in gui)
//  finishing - channel to allow caller to know when Run() is finished
func (c *Computer) Run(halting, finishing chan bool) {
	var h bool
	for {
		if len(halting) > 0 {
			h = <-halting
			if h {
				break
			}
		}

		if !c.Step() {
			break
		}
	}

	// Let caller know we are finished
	if finishing != nil {
		finishing <- true
	}
}

// Builds and returns a the status of the emulator via a ComputerStatus
//
// Parameters: None
//
// Returns:
//  status - a ComputerStatus module fully initialized
func (c *Computer) Status() (status ComputerStatus) {
	status.Flags[0], _ = c.registers.TestFlag(CPSR, N) // Negative Flag
	status.Flags[1], _ = c.registers.TestFlag(CPSR, Z) // Zero Flag
	status.Flags[2], _ = c.registers.TestFlag(CPSR, C) // Carry Flag
	status.Flags[3], _ = c.registers.TestFlag(CPSR, F) // Overflow Flag
	c.log.Println("Flags:", status.Flags)

	for i := 0; i < 16; i++ {
		status.Registers[i], _ = c.registers.ReadWord(uint32(i << 2))
	}

	sp, _ := c.cpu.FetchRegister(SP)
	stop := uint32(0x7000)
	stackSize := (stop-sp)/4 + 1
	if stackSize > 5 {
		stop = sp + 0x4*5
		stackSize = 5
	}
	status.Stack = make([]uint32, stackSize)
	c.log.Printf("Size of Stack: %#x (sp: %#x)", stackSize, sp)
	var i uint32
	for ; sp < stop; sp += 0x4 {
		status.Stack[i], _ = c.ram.ReadWord(sp)
		i++
	}

	status.Disassembly = make([]string, 8)
	i = 0
	for pc := status.Registers[15] - 8; pc < status.Registers[15]+4*6; pc += 4 {
		iBits, _ := c.ram.ReadWord(pc)
		instruction := Decode(c.cpu, iBits)
		status.Disassembly[i] = fmt.Sprintf("%x||%s", iBits, instruction.Disassemble())
		i++
	}

	status.Memory = make([]string, c.memSize)
	for i = 0; i < c.memSize; i++ {
		b, _ := c.ram.ReadByte(i)
		status.Memory[i] = fmt.Sprintf("%x", b)
	}

	status.Steps = c.step_counter
	status.Checksum = c.Checksum()

	return
}

// Performs a single execution cycle. Take no parameters and returns a boolean
// signifying if the cycle was completed (a cycle will not complete if the
// instrution fetched is 0x0).
func (c *Computer) Step() (status bool) {
	// For trace
	pc, _ := c.cpu.FetchRegister(PC)
	instructionBits := c.cpu.Fetch()

	// Don't continue if the instruction is useless
	if instructionBits == 0x0 {
		return false
	}

	instruction := c.cpu.Decode(instructionBits)
	status = c.cpu.Execute(instruction)

	// Write trace
	if c.traceFile != nil {
		c.traceFile.WriteString(c.Trace(pc) + "\n")
	}

	// Increment step counter
	c.step_counter++
	if !status {
		return false
	}

	return true
}

// Builds a three-line status output to debug simulator.
//
// Returns:
//  string contining trace output
func (c *Computer) Trace(program_counter uint32) (output string) {
	// Build Flags int
	cpsr, _ := c.cpu.FetchRegister(CPSR)
	flags := ExtractShiftBits(cpsr, V, 32)

	output = fmt.Sprintf("%06d %08X %08X %04b\t", c.step_counter, program_counter-4,
		c.ram.Checksum(), flags)
	for i := 0; i < 15; i++ {
		reg, _ := c.cpu.FetchRegister(uint32(i * 4))
		output += fmt.Sprintf("%2d=%08X", i, reg)
		if i == 3 || i == 9 {
			output += "\n\t"
		} else if i != 14 {
			output += "\t"
		}
	}
	c.log.Print(output)

	return
}

// Loads an ELF structed executable file into memory.
//
// Parameters:
//  filePath - a path to the ELF file to open
//  memory - a pointer to a suitable Memory
//
// Returns:
//  err - any error that might have occured
func (c *Computer) LoadELF(filePath string) (err error) {
	// TODO: Perhaps this method could be shortened? A lot of the code is
	// whitespace and logging code. Also, most of the code is not repeated
	// anywhere.

	// Get a clean system
	c.Reset()

	// Setup Logging
	defer c.log.SetPrefix(c.log.Prefix())
	c.log.SetPrefix("Loader: ")

	// Attempt to open file
	c.log.Println("Opening file", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		c.log.Printf("Error reading file (perhaps it doesn't exist)...")
		return
	}
	defer file.Close()

	// Test magic bytes
	c.log.Println("Testing magic bytes...")
	if err = verifyMagic(file); err != nil {
		c.log.Println(err)
		return
	}

	// Read ELF Header
	c.log.Println("Reading ELF header...")
	file.Seek(0, 0)
	elfHeader := new(elf.Header32)
	err = binary.Read(file, binary.LittleEndian, elfHeader)
	if err != nil {
		c.log.Println("Error reading ELF header...")
		return
	}

	c.log.Printf("Program header offset: %d", elfHeader.Phoff)

	// Set PC
	c.log.Printf("Entry Point: %d", elfHeader.Entry)
	c.registers.WriteWord(PC, elfHeader.Entry)

	c.log.Printf("# of program header entires: %d", elfHeader.Phnum)

	// Seek to Program Header start
	file.Seek(int64(elfHeader.Phoff), 0)

	// Read Program Headers
	c.log.Println("Reading program headers...")
	pHeader := new(elf.Prog32)
	for i := 0; uint16(i) < elfHeader.Phnum; i++ {
		// Seek to program header
		offset := int64(elfHeader.Phoff) + int64(i)*int64(elfHeader.Phentsize)
		file.Seek(offset, 0)

		// Read program header
		err = binary.Read(file, binary.LittleEndian, pHeader)
		if err != nil {
			c.log.Printf("Error reading program header %d...", i)
			return
		}

		c.log.Printf("Reading program header %d of %d - Offset: %d, Size: %d, Address: %d", i+1, elfHeader.Phnum, pHeader.Off, pHeader.Filesz, pHeader.Vaddr)
		// Seek to offset
		file.Seek(int64(pHeader.Off), 0)

		// Read to RAM
		b := make([]byte, 1)
		var i uint32 = 0
		for ; i < pHeader.Filesz; i++ {
			file.Read(b)
			err = c.ram.WriteByte(pHeader.Vaddr+i, b[0])
			if err != nil {
				err = errors.New("Insuffcient memory.")
				return
			}
		}
	}

	return
}

// Returns the checksum for the RAM
//
// Parameters: None
//
// Returns: checksum as int32
func (c *Computer) Checksum() (checksum int32) {
	checksum = c.ram.Checksum()
	return
}

// Enables tracing
//
// Parameters: None
//
// Returns:
//  err - any error that might have occured
func (c *Computer) EnableTracing() (err error) {
	c.traceFile, err = os.Create("trace.log")
	if err != nil {
		err = errors.New("Unable to open trace file.")
		c.log.Print(err)
	}

	return
}

// Disables tracing
//
// Parameters: None
//
// Returns:
//  err - any error that might have occured
func (c *Computer) DisableTracing() (err error) {
	if c.traceFile != nil {
		err = c.traceFile.Close()
		c.traceFile = nil
	}
	return
}

// Resets memory and registers to a clean state (all values zeroed out).
func (c *Computer) Reset() {
	for i := 0; uint32(i) < c.memSize; i += 4 {
		c.ram.WriteWord(uint32(i), 0x0)
	}

	for i := 0; uint32(i) < (SPSR_irq + 4); i += 4 {
		c.registers.WriteWord(uint32(i), 0x0)
	}

	if c.traceFile != nil {
		c.EnableTracing()
	}

	// Set SPs
	c.cpu.WriteRegister(SP, 0x7000)
	c.cpu.WriteRegister(SP_irq, 0x7FF0)
	c.cpu.WriteRegister(SP_svc, 0x78F0)

	// Set mode
	c.cpu.WriteRegister(CPSR, 0x1F)

	c.step_counter = 1
}

// Helper Methods

// Verifies if a given 4 bytes are the correct signature for an ELF header.
func verifyMagic(file *os.File) (err error) {
	magic := [4]byte{}

	err = binary.Read(file, binary.LittleEndian, &magic)
	if err != nil || magic[0] != 0x7f || magic[1] != 'E' || magic[2] != 'L' || magic[3] != 'F' {
		err = errors.New("ELF magic bytes were incorrect.")
	}

	return
}
