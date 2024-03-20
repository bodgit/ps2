// Package ecc implements the ECC algorithm used by PlayStation 2 memory cards.
package ecc

import "hash"

const (
	// BlockSize is the preferred block size.
	BlockSize = 1
	// Size is the number of ECC bytes returned.
	Size = 3
)

var parityTable, columnParityMasks [256]byte

func parity(b byte) byte {
	b ^= b >> 1
	b ^= b >> 2
	b ^= b >> 4

	return b & 1
}

func init() {
	for i := range parityTable {
		parityTable[i] = parity(byte(i))
	}

	for i := range columnParityMasks {
		for j, v := range []byte{0x55, 0x33, 0x0f, 0x00, 0xaa, 0xcc, 0xf0} {
			columnParityMasks[i] |= parityTable[i&int(v)] << j
		}
	}
}

type ecc struct {
	cp, lp0, lp1 byte
}

func (e *ecc) BlockSize() int { return BlockSize }

func (e *ecc) Reset() {
	e.cp = 0x77
	e.lp0 = 0x7f
	e.lp1 = 0x7f
}

func (e *ecc) Size() int { return Size }

func (e *ecc) Sum(data []byte) []byte {
	return append(data, e.cp, e.lp0&0x7f, e.lp1)
}

func (e *ecc) Write(p []byte) (int, error) {
	for _, i := range p {
		e.cp ^= columnParityMasks[i]

		if parityTable[i] == 1 {
			e.lp0 ^= ^i
			e.lp1 ^= i
		}
	}

	return len(p), nil
}

// New returns a hash.Hash implementation that computes the ECC from every
// byte written to it.
func New() hash.Hash {
	e := new(ecc)
	e.Reset()

	return e
}
