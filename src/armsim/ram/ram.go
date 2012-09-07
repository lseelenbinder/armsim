package ram

import (
  "errors"
  "log"
)

type RAM struct {
  memory []byte
}

func New(nBytes uint32) (ram *RAM) {
  log.SetPrefix("RAM: ")

  log.Println("Initializing", nBytes, "bytes of RAM...")
  ram = new(RAM)
  ram.memory = make([]byte, nBytes)

  return
}

func (r *RAM) WriteByte(address uint32, data byte) (err error) {
  err = r.catchAddressOutOfBounds(address)
  if err != nil {
    return
  }

  r.memory[address] = data
  return
}

func (r *RAM) ReadByte(address uint32) (b byte, err error) {
  err = r.catchAddressOutOfBounds(address)
  if err != nil {
    return
  }

  b = r.memory[address]
  return
}

func (r *RAM) WriteHalfWord(address uint32, data uint16) (err error) {
  if address&1 == 1 {
    log.Println("ERROR: Attempted to write halfword to an odd address.")
    err = errors.New("ERROR: Attempted to write halfword to an odd address.")
    return
  }

  log.Printf("Writing halfword %#x...", data)
  return r.writeMultiByte(address, 2, uint32(data))
}

func (r *RAM) ReadHalfWord(address uint32) (data16 uint16, err error) {
  var data uint32
  data, err = r.readMultiByte(address, 2)
  data16 = uint16(data)

  return
}

func (r *RAM) WriteWord(address uint32, data uint32) (err error) {
  if address&1 == 1 || address%4 != 0 {
    log.Println("ERROR: Attempted to write word to an address indivisible by 4.")
    err = errors.New("ERROR: Attempted to write word to an address indivisible by 4.")
    return
  }

  log.Printf("Writing word %#x...", data)
  return r.writeMultiByte(address, 4, data)
}

func (r *RAM) ReadWord(address uint32) (data uint32, err error) {
  data, err = r.readMultiByte(address, 4)
  return
}

func (r *RAM) Checksum() int32 {
  var checksum int32

  for i := 0; i < len(r.memory); i++ {
    checksum += int32(r.memory[i]) ^ int32(i)
  }

  return checksum
}

func (r *RAM) TestFlag(address uint32, bitPosition uint32) (flag bool, err error) {
  word, err := r.ReadWord(address)
  if err != nil {
    flag = false
    return
  }

  log.Printf("word: %#x mask: %#x", word, 1<<bitPosition)
  // Rather complicated method:
  // I shift 1 n times to make a mask and perform an AND leaving a single one
  // (or zero) at the bitPosition. After doing so, I have to shift back
  // to compare with 1 to obtain a boolean.
  flag = (word&(1<<bitPosition))>>bitPosition == 1
  return
}

func (r *RAM) SetFlag(address uint32, bitPosition uint32, flag bool) (err error) {
  var word uint32
  word, err = r.ReadWord(address)
  if err != nil {
    return
  }

  if flag {
    mask := uint32(1) << bitPosition
    r.WriteWord(address, word|mask)
  } else {
    mask := uint32(0xFFFFFFFF)
    mask ^= 1 << bitPosition
    r.WriteWord(address, word&mask)
  }

  log.Printf("word: %#x mask: %#x", word, 1<<bitPosition)
  return
}

func ExtractBits(word uint32, startBit uint32, endBit uint32) uint32 {
  mask := uint32(0)
  for i := 31; i >= 0; i-- {
    mask <<= 1
    if startBit <= uint32(i) && uint32(i) < endBit {
      mask |= 1
    }
  }

  log.Printf("word: %#x mask: %#x", word, mask)
  return word & mask
}

// Helpers

func (r *RAM) catchAddressOutOfBounds(address uint32) (err error) {
  if address > uint32(len(r.memory)) {
    log.Printf("ERROR: Could not read or write memory address %d. Address is out of range.", address)
    err = errors.New("ERROR: Could not read or write memory address. Address out of range.")
  }

  return
}

func (r *RAM) writeMultiByte(address uint32, nBytes int, data uint32) (err error) {
  err = r.catchAddressOutOfBounds(address)
  if err != nil {
    return
  }

  for i := 0; i < nBytes; i++ {
    r.memory[address+uint32(nBytes-1-i)] = byte(data >> uint(8*i))
  }

  return
}

func (r *RAM) readMultiByte(address uint32, nBytes int) (data uint32, err error) {
  err = r.catchAddressOutOfBounds(address)
  if err != nil {
    return
  }

  for i := 0; i < nBytes; i++ {
    data = data<<8 + uint32(r.memory[address+uint32(i)])
  }

  return
}
