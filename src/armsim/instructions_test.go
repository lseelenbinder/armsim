package armsim

import "testing"

func TestCLevel(t *testing.T) {
	c := NewComputer(32*1024, nil)
	c.LoadELF("../../test/ctest.exe")
	c.Run(nil, nil)
}

func TestBLevel(t *testing.T) {
	c := NewComputer(32*1024, nil)
	c.LoadELF("../../test/btest.exe")
	c.Run(nil, nil)
}

func TestMOV(t *testing.T) {
	// Test MOV immediate
	c := NewComputer(32, nil)
	c.registers.WriteWord(PC, 0x4)
	c.ram.WriteWord(0x4, 0xe3a02030)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 48 {
		t.Fatal("expected 48, got", word)
	}

	c.ram.WriteWord(0x8, 0xe3a02036)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 54 {
		t.Fatal("expected 54, got", word)
	}

	c.ram.WriteWord(0xc, 0xe3a03036)
	c.Step()
	if word, _ := c.registers.ReadWord(r3); word != 54 {
		t.Fatal("expected 54, got", word)
	}

	// Test MOV immediate with a rotate
	c.Reset()
	c.registers.WriteWord(PC, 0x4)
	c.ram.WriteWord(0x4, 0xe3a02130)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0xc {
		t.Fatalf("expected 0xC, got 0x%x", word)
	}

	// Test MOV immediate shift
	c.Reset()
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r1, 0x50)
	c.ram.WriteWord(0x4, 0xe1a020a1)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x28 {
		t.Fatalf("expected 0x28, got 0x%x", word)
	}

	// Test MOV immediate shift
	c.Reset()
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r1, 0x40)
	c.ram.WriteWord(0x4, 0xE1A021A1)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x8 {
		t.Fatalf("expected 0x28, got 0x%x", word)
	}

	// Test MOV register shift
	c.Reset()
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r1, 0x50)
	c.registers.WriteWord(r3, 0x1)
	c.ram.WriteWord(0x4, 0xe1a02331)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x28 {
		t.Fatalf("expected 0x28, got 0x%x", word)
	}

	// Test MOV register shift
	c.Reset()
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r1, 0x40)
	c.registers.WriteWord(r3, 0x3)
	c.ram.WriteWord(0x4, 0xe1a02331)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x8 {
		t.Fatalf("expected 0x8, got 0x%x", word)
	}
}

func TestMNV(t *testing.T) {
	// Test MNV, rely on MOV tests mostly
	c := NewComputer(32, nil)
	c.registers.WriteWord(PC, 0x4)
	c.ram.WriteWord(0x4, 0xE3E02030)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0xFFFFFFCF {
		t.Fatal("expected 0x0xFFFFFFCF, got", word)
	}

	c.ram.WriteWord(0x8, 0xE3E02036)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0xFFFFFFC9 {
		t.Fatal("expected 54, got", word)
	}
}

func TestADD(t *testing.T) {
	c := NewComputer(32, nil)
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r4, 0x10)
	c.ram.WriteWord(0x4, 0xE2842030)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x40 {
		t.Fatal("expected 0x40, got", word)
	}

	c.Reset()
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r4, 0x40)
	c.ram.WriteWord(0x4, 0xE2842030)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x70 {
		t.Fatal("expected 0x40, got", word)
	}
}

func TestSUB(t *testing.T) {
	c := NewComputer(32, nil)
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r4, 0x40)
	c.ram.WriteWord(0x4, 0xE2442030)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x10 {
		t.Fatal("expected 0x10, got", word)
	}

	c.Reset()
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r4, 0x60)
	c.ram.WriteWord(0x4, 0xE2442030)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x30 {
		t.Fatal("expected 0x30, got", word)
	}
}

func TestRSB(t *testing.T) {
	c := NewComputer(32, nil)
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r4, 0x20)
	c.ram.WriteWord(0x4, 0xE2642030)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x10 {
		t.Fatal("expected 0x10, got", word)
	}

	c.Reset()
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r4, 0x15)
	c.ram.WriteWord(0x4, 0xE2642030)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x1B {
		t.Fatal("expected 0x30, got", word)
	}
}

func TestAND(t *testing.T) {
	c := NewComputer(32, nil)
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r4, 0x10)
	c.ram.WriteWord(0x4, 0xE2042030)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x10 {
		t.Fatal("expected 0x10, got", word)
	}

	c.Reset()
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r4, 0x20)
	c.ram.WriteWord(0x4, 0xE2042030)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x20 {
		t.Fatal("expected 0x20, got", word)
	}
}

func TestBIC(t *testing.T) {
	c := NewComputer(32, nil)
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r4, 0x10)
	c.ram.WriteWord(0x4, 0xE2E22030)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x0 {
		t.Fatal("expected 0x0, got", word)
	}

	c.Reset()
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r4, 0x20)
	c.ram.WriteWord(0x4, 0xE2E22030)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x0 {
		t.Fatal("expected 0x0, got", word)
	}
}

func TestEOR(t *testing.T) {
	c := NewComputer(32, nil)
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r4, 0x10)
	c.ram.WriteWord(0x4, 0xE2242030)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x20 {
		t.Fatal("expected 0x20, got", word)
	}

	c.Reset()
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r4, 0x15)
	c.ram.WriteWord(0x4, 0xE2242030)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x25 {
		t.Fatal("expected 0x25, got", word)
	}
}

func TestORR(t *testing.T) {
	c := NewComputer(32, nil)
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r4, 0x10)
	c.ram.WriteWord(0x4, 0xE3842030)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x30 {
		t.Fatal("expected 0x30, got", word)
	}

	c.Reset()
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r4, 0x15)
	c.ram.WriteWord(0x4, 0xE3842030)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x35 {
		t.Fatal("expected 0x35, got", word)
	}
}

func TestMUL(t *testing.T) {
	c := NewComputer(32, nil)
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r4, 0x10) // Rs
	c.registers.WriteWord(r3, 0x30) // Rm
	c.ram.WriteWord(0x4, 0xE0020394)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x300 {
		t.Fatal("expected 0x300, got", word)
	}

	c.Reset()
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r4, 0x15) // Rs
	c.registers.WriteWord(r3, 0x30) // Rm
	c.ram.WriteWord(0x4, 0xE0020394)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x3f0 {
		t.Fatal("expected 0x3f0, got", word)
	}
}
