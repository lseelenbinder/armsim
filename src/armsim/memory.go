// Filename: memory.go
// Contents: Definition for Memory struct and related methods

package armsim

import (
	"errors"
	"log"
)

// A Memory holds a memory slice (a variable-length slice of bytes) used to
// implement RAM or other like memory structures.
type Memory struct {
	memory []byte
}

// Initializes a Memory
//
// Parameters:
//  nBytes - size of the new memory unit in bytes
//
// Returns: A pointer to the newly created Memory
func NewMemory(nBytes uint32) (m *Memory) {
	log.SetPrefix("Memory: ")

	log.Printf("Initializing %d bytes of Memory...", nBytes)
	m = new(Memory)
	// Initialize byte slice
	m.memory = make([]byte, nBytes)

	return
}

// Writes a byte of data to memory at a specified address.
//
// Parameters:
//  address - 32-bit address of write location in memory
//  data - byte of data to write
//
// Returns:
//  err - any error that may have occurred
func (m *Memory) WriteByte(address uint32, data byte) (err error) {
	err = m.catchAddressOutOfBounds(address)
	if err != nil {
		return
	}

	m.memory[address] = data
	return
}

// Reads a byte of data from memory at a specified address.
//
// Parameters:
//  address - 32-bit address of read location in memory
//
// Returns:
//  data - byte of data at address
//  err - any error that may have occurred
func (m *Memory) ReadByte(address uint32) (data byte, err error) {
	err = m.catchAddressOutOfBounds(address)
	if err != nil {
		return
	}

	data = m.memory[address]
	return
}

// Writes a halfword (16 bits) of data to memory at a specified address.
//
// Parameters:
//  address - 32-bit address of write location in memory, must be divisible by 2
//  data - halfword of data to write
//
// Returns:
//  err - any error that may have occurred
func (m *Memory) WriteHalfWord(address uint32, data uint16) (err error) {
	if address&1 == 1 {
		log.Println("ERROR: Attempted to write halfword to an odd address.")
		err = errors.New("ERROR: Attempted to write halfword to an odd address.")
		return
	}

	log.Printf("Writing halfword %#x...", data)
	return m.writeMultiByte(address, 2, uint32(data))
}

// Reads a halfword (16 bits) of data from memory at a specified address.
//
// Parameters:
//  address - 32-bit address of read location in memory
//
// Returns:
//  data - halfword of data at address
//  err - any error that may have occurred
func (m *Memory) ReadHalfWord(address uint32) (data uint16, err error) {
	var data32 uint32
	data32, err = m.readMultiByte(address, 2)
	data = uint16(data32)

	return
}

// Writes a word (32 bits) of data to memory at a specified address.
//
// Parameters:
//  address - 32-bit address of write location in memory, must be divisible by 4
//  data - word of data to write
//
// Returns:
//  err - any error that may have occurred
func (m *Memory) WriteWord(address uint32, data uint32) (err error) {
	if address&1 == 1 || address%4 != 0 {
		log.Println("ERROR: Attempted to write word to an address indivisible by 4.")
		err = errors.New("ERROR: Attempted to write word to an address indivisible by 4.")
		return
	}

	log.Printf("Writing word %#x...", data)
	return m.writeMultiByte(address, 4, data)
}

// Reads a word (32 bits) of data from memory at a specified address.
//
// Parameters:
//  address - 32-bit address of read location in memory
//
// Returns:
//  data - word of data at address
//  err - any error that may have occurred
func (m *Memory) ReadWord(address uint32) (data uint32, err error) {
	data, err = m.readMultiByte(address, 4)
	return
}

// Calculates a simple checksum based on the whole memory.
//
// Returns:
//  checksum - 32-bit integer
func (m *Memory) Checksum() (checksum int32) {
	for i := 0; i < len(m.memory); i++ {
		checksum += int32(m.memory[i]) ^ int32(i)
	}

	return
}

// Checks a specified bit in a word of data.
//
// Parameters:
//  address - 32-bit address of read location in memory
//  bitPosition - location in the word to test (based on endianness
//
// Returns:
//  flag - the on/off state of the tested bit
//  err - any error that may have occurred
func (m *Memory) TestFlag(address uint32, bitPosition uint32) (flag bool, err error) {
	word, err := m.ReadWord(address)
	if err != nil {
		return
	}

	// Somewhat complicated method:
	// I right-shift the word bitPosition times to get the tested bit in the far
	// right position and then bitwise-and it with 1 and compare with 1 to obtain
	// a boolean flag.
	flag = (word>>bitPosition)&1 == 1
	return
}

// Sets a specified bit in a word to a 1 or 0.
//
// Parameters:
//  address - 32-bit address of read location in memory
//  bitPosition - location in the word to test (based on endianness
//  flag - a boolean, determines whether to set a bit to 1 or 0
//
// Returns:
//  err - any error that may have occurred
func (m *Memory) SetFlag(address uint32, bitPosition uint32, flag bool) (err error) {
	word, err := m.ReadWord(address)
	if err != nil {
		return
	}

	// Build a mask
	var mask uint32
	if flag {
		// Need to set bitPosition bit to 1
		mask = 1 << bitPosition

		// Bitwise-or word with mask to set bit
		m.WriteWord(address, word|mask)
	} else {
		// Need to set bitPosition bit to 0
		mask = 0xFFFFFFFF
		mask ^= 1 << bitPosition

		// Bitwise-and word with mask to set bit
		m.WriteWord(address, word&mask)
	}

	log.Printf("word: %#x mask: %#x", word, mask)
	return
}

// Extracts bits from a word.
//
// Parameters:
//  word - a word of data
//  startBit - the least-significant bit to extract from
//  endBit - the most-significant bit to extract up to
//
// Returns: a new word containing the extracted bits and the rest set to zero
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

// Checks if an address is in the range of the memory. Returns nil or an error.
func (m *Memory) catchAddressOutOfBounds(address uint32) (err error) {
	if address > uint32(len(m.memory)) {
		log.Printf("ERROR: Could not read or write memory address %d. Address is out of range.", address)
		err = errors.New("ERROR: Could not read or write memory address. Address out of range.")
	}

	return
}

// Writes multiple bytes at a time in correct endianness
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

// Reads multiple bytes at a time in correct endianness
func (m *Memory) readMultiByte(address uint32, nBytes int) (data uint32, err error) {
	err = m.catchAddressOutOfBounds(address)
	if err != nil {
		return
	}

	for i := 0; i < nBytes; i++ {
		data <<= 8
		data |= uint32(m.memory[address+uint32(i)])
	}

	return
}
