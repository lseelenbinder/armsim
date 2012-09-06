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
  defer r.catchAddressOutOfBounds(address)

  r.memory[address] = data
  return true
}

func (r *RAM) ReadByte(address uint32) (byte, bool) {
  defer r.catchAddressOutOfBounds(address)

  return r.memory[address], true
}

func (r *RAM) WriteHalfWord(address uint32, data uint16) bool {
  defer r.catchAddressOutOfBounds(address)

  if address & 0x00000001 == 1 {
    log.Println("ERROR: Attempted to write halfword to an odd address.")
    return false
  }

  var b1, b2 byte = byte(data >> 8), byte(data)
  log.Printf("Writing halfword %#x as %#x & %#x...", data, b1, b2)

  r.memory[address] = b1
  r.memory[address+1] = b2
  return true
}

func (r *RAM) ReadHalfWord(address uint32) (uint16, bool) {
  defer r.catchAddressOutOfBounds(address)

  data := uint16(r.memory[address+1]) << 8 + uint16(r.memory[address])

  return data, true
}

func (r *RAM) catchAddressOutOfBounds(address uint32) bool {
    recover()
    if address > uint32(len(r.memory)) {
      log.Printf("ERROR: Could not read memory address %d. Address is out of range.", address)
      return false
    }

    return true
}
