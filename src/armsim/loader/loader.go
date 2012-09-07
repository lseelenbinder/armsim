package loader

import (
  "armsim/ram"
  "debug/elf"
  "encoding/binary"
  "log"
  "os"
)

func LoadELF(filePath string, memory *ram.RAM) bool {
  log.SetPrefix("Loader: ")
  // Attempt to open file
  log.Printf("Opening file %s", filePath)
  file, err := os.Open(filePath)
  if err != nil {
    log.Printf("Error reading file (perhaps it doesn't exist)...")
    return false
  }

  // Test magic bytes
  log.Println("Testing magic bytes...")
  magic := [4]byte{}
  err = binary.Read(file, binary.LittleEndian, &magic)
  log.Printf("Magic bytes: %s", magic)

  if err != nil || magic[0] != 0x7f || magic[1] != 'E' || magic[2] != 'L' || magic[3] != 'F' {
    log.Println("ELF magic bytes were incorrect.")
    return false
  }

  // Read ELF Header
  log.Println("Reading ELF header...")
  file.Seek(0, 0)
  elfHeader := new(elf.Header32)
  err = binary.Read(file, binary.LittleEndian, elfHeader)
  if err != nil {
    log.Println("Error reading ELF header...")
    return false
  }

  log.Printf("Program header offset: %d", elfHeader.Phoff)
  log.Printf("# of program header entires: %d", elfHeader.Phnum)

  // Seek to Program Header start
  file.Seek(int64(elfHeader.Phoff), 0)

  // Read Program Headers
  log.Println("Reading program headers...")
  pHeader := new(elf.Prog32)
  for i := 0; uint16(i) < elfHeader.Phnum; i++ {
    // Seek to program header
    offset := int64(elfHeader.Phoff) + int64(i)*int64(elfHeader.Phentsize)
    file.Seek(offset, 0)

    // Read program header
    err = binary.Read(file, binary.LittleEndian, pHeader)
    if err != nil {
      log.Println("Error reading program header %d...", i)
      return false
    }

    log.Printf("Reading program header %d of %d - Offset: %d, Size: %d, Address: %d", i+1, elfHeader.Phnum, pHeader.Off, pHeader.Filesz, pHeader.Vaddr)
    // Seek to offset
    file.Seek(int64(pHeader.Off), 0)

    // Read to RAM
    b := make([]byte, 1)
    var i uint32 = 0
    for ; i < pHeader.Filesz; i++ {
      file.Read(b)
      memory.WriteByte(pHeader.Vaddr+i, b[0])
    }
  }

  file.Close()
  return true
}
