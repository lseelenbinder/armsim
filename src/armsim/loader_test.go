package armsim

import "testing"

func TestLoadELF(t *testing.T) {
	var memory *Memory
	var checksum int32

	// Setup RAM
	memory = NewMemory(32 * 1024)

	// Load Non-existent Test File
	err := LoadELF("asdfasdfaitheirhasdifhadf", memory)
	if err == nil {
		t.Fatalf("should have failed file error")
	}

	// Load Non-ELF Test File
	err = LoadELF("loader.go", memory)
	if err == nil {
		t.Fatalf("should have failed with magic number err")
	}

	// Load Test File 1
	memory = NewMemory(32 * 1024)
	t.Log("Checksum of empty RAM: ", memory.Checksum())
	err = LoadELF("../../test/test1.exe", memory)
	checksum = memory.Checksum()
	if err != nil || checksum != 536861081 {
		t.Fatalf("Checksum did not match for test1.exe. Expected 536861081. Got %d", checksum)
	}

	// Clear RAM
	memory = NewMemory(32 * 1024)

	// Load Test File 2
	err = LoadELF("../../test/test2.exe", memory)
	checksum = memory.Checksum()
	if err != nil || checksum != 536864433 {
		t.Fatalf("Checksum did not match for test2.exe. Expected 536864433. Got %d", checksum)
	}

	// Clear RAM
	memory = NewMemory(32 * 1024)

	// Load Test File 3
	err = LoadELF("../../test/test3.exe", memory)
	checksum = memory.Checksum()
	if err != nil || checksum != 536861199 {
		t.Fatalf("Checksum did not match for test3.exe. Expected 536861199. Got %d", checksum)
	}

	// Clear RAM
	memory = NewMemory(8)

	// Load Test 3 into insuffcient memory
	err = LoadELF("../../test/test3.exe", memory)
	if err == nil {
		t.Fatal("Should have failed with insuffcient memory error.")
	}
}
