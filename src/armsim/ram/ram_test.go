package ram

import "testing"

func TestNew(t *testing.T) {
  // Test the RAM initializer at 32k
  runNew(t, 32*1024)
  // Test the RAM initializer at 32m
  runNew(t, 32*1024*1024)
}

// Helper method for the RAM initializer testing code
func runNew(t *testing.T, nBytes uint32) {
  // Setup
  ram := New(nBytes)

  // Ensure RAM is the same length as nBytes
  if uint32(len(ram.memory)) != nBytes {
    t.Fatal("Did not allocate RAM to specified nBytes")
  }

  // Ensure RAM is initialized to zero
  for i := uint32(0); i < nBytes; i++ {
    if ram.memory[i] != byte(0) {
      t.Fatal("Did not initalize RAM sector", i)
    }
  }
}

func TestWriteByte(t *testing.T) {
  // Setup
  ram := New(32)
  var byteToWrite byte = 95

  // Test under normal conditions
  err := ram.WriteByte(0, byteToWrite)
  if err != nil && byteToWrite != ram.memory[0] {
    t.Fatal("Did not properly write byte.")
  }

  // Test address out of range
  err = ram.WriteByte(33, byteToWrite)
  if err == nil {
    t.Fatal("Attempted to write out of range address.")
  }
}

func TestReadByte(t *testing.T) {
  // Setup
  ram := New(32)
  var dataToWrite byte = 0xFF

  // Test normal reading of bytes
  ram.WriteByte(0, dataToWrite)
  b, err := ram.ReadByte(0)
  if err != nil || b != 0xFF {
    t.Fatalf("expected %#x got %#x", dataToWrite, b)
  }

  // Test address out of range
  b, err = ram.ReadByte(80)
  if err == nil {
    t.Fatal("Attempted to read out of range address.")
  }
}

func TestWriteHalfWord(t *testing.T) {
  // Setup
  ram := New(32)

  // Test normal writing of halfwords
  // Lots of tests because this is faily complicated with shifts and casts.
  var hwToWrite uint16

  hwToWrite = 0xFFFF
  err := ram.WriteHalfWord(0, hwToWrite)
  b1, _ := ram.ReadByte(0)
  b2, _ := ram.ReadByte(1)

  if err != nil || b1 != 0xFF || b2 != 0xFF {
    t.Fatal("Did not properly write half word. Expected:", hwToWrite, "Got:", b1, b2)
  }

  hwToWrite = 0xFF00
  err = ram.WriteHalfWord(0, hwToWrite)
  b1, _ = ram.ReadByte(0)
  b2, _ = ram.ReadByte(1)

  if err != nil || b1 != 0xFF || b2 != 0x00 {
    t.Fatal("Did not properly write half word. Expected:", hwToWrite, "Got:", b1, b2)
  }

  hwToWrite = 0x0F0D
  err = ram.WriteHalfWord(0, hwToWrite)
  b1, _ = ram.ReadByte(0)
  b2, _ = ram.ReadByte(1)

  if err != nil || b1 != 0x0F || b2 != 0x0D {
    t.Fatal("Did not properly write half word. Expected:", hwToWrite, "Got:", b1, b2)
  }

  // Test writing to odd address
  err = ram.WriteHalfWord(1, hwToWrite)

  if err == nil {
    t.Fatal("Wrote halfword to odd address. This should not happen.")
  }

  // Test writing to out-of-bounds address
  err = ram.WriteHalfWord(34, hwToWrite)
  if err == nil {
    t.Fatal("Wrote halfword to an out-of-bounds address.")
  }
}

func TestReadHalfWord(t *testing.T) {
  // Setup
  ram := New(32)
  var hwToWrite uint16 = 0xFFFF

  // Test normal reading of bytes
  ram.WriteHalfWord(0, hwToWrite)
  b, err := ram.ReadHalfWord(0)
  if err != nil || b != hwToWrite {
    t.Fatalf("expected %#x got %#x", hwToWrite, b)
  }

  // Test address out of range
  _, err = ram.ReadHalfWord(80)
  if err == nil {
    t.Fatal("Attempted to read out of range address.")
  }
}

func TestWriteWord(t *testing.T) {
  // Setup
  ram := New(32)

  // Test normal writing of words
  // Lots of tests because this is faily complicated with shifts and casts.

  var wordToWrite uint32

  wordToWrite = 0xFFFFFFFF
  err := ram.WriteWord(0, wordToWrite)
  hw1, _ := ram.ReadHalfWord(0)
  hw2, _ := ram.ReadHalfWord(2)

  if err != nil || hw1 != 0xFFFF || hw2 != 0xFFFF {
    t.Fatalf("expected: %#x to be written instead got %#x & %#x", wordToWrite, hw1, hw2)
  }

  wordToWrite = 0xF0F0F0F0
  err = ram.WriteWord(0, wordToWrite)
  hw1, _ = ram.ReadHalfWord(0)
  hw2, _ = ram.ReadHalfWord(2)

  if err != nil || hw1 != 0xF0F0 || hw2 != 0xF0F0 {
    t.Fatalf("expected: %#x to be written instead got %#x & %#x", wordToWrite, hw1, hw2)
  }

  wordToWrite = 0xF0F0F0F0
  err = ram.WriteWord(4, wordToWrite)
  hw1, _ = ram.ReadHalfWord(4)
  hw2, _ = ram.ReadHalfWord(6)

  if err != nil || hw1 != 0xF0F0 || hw2 != 0xF0F0 {
    t.Fatalf("expected: %#x to be written instead got %#x & %#x", wordToWrite, hw1, hw2)
  }

  wordToWrite = 0x2
  err = ram.WriteWord(4, wordToWrite)
  hw1, _ = ram.ReadHalfWord(4)
  hw2, _ = ram.ReadHalfWord(6)

  if err != nil || hw1 != 0x0000 || hw2 != 0x0002 {
    t.Fatalf("expected: %#x to be written instead got %#x & %#x", wordToWrite, hw1, hw2)
  }

  // Test writing to odd address
  err = ram.WriteWord(1, wordToWrite)
  if err == nil {
    t.Fatal("Wrote word to an address indivisible by 4. This should not happen.")
  }

  err = ram.WriteWord(3, wordToWrite)
  if err == nil {
    t.Fatal("Wrote word to an address indivisible by 4. This should not happen.")
  }

  err = ram.WriteWord(6, wordToWrite)
  if err == nil {
    t.Fatal("Wrote word to an address indivisible by 4. This should not happen.")
  }

  // Test writing to out-of-bounds address
  err = ram.WriteWord(34, wordToWrite)
  if err == nil {
    t.Fatal("Wrote word to an out-of-bounds address.")
  }
}

func TestReadWord(t *testing.T) {
  // Setup
  ram := New(32)
  var wordToWrite uint32 = 0xFFFFFFFF

  // Test normal reading of bytes
  ram.WriteWord(0, wordToWrite)
  b, err := ram.ReadWord(0)
  if err != nil || b != wordToWrite {
    t.Fatalf("expected %#x got %#x", wordToWrite, b)
  }

  wordToWrite = 0xFDF67859
  ram.WriteWord(0, wordToWrite)
  b, err = ram.ReadWord(0)
  if err != nil || b != wordToWrite {
    t.Fatalf("expected %#x got %#x", wordToWrite, b)
  }

  // Test address out of range
  _, err = ram.ReadWord(80)
  if err == nil {
    t.Fatal("Attempted to read out of range address.")
  }
}

func TestChecksum(t *testing.T) {
  // Test empty ram
  ram := New(0)
  if ram.Checksum() != 0 {
    t.Fatal("zero byte ram checksum is not 0")
  }

  // Note: I used WolframAlpha to calculate these small checksums. They seem to
  // be correct.
  ram = New(5)
  if ram.Checksum() != 10 {
    t.Fatal("checksum for 5 bytes of empty ram should be 10")
  }

  ram = New(5)
  ram.WriteByte(0, 0x1)
  check := ram.Checksum()
  if check != 11 {
    t.Fatalf("expected checksum of %d; got: %d", 0x1^0, check)
  }

  ram = New(5)
  ram.WriteByte(4, 11)
  check = ram.Checksum()
  if check != 21 {
    t.Fatalf("expected checksum of %d; got: %d", 21, check)
  }

  ram = New(5)
  ram.WriteByte(3, 0x65)
  check = ram.Checksum()
  if check != 109 {
    t.Fatalf("expected checksum of %d; got: %d", 109, check)
  }

}

func TestTestFlag(t *testing.T) {
  // Setup
  // Newialize a single word for simplicity
  ram := New(4)

  ram.WriteWord(0, 0x0001)
  flag, err := ram.TestFlag(0, 0)
  if err != nil || flag != true {
    t.Fatal("tested bit was false; expected true")
  }

  ram.WriteWord(0, 0x2)
  flag, err = ram.TestFlag(0, 0)
  if err != nil || flag != false {
    t.Fatal("tested bit was true; expected false")
  }
  flag, err = ram.TestFlag(0, 1)
  if err != nil || flag != true {
    t.Fatal("tested bit was false; expected true")
  }

  ram.WriteWord(0, 0xFFFFFFFF)
  for i := 0; i < 32; i++ {
    flag, err = ram.TestFlag(0, uint32(i))
    if err != nil || flag != true {
      t.Fatal("tested bit was false; expected true")
    }
  }

  ram.WriteWord(0, 0x0)
  for i := 0; i < 32; i++ {
    flag, err = ram.TestFlag(0, uint32(i))
    if err != nil || flag != false {
      t.Fatal("tested bit was true; expected false")
    }
  }

  // Test out-of-bounds read
  _, err = ram.TestFlag(5, 0)
  if err == nil {
    t.Fatal("Shouldn't have been able to read this byte.")
  }

}

func TestSetFlag(t *testing.T) {
  // Setup
  ram := New(4)

  ram.WriteWord(0, 0xFFFFFFFF)
  ram.SetFlag(0, 0, true)
  flag, err := ram.TestFlag(0, 0)
  if err != nil || flag != true {
    t.Fatal("tested bit was not set to true")
  }

  err = ram.SetFlag(0, 0, false)
  flag, _ = ram.TestFlag(0, 0)
  if err != nil || flag != false {
    t.Fatal("tested bit was not set to false")
  }

  ram.WriteWord(0, 0xFFFFFFFF)
  err = ram.SetFlag(0, 5, false)
  flag, _ = ram.TestFlag(0, 5)
  if err != nil || flag != false {
    t.Fatal("tested bit was not set to false")
  }

  // Out-of-bounds testing
  err = ram.SetFlag(5, 5, false)
  if err == nil {
    t.Fatal("should have panicked")
  }
}

func TestExtractBits(t *testing.T) {
  err := ExtractBits(0xb5, 1, 3)
  if err != 0x04 {
    t.Fatalf("expected: 0x04 got: %#x", err)
  }

  err = ExtractBits(0xb5, 0, 3)
  if err != 0x05 {
    t.Fatalf("expected: 0x05 got: %#x", err)
  }

  // Implicity works
  err = ExtractBits(0xb5, 0, 33)
  if err != 0xb5 {
    t.Fatalf("expected: 0x05 got: %#x", err)
  }

  // Explicitly fails due to typing
  // ExtractBits(0xb5, -1, 33)
}
