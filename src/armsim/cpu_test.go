package armsim

import (
	"testing"
	"time"
)

func TestNewCPU(t *testing.T) {
	var cpu *CPU
	ram := NewMemory(32 * 1024)
	registers := NewMemory(16 * 4)

	// initialize CPU
	cpu = NewCPU(ram, registers)

	if cpu.ram != ram {
		t.Fatal("RAM not referenced in CPU.")
	}
	if cpu.registers != registers {
		t.Fatal("Registers not referenced in CPU.")
	}
}

func TestFetch(t *testing.T) {
	// Setup
	ram := NewMemory(32 * 1024)
	registers := NewMemory(16 * 4)
	cpu := NewCPU(ram, registers)
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
	ram := NewMemory(32 * 1024)
	registers := NewMemory(16 * 4)
	cpu := NewCPU(ram, registers)

	// Does nothing
	cpu.Decode()
}

func TestExecute(t *testing.T) {
	// Setup
	ram := NewMemory(32 * 1024)
	registers := NewMemory(16 * 4)
	cpu := NewCPU(ram, registers)

	// Should wait at least 0.25 seconds and no more than 0.26 seconds
	start := time.Now().Nanosecond()
	cpu.Execute()
	end := time.Now().Nanosecond()

	// Accuracy of 100 milliseconds (0.1 seconds)
	if int64(end-start)-(time.Duration(250)*time.Millisecond).Nanoseconds() >
		(time.Duration(100) * time.Millisecond).Nanoseconds() {
		t.Fatal("Did not wait a quarter second.")
	}
}
