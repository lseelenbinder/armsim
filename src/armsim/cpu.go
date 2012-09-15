package armsim

import (
	"log"
	"time"
)

// Initialize register position constants
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
	SP = r13
	LR = r14
	PC = r15
)

type CPU struct {
	ram       *Memory
	registers *Memory
}

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

func (cpu *CPU) Fetch() (instruction uint32) {
	address, err := cpu.registers.ReadWord(PC)
	if err != nil {
		log.Panic("Unable to read PC.")
	}
	log.Printf("Current PC: %#x", address)

	instruction, err = cpu.ram.ReadWord(address)
	if err != nil {
		log.Panic("Unable to read next instruction.")
	}
	cpu.registers.WriteWord(PC, address+4)
	return
}

func (cpu *CPU) Decode() {
	// Does nothing; this is a stub.
}

func (cpu *CPU) Execute() {
	time.Sleep(time.Duration(250) * time.Millisecond)
}
