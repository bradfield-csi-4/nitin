package vm

import "fmt"

const (
	_ = byte(iota)
	LOAD
	STORE
	ADD
	SUB
	HALT = 0xFF
)

const instructionSize = 3

func compute(memory []byte) {
	registers := [3]byte{8}
	pc := &registers[0]

	var arg1 *byte
	var arg2 *byte

	for {
		// FETCH
		instruction := memory[*pc : *pc+instructionSize]

		// DECODE
		op := &instruction[0]
		arg1 = &instruction[1]
		arg2 = &instruction[2]

		// EXECUTE
		if *op == HALT {
			return
		}

		switch *op {
		case LOAD:
			registers[*arg1] = memory[*arg2]
		case STORE:
			memory[*arg2] = registers[*arg1]
		case ADD:
			registers[*arg1] += registers[*arg2]
		case SUB:
			registers[*arg1] -= registers[*arg2]
		default:
			panic(fmt.Sprintf("Invalid opcode: 0x%02x", *op))
		}

		*pc += instructionSize
	}
}
