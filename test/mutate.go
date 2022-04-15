package test

import (
	. "github.com/neatlab/common-go"
)

func MutateByteSlice(bytez []byte) []byte {

	if len(bytez) == 0 {
		panic("Cannot mutate an empty bytez")
	}

	mBytez := make([]byte, len(bytez))
	copy(mBytez, bytez)
	bytez = mBytez

	switch RandInt() % 2 {
	case 0:
		bytez[RandInt()%len(bytez)] += byte(RandInt()%255 + 1)
	case 1:
		pos := RandInt() % len(bytez)
		bytez = append(bytez[:pos], bytez[pos+1:]...)
	}
	return bytez
}
