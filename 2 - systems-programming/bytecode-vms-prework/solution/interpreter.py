#!/usr/bin/env python2.7

import marshal
import operator
import sys


class Op:
    BINARY_MULTIPLY   = 20
    BINARY_DIVIDE     = 21
    BINARY_MODULO     = 22
    BINARY_ADD        = 23
    BINARY_SUBTRACT   = 24
    INPLACE_ADD       = 55
    GET_ITER          = 68
    PRINT_ITEM        = 71
    PRINT_NEWLINE     = 72
    RETURN_VALUE      = 83
    POP_BLOCK         = 87

    # Opcodes >= 90 take an argument
    HAVE_ARGUMENT     = 90

    STORE_NAME        = 90
    FOR_ITER          = 93
    LOAD_CONST        = 100
    LOAD_NAME         = 101
    COMPARE_OP        = 107
    JUMP_ABSOLUTE     = 113
    POP_JUMP_IF_FALSE = 114
    LOAD_GLOBAL       = 116
    SETUP_LOOP        = 120
    LOAD_FAST         = 124
    CALL_FUNCTION     = 131
    MAKE_FUNCTION     = 132

BINARY_OPS = {
    Op.BINARY_MULTIPLY: operator.mul,
    Op.BINARY_DIVIDE:   operator.div,
    Op.BINARY_MODULO:   operator.mod,
    Op.BINARY_ADD:      operator.add,
    Op.BINARY_SUBTRACT: operator.sub,

    # TODO: Actually do these in-place
    Op.INPLACE_ADD:     operator.add,
}

COMPARE_OPS = {
    2: operator.eq,
}

class Frame:
    def __init__(self, code, m_globals):
        self.code = code
        self.bytecode = [ord(byte) for byte in code.co_code]
        self.m_globals = m_globals

        self.pc = 0
        self.locals = [None] * len(code.co_varnames)
        self.stack = []

    def is_at_end(self):
        return self.pc >= len(self.bytecode)

    def fetch_op(self):
        self.pc += 1
        return self.bytecode[self.pc - 1]

    def fetch_arg(self):
        self.pc += 2
        b0 = self.bytecode[self.pc - 2]
        b1 = self.bytecode[self.pc - 1]
        return b0 + 256 * b1

    def run(self):
        while not self.is_at_end():
            op = self.fetch_op()
            if op >= Op.HAVE_ARGUMENT:
                arg = self.fetch_arg()

            if op in BINARY_OPS:
                b = self.stack.pop()
                a = self.stack.pop()
                self.stack.append(BINARY_OPS[op](a, b))
            elif op == Op.GET_ITER:
                self.stack.append(iter(self.stack.pop()))
            elif op == Op.PRINT_ITEM:
                print self.stack.pop()
            elif op == Op.PRINT_NEWLINE:
                print
            elif op == Op.RETURN_VALUE:
                return self.stack.pop()
            elif op == Op.POP_BLOCK:
                # TODO: Actually handle blocks
                pass
            elif op == Op.STORE_NAME:
                self.m_globals[self.code.co_names[arg]] = self.stack.pop()
            elif op == Op.FOR_ITER:
                try:
                    self.stack.append(next(self.stack[-1]))
                except StopIteration:
                    self.stack.pop()
                    self.pc += arg
            elif op == Op.LOAD_CONST:
                self.stack.append(self.code.co_consts[arg])
            elif op == Op.LOAD_NAME:
                # TODO: Handle the case where LOAD_NAME isn't loading a global
                self.stack.append(self.m_globals[self.code.co_names[arg]])
            elif op == Op.COMPARE_OP:
                b = self.stack.pop()
                a = self.stack.pop()
                self.stack.append(COMPARE_OPS[arg](a, b))
            elif op == Op.JUMP_ABSOLUTE:
                self.pc = arg
            elif op == Op.POP_JUMP_IF_FALSE:
                if not self.stack.pop():
                    self.pc = arg
            elif op == Op.LOAD_GLOBAL:
                self.stack.append(self.m_globals[self.code.co_names[arg]])
            elif op == Op.SETUP_LOOP:
                # TODO: Actually handle blocks
                pass
            elif op == Op.LOAD_FAST:
                self.stack.append(self.locals[arg])
            elif op == Op.CALL_FUNCTION:
                # TODO: Handle keyword args
                assert arg < 256

                # Get positional args from stack
                pos_args = []
                for i in range(arg):
                    pos_args.append(self.stack.pop())
                pos_args = pos_args[::-1] # Reverse

                code = self.stack.pop()
                if hasattr(code, 'co_code'):
                    # Pass positional args into the function's context
                    frame = Frame(code, self.m_globals)
                    for i in range(arg):
                        frame.locals[i] = pos_args[i]

                    self.stack.append(frame.run())
                else:
                    # TODO: Handle builtins in a less hacky way
                    self.stack.append(apply(code, pos_args))

            elif op == Op.MAKE_FUNCTION:
                self.stack.append(self.stack.pop())
            else:
                raise Exception('Unknown opcode %s' % op)

def load_code(fname):
    f = open(fname, "rb")
    f.read(8) # skip over magic and moddate
    return marshal.load(f)

if __name__ == '__main__':
    if len(sys.argv) != 2:
        print 'Usage: python2.7', sys.argv[0], '<pyc file>'
        sys.exit(1)

    # 8.pyc relies on the 'range' builtin
    m_globals = {"range": range}

    code = load_code(sys.argv[1])
    frame = Frame(code, m_globals)
    frame.run()
