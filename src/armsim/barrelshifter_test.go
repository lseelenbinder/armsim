package armsim

import (
	"os"
	"testing"
)

func TestROR(t *testing.T) {
	// Test without a rotate value
	b := BarrelShifter{ROR, 0, 0xF, 0, 0, false, nil}

	if result := b.Shift(); result != 0xF {
		t.Fatalf("Expected 0xF; got 0x%x", result)
	}

	b = BarrelShifter{ROR, 2, 0xF, 0, 0, false, nil}

	if result := b.Shift(); result != 0xc0000003 {
		t.Fatalf("Expected 0xc3; got 0x%x", result)
	}
}

func TestLSL(t *testing.T) {
	// Test without a rotate value
	b := BarrelShifter{LSL, 0, 0xF, 0, 0, false, nil}

	if result := b.Shift(); result != 0xF {
		t.Fatalf("Expected 0xF; got 0x%x", result)
	}

	b = BarrelShifter{LSL, 2, 0xF, 0, 0, false, nil}

	if result := b.Shift(); result != 0x3C {
		t.Fatalf("Expected 0x3C; got 0x%x", result)
	}

	b = BarrelShifter{LSL, 4, 0xF, 0, 0, false, nil}

	if result := b.Shift(); result != 0xF0 {
		t.Fatalf("Expected 0xF0; got 0x%x", result)
	}

	b = BarrelShifter{LSL, 32, 0xF, 0, 0, false, nil}

	if result := b.Shift(); result != 0x0 {
		t.Fatalf("Expected 0x0; got 0x%x", result)
	}
}

func TestOperand2(t *testing.T) {

}

func TestBDisassemble(t *testing.T) {
	c := NewComputer(32*1024, os.Stderr)

	// Test immediate
	b := NewFromOperand2(0xfb5, true, c.cpu)
	if a := b.Disassemble(); a != "#724" {
		t.Fatal("Expected '#724', got", a)
	}

	b = NewFromOperand2(0x4a1, true, c.cpu)
	if a := b.Disassemble(); a != "#2701131776" {
		t.Fatal("Expected '#2701131776', got", a)
	}

	// Test immediate shift
	b = NewFromOperand2(0x141, false, c.cpu)
	if a := b.Disassemble(); a != "r1, asr #2" {
		t.Fatal("Expected 'asr #2', got", a)
	}

	b = NewFromOperand2(0x121, false, c.cpu)
	if a := b.Disassemble(); a != "r1, lsr #2" {
		t.Fatal("Expected 'lsr #2', got", a)
	}

	// Test register shift TODO
}
