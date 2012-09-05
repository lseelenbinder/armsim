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
