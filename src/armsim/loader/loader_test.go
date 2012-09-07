package loader

import (
  "armsim/ram"
  "testing"
)

func TestLoadELF(t *testing.T) {
  // Setup RAM
  memory := ram.Init(32 * 1024)

  // Load Non-existent Test File
  success := LoadELF("asdfasdfaitheirhasdifhadf", &memory)
  checksum := memory.Checksum()
  if success {
    t.Fatalf("should have failed without success")
  }

  // Load Non-ELF Test File
  success = LoadELF("loader.go", &memory)
  checksum = memory.Checksum()
  if success {
    t.Fatalf("should have failed without magic number success")
  }

  // Load Test File 1
  memory = ram.Init(32 * 1024)
  success = LoadELF("../../../test/test1.exe", &memory)
  checksum = memory.Checksum()
  if !success || checksum != 536861081 {
    t.Fatalf("Checksum did not match for test1.exe. Expected 536861081. Got %d", checksum)
  }

  // Clear RAM
  memory = ram.Init(32 * 1024)

  // Load Test File 2
  LoadELF("../../../test/test2.exe", &memory)
  checksum = memory.Checksum()
  if !success || checksum != 536864433 {
    t.Fatalf("Checksum did not match for test2.exe. Expected 536864433. Got %d", checksum)
  }

  // Clear RAM
  memory = ram.Init(32 * 1024)

  // Load Test File 3
  LoadELF("../../../test/test3.exe", &memory)
  checksum = memory.Checksum()
  if !success || checksum != 536861199 {
    t.Fatalf("Checksum did not match for test3.exe. Expected 536861199. Got %d", checksum)
  }
}
