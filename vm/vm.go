package vm

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

	var r1 *byte
	var r2 *byte
	var addr *byte

OuterLoop:
	for {
		// FETCH
		instruction := memory[*pc : *pc+instructionSize]

		// DECODE
		opCode := &instruction[0]
		r1 = &instruction[1]
		r2 = &instruction[2]
		addr = &instruction[2]

		// EXECUTE
		switch *opCode {
		case LOAD:
			registers[*r1] = memory[*addr]
		case STORE:
			memory[*addr] = registers[*r1]
		case ADD:
			registers[*r1] += registers[*r2]
		case SUB:
			registers[*r1] -= registers[*r2]
		case HALT:
			break OuterLoop
		}

		*pc += instructionSize
	}
}
