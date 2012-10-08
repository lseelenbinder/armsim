package armsim

import "testing"

func TestMOV(t *testing.T) {
	// Test MOV immediate
	c := NewComputer(32 * 1024, nil)
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
}
