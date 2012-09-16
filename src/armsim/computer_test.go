package armsim

import "testing"

func TestComputer(t *testing.T) {
	// Ensure Computer type exists with ram, registers, and etc.
	c := new(Computer)
	_ = c.ram
	_ = c.registers
	_ = c.cpu
}

func TestNewComputer(t *testing.T) {
	computer := NewComputer(32 * 1024)

	if computer == nil {
		t.Fatal("Did not initialize Computer.")
	}

	// Test RAM
	if computer.ram == nil {
		t.Fatal("Did not initialize RAM.")
	}
	_, err := computer.ram.ReadByte(32*1024 - 1)
	if err != nil {
		t.Fatal("Did not initialize RAM to correct size.")
	}

	// Test Registers
	if computer.registers == nil {
		t.Fatal("Did not initialize Registers.")
	}
	_, err = computer.registers.ReadByte(63)
	if err != nil {
		t.Fatal("Did not initialize registers to correct size (too small).")
	}
	_, err = computer.registers.ReadByte(65)
	if err == nil {
		t.Fatal("Did not initialize registers to correct size (too big).")
	}

	// Test CPU
	if computer.cpu == nil {
		t.Fatal("Did not initialize CPU.")
	}
	if computer.cpu.ram != computer.ram {
		t.Fatal("Did not initialize CPU with correct RAM.")
	}
	if computer.cpu.registers != computer.registers {
		t.Fatal("Did not initialize CPU with correct registers.")
	}
}

func TestRun(t *testing.T) {
	// Setup
	c := NewComputer(32 * 1024)

	// Simulate loading program and PC
	c.ram.WriteWord(0x0, 0x67)
	c.ram.WriteWord(0x4, 0x67)
	c.ram.WriteWord(0x8, 0x67)
	c.registers.WriteWord(PC, 0x0)

	c.Run()
	pc, _ := c.registers.ReadWord(PC)

	// Should be the last position + 8
	// (a fetch for the 0x0 value and an increment)
	if pc != 0x08+0x8 {
		t.Fatal("Did not run until 0x0 was fetched.")
	}

	// Simulate loading program and PC
	c.ram.WriteWord(0x80, 0x67)
	c.ram.WriteWord(0x84, 0x67)
	c.ram.WriteWord(0x88, 0x67)
	c.registers.WriteWord(PC, 0x80)

	c.Run()
	pc, _ = c.registers.ReadWord(PC)

	// Should be the last position + 8
	// (a fetch for the 0x0 value and an increment)
	if pc != 0x88+0x8 {
		t.Fatal("Did not run until 0x0 was fetched.")
	}
}

func TestStep(t *testing.T) {
	// Setup
	c := NewComputer(32 * 1024)

	// Simulate loading program and PC
	c.ram.WriteWord(0x0, 0x67)
	c.ram.WriteWord(0x4, 0x67)
	c.ram.WriteWord(0x8, 0x67)
	c.registers.WriteWord(PC, 0x0)

	c.Step()
	pc, _ := c.registers.ReadWord(PC)
	if pc != 0x04 {
		t.Fatal("Did not step to 0x4.")
	}

	c.Step()
	pc, _ = c.registers.ReadWord(PC)
	if pc != 0x08 {
		t.Fatal("Did not step to 0x8.")
	}

	c.Step()
	pc, _ = c.registers.ReadWord(PC)
	if pc != 0x0C {
		t.Fatal("Did not step to 0xC.")
	}
}
