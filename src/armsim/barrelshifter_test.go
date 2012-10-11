package armsim

import "testing"

func TestROR(t *testing.T) {
	// Test without a rotate value
	b := BarrelShifter{ROR, 0, 0xF}

	if result := b.Shift(); result != 0xF {
		t.Fatalf("Expected 0xF; got 0x%x", result)
	}

	b = BarrelShifter{ROR, 2, 0xF}

	if result := b.Shift(); result != 0xc0000003 {
		t.Fatalf("Expected 0xc3; got 0x%x", result)
	}
}

func TestLSL(t *testing.T) {
	// Test without a rotate value
	b := BarrelShifter{LSL, 0, 0xF}

	if result := b.Shift(); result != 0xF {
		t.Fatalf("Expected 0xF; got 0x%x", result)
	}

	b = BarrelShifter{LSL, 2, 0xF}

	if result := b.Shift(); result != 0x3C {
		t.Fatalf("Expected 0x3C; got 0x%x", result)
	}

	b = BarrelShifter{LSL, 4, 0xF}

	if result := b.Shift(); result != 0xF0 {
		t.Fatalf("Expected 0xF0; got 0x%x", result)
	}

	b = BarrelShifter{LSL, 32, 0xF}

	if result := b.Shift(); result != 0x0 {
		t.Fatalf("Expected 0x0; got 0x%x", result)
	}
}
