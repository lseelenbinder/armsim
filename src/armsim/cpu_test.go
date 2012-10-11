// Filename: cpu_test.go
// Contents: Tests for CPU struct.

package armsim

import "testing"

func TestNewCPU(t *testing.T) {
	var cpu *CPU
	ram := NewMemory(32*1024, nil)
	registers := NewMemory(16*4, nil)

	// initialize CPU
	cpu = NewCPU(ram, registers, nil)

	if cpu.ram != ram {
		t.Fatal("RAM not referenced in CPU.")
	}
	if cpu.registers != registers {
		t.Fatal("Registers not referenced in CPU.")
	}
}

func TestFetch(t *testing.T) {
	// Setup
	ram := NewMemory(32*1024, nil)
	registers := NewMemory(16*4, nil)
	cpu := NewCPU(ram, registers, nil)
	var word uint32
	var test uint32
	var pc uint32

	// Test Normal Fetch
	test = 0xFF00FF14
	ram.WriteWord(0x0, test)
	// Simulate setting PC to 0
	registers.WriteWord(PC, 0)
	word = cpu.Fetch()
	if word != test {
		t.Fatalf("Incorrect word fetched. Expected %#x got %#x", test, word)
	}
	pc, _ = registers.ReadWord(PC)
	if pc != 4 {
		t.Fatal("PC was not incremented.")
	}

	// Test Normal Fetch
	test = 0xFF00FF15
	// PC is now at 0x4
	ram.WriteWord(0x4, test)
	word = cpu.Fetch()
	if word != test {
		t.Fatalf("Incorrect word fetched. Expected %#x got %#x", test, word)
	}
	pc, _ = registers.ReadWord(PC)
	if pc != 8 {
		t.Fatal("PC was not incremented.")
	}
}

func TestDecode(t *testing.T) {
	// Setup
	ram := NewMemory(32*1024, nil)
	registers := NewMemory(16*4, nil)
	cpu := NewCPU(ram, registers, nil)

	// I'm depending on my Instruction unit tests
	cpu.Decode(0x0)
}

func TestExecute(t *testing.T) {
	// Setup
	ram := NewMemory(32*1024, nil)
	registers := NewMemory(16*4, nil)
	cpu := NewCPU(ram, registers, nil)
	instruction := Decode(0x0, nil)

	cpu.Execute(instruction)
}
