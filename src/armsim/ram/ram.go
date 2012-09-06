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
  if address & 1 == 1 {
    log.Println("ERROR: Attempted to write halfword to an odd address.")
    return false
  }

  log.Printf("Writing halfword %#x...", data)
  return r.writeMultiByte(address, 2, uint32(data))
}

func (r *RAM) ReadHalfWord(address uint32) (uint16, bool) {
  data, success := r.readMultiByte(address, 2)
  return uint16(data), success
}

func (r *RAM) WriteWord(address uint32, data uint32) bool {
  if address & 1 == 1 || address & 3 == 1 {
    log.Println("ERROR: Attempted to write word to an address indivisible by 4.")
    return false
  }

  log.Printf("Writing word %#x...", data)
  return r.writeMultiByte(address, 4, data)
}

func (r *RAM) ReadWord(address uint32) (uint32, bool) {
  return r.readMultiByte(address, 4)
}

// Helpers

func (r *RAM) catchAddressOutOfBounds(address uint32) bool {
    if err := recover(); err != nil {
      if address > uint32(len(r.memory)) {
        log.Printf("ERROR: Could not read or write memory address %d. Address is out of range.", address)
      } else {
        log.Fatalln("ERROR: Unknown Error")
      }
      return false
    }
    return true
}

func (r *RAM) writeMultiByte(address uint32, nBytes int, data uint32) bool {
  defer r.catchAddressOutOfBounds(address)

  for i := 0; i < nBytes; i++ {
    r.memory[address + uint32(nBytes - 1 - i)] = byte(data >> uint(8 * i))
  }

  return true
}

func (r *RAM) readMultiByte(address uint32, nBytes int) (uint32, bool) {
  defer r.catchAddressOutOfBounds(address)

  var data uint32 = 0
  for i := 0; i < nBytes; i++ {
    data = data << 8 + uint32(r.memory[address + uint32(i)])
  }

  return data, true
}
