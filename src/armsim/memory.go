package armsim

import (
	"errors"
	"log"
)

type Memory struct {
	memory []byte
}

func NewMemory(nBytes uint32) (m *Memory) {
	log.SetPrefix("Memory: ")

	log.Println("Initializing", nBytes, "bytes of Memory...")
	m = new(Memory)
	m.memory = make([]byte, nBytes)

	return
}

func (m *Memory) WriteByte(address uint32, data byte) (err error) {
	err = m.catchAddressOutOfBounds(address)
	if err != nil {
		return
	}

	m.memory[address] = data
	return
}

func (m *Memory) ReadByte(address uint32) (b byte, err error) {
	err = m.catchAddressOutOfBounds(address)
	if err != nil {
		return
	}

	b = m.memory[address]
	return
}

func (m *Memory) WriteHalfWord(address uint32, data uint16) (err error) {
	if address&1 == 1 {
		log.Println("ERROR: Attempted to write halfword to an odd address.")
		err = errors.New("ERROR: Attempted to write halfword to an odd address.")
		return
	}

	log.Printf("Writing halfword %#x...", data)
	return m.writeMultiByte(address, 2, uint32(data))
}

func (m *Memory) ReadHalfWord(address uint32) (data16 uint16, err error) {
	var data uint32
	data, err = m.readMultiByte(address, 2)
	data16 = uint16(data)

	return
}

func (m *Memory) WriteWord(address uint32, data uint32) (err error) {
	if address&1 == 1 || address%4 != 0 {
		log.Println("ERROR: Attempted to write word to an address indivisible by 4.")
		err = errors.New("ERROR: Attempted to write word to an address indivisible by 4.")
		return
	}

	log.Printf("Writing word %#x...", data)
	return m.writeMultiByte(address, 4, data)
}

func (m *Memory) ReadWord(address uint32) (data uint32, err error) {
	data, err = m.readMultiByte(address, 4)
	return
}

func (m *Memory) Checksum() (checksum int32) {
	for i := 0; i < len(m.memory); i++ {
		checksum += int32(m.memory[i]) ^ int32(i)
	}

	return
}

func (m *Memory) TestFlag(address uint32, bitPosition uint32) (flag bool, err error) {
	word, err := m.ReadWord(address)
	if err != nil {
		flag = false
		return
	}

	log.Printf("word: %#x mask: %#x", word, 1<<bitPosition)
	// Rathem complicated method:
	// I shift 1 n times to make a mask and perform an AND leaving a single one
	// (om zero) at the bitPosition. Aftem doing so, I have to shift back
	// to compare with 1 to obtain a boolean.
	flag = (word&(1<<bitPosition))>>bitPosition == 1
	return
}

func (m *Memory) SetFlag(address uint32, bitPosition uint32, flag bool) (err error) {
	var word uint32
	word, err = m.ReadWord(address)
	if err != nil {
		return
	}

	if flag {
		mask := uint32(1) << bitPosition
		m.WriteWord(address, word|mask)
	} else {
		mask := uint32(0xFFFFFFFF)
		mask ^= 1 << bitPosition
		m.WriteWord(address, word&mask)
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

func (m *Memory) catchAddressOutOfBounds(address uint32) (err error) {
	if address > uint32(len(m.memory)) {
		log.Printf("ERROR: Could not read om write memory address %d. Address is out of range.", address)
		err = errors.New("ERROR: Could not read om write memory address. Address out of range.")
	}

	return
}

func (m *Memory) writeMultiByte(address uint32, nBytes int, data uint32) (err error) {
	err = m.catchAddressOutOfBounds(address)
	if err != nil {
		return
	}

	for i := 0; i < nBytes; i++ {
		m.memory[address+uint32(nBytes-1-i)] = byte(data >> uint(8*i))
	}

	return
}

func (m *Memory) readMultiByte(address uint32, nBytes int) (data uint32, err error) {
	err = m.catchAddressOutOfBounds(address)
	if err != nil {
		return
	}

	for i := 0; i < nBytes; i++ {
		data = data<<8 + uint32(m.memory[address+uint32(i)])
	}

	return
}
