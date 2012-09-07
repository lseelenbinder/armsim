package ram

import "testing"

func TestInit(t *testing.T) {
  // Test the RAM initializer at 32k
  runInit(t, 32*1024)
  // Test the RAM initializer at 32m
  runInit(t, 32*1024*1024)
}

// Helper method for the RAM initializer testing code
func runInit(t *testing.T, nBytes uint32) {
  // Setup
  ram := Init(nBytes)

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
  ram := Init(32)
  var byteToWrite byte = 95

  // Test under normal conditions
  result := ram.WriteByte(0, byteToWrite)
  if result == true && byteToWrite != ram.memory[0] {
    t.Fatal("Did not properly write byte.")
  }

  // Test address out of range
  result = ram.WriteByte(33, byteToWrite)
  if result == true {
    t.Fatal("Attempted to write out of range address.")
  }
}

func TestReadByte(t *testing.T) {
  // Setup
  ram := Init(32)
  var dataToWrite byte = 0xFF

  // Test normal reading of bytes
  ram.WriteByte(0, dataToWrite)
  b, result := ram.ReadByte(0)
  if !result || b != 0xFF {
    t.Fatalf("expected %#x got %#x", dataToWrite, b)
  }

  // Test address out of range
  b, result = ram.ReadByte(80)
  if result {
    t.Fatal("Attempted to read out of range address.")
  }
}

func TestWriteHalfWord(t *testing.T) {
  // Setup
  ram := Init(32)

  // Test normal writing of halfwords
  // Lots of tests because this is faily complicated with shifts and casts.
  var hwToWrite uint16

  hwToWrite = 0xFFFF
  result := ram.WriteHalfWord(0, hwToWrite)
  b1, _ := ram.ReadByte(0)
  b2, _ := ram.ReadByte(1)

  if !result || b1 != 0xFF || b2 != 0xFF {
    t.Fatal("Did not properly write half word. Expected:", hwToWrite, "Got:", b1, b2)
  }

  hwToWrite = 0xFF00
  result = ram.WriteHalfWord(0, hwToWrite)
  b1, _ = ram.ReadByte(0)
  b2, _ = ram.ReadByte(1)

  if !result || b1 != 0xFF || b2 != 0x00 {
    t.Fatal("Did not properly write half word. Expected:", hwToWrite, "Got:", b1, b2)
  }

  hwToWrite = 0x0F0D
  result = ram.WriteHalfWord(0, hwToWrite)
  b1, _ = ram.ReadByte(0)
  b2, _ = ram.ReadByte(1)

  if !result || b1 != 0x0F || b2 != 0x0D {
    t.Fatal("Did not properly write half word. Expected:", hwToWrite, "Got:", b1, b2)
  }

  // Test writing to odd address
  result = ram.WriteHalfWord(1, hwToWrite)

  if result {
    t.Fatal("Wrote halfword to odd address. This should not happen.")
  }

  // Test writing to out-of-bounds address
  result = ram.WriteHalfWord(34, hwToWrite)
  if result {
    t.Fatal("Wrote halfword to an out-of-bounds address.")
  }
}

func TestReadHalfWord(t *testing.T) {
  // Setup
  ram := Init(32)
  var hwToWrite uint16 = 0xFFFF

  // Test normal reading of bytes
  ram.WriteHalfWord(0, hwToWrite)
  b, result := ram.ReadHalfWord(0)
  if !result || b != hwToWrite {
    t.Fatalf("expected %#x got %#x", hwToWrite, b)
  }

  // Test address out of range
  _, result = ram.ReadHalfWord(80)
  if result {
    t.Fatal("Attempted to read out of range address.")
  }
}

func TestWriteWord(t *testing.T) {
  // Setup
  ram := Init(32)

  // Test normal writing of words
  // Lots of tests because this is faily complicated with shifts and casts.

  var wordToWrite uint32

  wordToWrite = 0xFFFFFFFF
  result := ram.WriteWord(0, wordToWrite)
  hw1, _ := ram.ReadHalfWord(0)
  hw2, _ := ram.ReadHalfWord(2)

  if !result || hw1 != 0xFFFF || hw2 != 0xFFFF {
    t.Fatalf("expected: %#x to be written instead got %#x & %#x", wordToWrite, hw1, hw2)
  }

  wordToWrite = 0xF0F0F0F0
  result = ram.WriteWord(0, wordToWrite)
  hw1, _ = ram.ReadHalfWord(0)
  hw2, _ = ram.ReadHalfWord(2)

  if !result || hw1 != 0xF0F0 || hw2 != 0xF0F0 {
    t.Fatalf("expected: %#x to be written instead got %#x & %#x", wordToWrite, hw1, hw2)
  }

  wordToWrite = 0xF0F0F0F0
  result = ram.WriteWord(4, wordToWrite)
  hw1, _ = ram.ReadHalfWord(4)
  hw2, _ = ram.ReadHalfWord(6)

  if !result || hw1 != 0xF0F0 || hw2 != 0xF0F0 {
    t.Fatalf("expected: %#x to be written instead got %#x & %#x", wordToWrite, hw1, hw2)
  }

  wordToWrite = 0x2
  result = ram.WriteWord(4, wordToWrite)
  hw1, _ = ram.ReadHalfWord(4)
  hw2, _ = ram.ReadHalfWord(6)

  if !result || hw1 != 0x0000 || hw2 != 0x0002 {
    t.Fatalf("expected: %#x to be written instead got %#x & %#x", wordToWrite, hw1, hw2)
  }

  // Test writing to odd address
  result = ram.WriteWord(1, wordToWrite)
  if result {
    t.Fatal("Wrote word to an address indivisible by 4. This should not happen.")
  }

  result = ram.WriteWord(3, wordToWrite)
  if result {
    t.Fatal("Wrote word to an address indivisible by 4. This should not happen.")
  }

  result = ram.WriteWord(6, wordToWrite)
  if result {
    t.Fatal("Wrote word to an address indivisible by 4. This should not happen.")
  }

  // Test writing to out-of-bounds address
  result = ram.WriteWord(34, wordToWrite)
  if result {
    t.Fatal("Wrote word to an out-of-bounds address.")
  }
}

func TestReadWord(t *testing.T) {
  // Setup
  ram := Init(32)
  var wordToWrite uint32 = 0xFFFFFFFF

  // Test normal reading of bytes
  ram.WriteWord(0, wordToWrite)
  b, result := ram.ReadWord(0)
  if !result || b != wordToWrite {
    t.Fatalf("expected %#x got %#x", wordToWrite, b)
  }

  wordToWrite = 0xFDF67859
  ram.WriteWord(0, wordToWrite)
  b, result = ram.ReadWord(0)
  if !result || b != wordToWrite {
    t.Fatalf("expected %#x got %#x", wordToWrite, b)
  }

  // Test address out of range
  _, result = ram.ReadWord(80)
  if result {
    t.Fatal("Attempted to read out of range address.")
  }
}

func TestChecksum(t *testing.T) {
  // Test empty ram
  ram := Init(0)
  if ram.Checksum() != 0 {
    t.Fatal("zero byte ram checksum is not 0")
  }

  // Note: I used WolframAlpha to calculate these small checksums. They seem to
  // be correct.
  ram = Init(5)
  if ram.Checksum() != 10 {
    t.Fatal("checksum for 5 bytes of empty ram should be 10")
  }

  ram = Init(5)
  ram.WriteByte(0, 0x1)
  check := ram.Checksum()
  if check != 11 {
    t.Fatalf("expected checksum of %d; got: %d", 0x1^0, check)
  }

  ram = Init(5)
  ram.WriteByte(4, 11)
  check = ram.Checksum()
  if check != 21 {
    t.Fatalf("expected checksum of %d; got: %d", 21, check)
  }

  ram = Init(5)
  ram.WriteByte(3, 0x65)
  check = ram.Checksum()
  if check != 109 {
    t.Fatalf("expected checksum of %d; got: %d", 109, check)
  }

}

func TestTestFlag(t *testing.T) {
  // Setup
  // Initialize a single word for simplicity
  ram := Init(4)

  ram.WriteWord(0, 0x0001)
  if ram.TestFlag(0, 0) != true {
    t.Fatal("tested bit was false; expected true")
  }

  ram.WriteWord(0, 0x2)
  if ram.TestFlag(0, 0) != false {
    t.Fatal("tested bit was true; expected false")
  }
  if ram.TestFlag(0, 1) != true {
    t.Fatal("tested bit was false; expected true")
  }

  ram.WriteWord(0, 0xFFFFFFFF)
  for i := 0; i < 32; i++ {
    if ram.TestFlag(0, uint32(i)) == false {
      t.Fatal("tested bit was false; expected true")
    }
  }

  ram.WriteWord(0, 0x0)
  for i := 0; i < 32; i++ {
    if ram.TestFlag(0, uint32(i)) == true {
      t.Fatal("tested bit was true; expected false")
    }
  }

  // Test out-of-bounds read
  func() {
    defer func() {
      recover()
    }()

    ram.TestFlag(5, 0)
    t.Fatal("Shouldn't have been able to read this byte.")
  }()
}

func TestSetFlag(t *testing.T) {
  // Setup
  ram := Init(4)

  ram.WriteWord(0, 0xFFFFFFFF)
  ram.SetFlag(0, 0, true)
  if ram.TestFlag(0, 0) != true {
    t.Fatal("tested bit was not set to true")
  }

  ram.SetFlag(0, 0, false)
  if ram.TestFlag(0, 0) != false {
    t.Fatal("tested bit was not set to false")
  }

  ram.WriteWord(0, 0xFFFFFFFF)
  ram.SetFlag(0, 5, false)
  if ram.TestFlag(0, 5) != false {
    t.Fatal("tested bit was not set to false")
  }

  // Out-of-bounds testing
  func() {
    defer func() {
      recover()
    }()

    ram.SetFlag(5, 5, false)
    t.Fatal("should have panicked")
  }()
}

func TestExtractBits(t *testing.T) {
  result := ExtractBits(0xb5, 1, 3)
  if result != 0x04 {
    t.Fatalf("expected: 0x04 got: %#x", result)
  }

  result = ExtractBits(0xb5, 0, 3)
  if result != 0x05 {
    t.Fatalf("expected: 0x05 got: %#x", result)
  }

  // Implicity works
  result = ExtractBits(0xb5, 0, 33)
  if result != 0xb5 {
    t.Fatalf("expected: 0x05 got: %#x", result)
  }

  // Explicitly fails due to typing
  // ExtractBits(0xb5, -1, 33)
}
