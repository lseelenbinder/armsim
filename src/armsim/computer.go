package armsim

type Computer struct {
	ram       *Memory
	registers *Memory
	cpu       *CPU
}

func NewComputer(memSize uint32) (c *Computer) {
	c = new(Computer)

	// Initialze RAM of memSize
	c.ram = NewMemory(memSize)

	// Initialze a register bank to contain all 16 registers
	c.registers = NewMemory(r15 + 4)

	// Initialze CPU with RAM and registers
	c.cpu = NewCPU(c.ram, c.registers)

	return
}

func (c *Computer) Run() {
	for {
		if !c.Step() {
			break
		}
	}
}

func (c *Computer) Step() bool {
	instruction := c.cpu.Fetch()

	// Don't continue if the instruction is useless
	if instruction == 0x0 {
		return false
	}

	// Not easily testable at the moment
	c.cpu.Decode()
	c.cpu.Execute()

	return true
}
