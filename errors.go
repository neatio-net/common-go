package common

import (
	"fmt"
)

type StackError struct {
	Err   interface{}
	Stack []byte
}

func (se StackError) String() string {
	return fmt.Sprintf("Error: %v\nStack: %s", se.Err, se.Stack)
}

func (se StackError) Error() string {
	return se.String()
}

func PanicSanity(v interface{}) {
	panic(Fmt("Panicked on a Sanity Check: %v", v))
}

func PanicCrisis(v interface{}) {
	panic(Fmt("Panicked on a Crisis: %v", v))
}

func PanicConsensus(v interface{}) {
	panic(Fmt("Panicked on a Consensus Failure: %v", v))
}

func PanicQ(v interface{}) {
	panic(Fmt("Panicked questionably: %v", v))
}
