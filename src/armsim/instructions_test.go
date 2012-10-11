package armsim

import "testing"

func TestCLevel(t *testing.T) {
	c := NewComputer(32*1024, nil)
	c.LoadELF("../../test/ctest.exe")
	c.Run(nil, nil)
}

func TestMOV(t *testing.T) {
	// Test MOV immediate
	c := NewComputer(32, nil)
	c.registers.WriteWord(PC, 0x4)
	c.ram.WriteWord(0x4, 0xe3a02030)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 48 {
		t.Fatal("Exepected 48, got", word)
	}

	c.ram.WriteWord(0x8, 0xe3a02036)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 54 {
		t.Fatal("Exepected 54, got", word)
	}

	c.ram.WriteWord(0xc, 0xe3a03036)
	c.Step()
	if word, _ := c.registers.ReadWord(r3); word != 54 {
		t.Fatal("Exepected 54, got", word)
	}

	// Test MOV immediate with a rotate
	c.Reset()
	c.registers.WriteWord(PC, 0x4)
	c.ram.WriteWord(0x4, 0xe3a02130)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0xc {
		t.Fatalf("Exepected 0x12, got 0x%x", word)
	}

	// Test MOV immediate shift
	c.Reset()
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r1, 0x50)
	c.ram.WriteWord(0x4, 0xe1a020a1)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x28 {
		t.Fatalf("Exepected 0x28, got 0x%x", word)
	}

	// Test MOV immediate shift
	c.Reset()
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r1, 0x40)
	c.ram.WriteWord(0x4, 0xE1A021A1)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x8 {
		t.Fatalf("Exepected 0x28, got 0x%x", word)
	}

	// Test MOV register shift
	c.Reset()
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r1, 0x50)
	c.registers.WriteWord(r3, 0x1)
	c.ram.WriteWord(0x4, 0xe1a02331)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x28 {
		t.Fatalf("Exepected 0x28, got 0x%x", word)
	}

	// Test MOV register shift
	c.Reset()
	c.registers.WriteWord(PC, 0x4)
	c.registers.WriteWord(r1, 0x40)
	c.registers.WriteWord(r3, 0x3)
	c.ram.WriteWord(0x4, 0xe1a02331)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0x8 {
		t.Fatalf("Exepected 0x8, got 0x%x", word)
	}
}
func TestMNV(t *testing.T) {
	// Test MNV, rely on MOV tests mostly
	c := NewComputer(32, nil)
	c.registers.WriteWord(PC, 0x4)
	c.ram.WriteWord(0x4, 0xE3E02030)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0xFFFFFFCF {
		t.Fatal("Exepected 0x0xFFFFFFCF, got", word)
	}

	c.ram.WriteWord(0x8, 0xE3E02036)
	c.Step()
	if word, _ := c.registers.ReadWord(r2); word != 0xFFFFFFC9 {
		t.Fatal("Exepected 54, got", word)
	}
}