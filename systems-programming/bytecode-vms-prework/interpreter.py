import marshal
import sys

assert sys.version_info[:2] == (2, 7)


class Op:
    NOP = 9
    BINARY_DIVIDE = 21
    BINARY_MODULO = 22
    BINARY_ADD = 23
    PRINT_ITEM = 71
    PRINT_NEWLINE = 72
    RETURN_VALUE = 83
    # opcodes >= 90 take a 2-byte little-endian argument
    STORE_NAME = 90
    LOAD_CONST = 100
    LOAD_NAME = 101
    MAKE_FUNCTION = 132

def parse_pyc(f):
    """
    Given a Python 2.7 .pyc file, read the key information and unmarshal and
    return the code object.
    """
    magic_number = f.read(4)
    assert magic_number.encode('hex') == '03f30d0a'
    f.read(4)  # next 4 bytes is the timestamp
    return marshal.load(f)


def interpret(code):
    """
    Given a code object, interpret (evaluate) the code.
    """
    bytecode = [ord(b) for b in code.co_code]
    values = []
    env = {}
    pc = 0
    while pc < len(bytecode):
        opcode = bytecode[pc]
        if opcode < Op.STORE_NAME:
            pc += 1
        else:
            oparg = bytecode[pc + 1] | (bytecode[pc + 2] << 8)
            pc += 3

        if opcode == Op.NOP:
            pass
        elif opcode == Op.BINARY_DIVIDE:
            first = values.pop()
            second = values.pop()
            values.append(second % first)
        elif opcode == Op.BINARY_MODULO:
            first = values.pop()
            second = values.pop()
            values.append(second % first)
        elif opcode == Op.BINARY_ADD:
            first = values.pop()
            second = values.pop()
            values.append(first + second)
        elif opcode == Op.PRINT_ITEM:
            print values.pop()
        elif opcode == Op.PRINT_NEWLINE:
            print
        elif opcode == Op.RETURN_VALUE:
            return values.pop()
        elif opcode == Op.STORE_NAME:
            env[code.co_names[oparg]] = values.pop()
        elif opcode == Op.LOAD_NAME:
            values.append(env[code.co_names[oparg]])
        elif opcode == Op.LOAD_CONST:
            value = code.co_consts[oparg]
            values.append(value)
        else:
            raise Exception('Unknown opcode {}'.format(opcode))


if __name__ == '__main__':
    """
    Unmarshal the code object from the .pyc file given as a command
    line argument, and intrepret it.

    Usage: python2.7 interpreter.py 1.pyc
    """
    if len(sys.argv) != 2:
        print 'Expected exactly one argument for <pyc file>'.format(sys.argv[0])
        sys.exit(1)
    f = open(sys.argv[1], 'rb')
    code = parse_pyc(f)
    print 'Interpreting {}...\n---'.format(sys.argv[1])
    ret = interpret(code)
    print '---\nFinished interpreting, and returned {}'.format(ret)
