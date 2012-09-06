package ram

import "log"

type RAM struct {
  memory []byte
}

func Init(nBytes uint32) RAM {
  log.SetPrefix("RAM: ")
  log.Println("Initializing", nBytes, "bytes of RAM...")
  ram := RAM{make([]byte, nBytes)}

  return ram
}

func (r *RAM) WriteByte(address uint32, data byte) bool {
  log.SetPrefix("RAM: ")

  defer func() {
    recover()
    if address > uint32(len(r.memory)) {
      log.Println("ERROR: Could not write byte. Address is out of range.")
    }
  }()

  r.memory[address] = data
  return true
}

func (r *RAM) ReadByte(address uint32) (byte, bool) {
  log.SetPrefix("RAM: ")

  defer func() {
    recover()
    if address > uint32(len(r.memory)) {
      log.Println("ERROR: Could not read byte. Address is out of range.")
    }
  }()

  return r.memory[address], true
}
