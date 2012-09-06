package ram

import "testing"

func TestInit(t *testing.T) {
  // Test the RAM initializer at 32k
  runInit(t, 32 * 1024)
  // Test the RAM initializer at 32m
  runInit(t, 32 * 1024 * 1024)
}

// Helper method for the RAM initializer testing code
func runInit(t *testing.T, nBytes uint32) {
  // Setup
  ram := Init(nBytes)

  // Ensure RAM is the same length as nBytes
  if (uint32(len(ram.memory)) != nBytes) {
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
  var byteToWrite byte = 95;

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

  // Test normal reading of bytes
  ram.WriteByte(0, 95)
  b, result := ram.ReadByte(0)
  if result && b != 95 {
    t.Fatal("Byte not read properly.")
  }

  // Test address out of range
  ram.WriteByte(1, 95)
  b, result = ram.ReadByte(80)
  if result {
    t.Fatal("Attempted to read out of range address.")
  }
}
