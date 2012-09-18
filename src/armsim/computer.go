// Filename: computer.go
// Contents: The Computer struct and related methods.

package armsim

import (
	"debug/elf"
	"encoding/binary"
	"errors"
	"fmt"
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
	registers *Memory

	// A reference to the CPU for the simulator
	cpu *CPU

	// A simple counter to track number of execution cycles
	step_counter uint64

	// Logger class
	log *log.Logger
}

// Initializes a Computer
//
// Parameters:
//  memSize - a unsigned 32-bit integer specifying the size of the RAM in the
//  computer
//
// Returns:
//  a pointer to the newly created Computer
func NewComputer(memSize uint32) (c *Computer) {
	c = new(Computer)

	// Setup logging
	c.log = log.New(os.Stdout, "Computer: ", 0)

	// Initialize RAM of memSize
	c.ram = NewMemory(memSize)

	// Initialize a register bank to contain all 16 registers + CPSR
	c.registers = NewMemory(CPSR + 4)

	// Initialize CPU with RAM and registers
	c.cpu = NewCPU(c.ram, c.registers)

	return
}

// Simulates the running of the a computer. It executes the fetch, execute,
// decode cycle until fetch returns false (signifying an instruction of 0x0).
func (c *Computer) Run() {
	for {
		if !c.Step() {
			break
		}
	}
}

// Performs a single execution cycle. Take no parameters and returns a boolean
// signifying if the cycle was completed (a cycle will not complete if the
// instrution fetched is 0x0).
func (c *Computer) Step() bool {
	instruction := c.cpu.Fetch()

	// Don't continue if the instruction is useless
	if instruction == 0x0 {
		return false
	}

	// Not easily testable at the moment
	c.cpu.Decode()
	c.cpu.Execute()

	// Increment step counter
	c.step_counter++

	return true
}

// Builds a three-line status output to debug simulator.
//
// Returns:
//  string contining trace output
func (c *Computer) Trace() (output string) {
	program_counter, _ := c.registers.ReadWord(PC)

	// Build Flags int
	cpsr, _ := c.registers.ReadWord(CPSR)
	flags := ExtractBits(cpsr, N, F) >> F

	output = fmt.Sprintf("%06d %08x %08x %04d\t", c.step_counter, program_counter,
		c.ram.Checksum(), flags)
	for i := 0; i < 16; i++ {
		reg, _ := c.registers.ReadWord(uint32(i))
		output += fmt.Sprintf("%2d=%08x", i, reg)
		if i == 3 || i == 9 {
			output += "\n\t"
		} else if i != 15 {
			output += "\t"
		}
	}
	c.log.Println("*****TRACE*****")
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
