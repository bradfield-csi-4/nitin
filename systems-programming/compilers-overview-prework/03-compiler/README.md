## Overview

The goal of this exercise is to convert a simple Go function into bytecode, and then execute the bytecode in a virtual machine.

## Simplifying Assumptions

Solving the problem in full generality (i.e. supporting an arbitrary Go function) is well outside the scope of this exercise. Instead, we recommend making the following simplifying assumptions:

- The input will always be a file with a package declaration followed by a single function
- The single function will always have the prototype `func f(x, y byte) byte`
- Every value will have the type `byte`
- There are no recursive function calls (though you may consider adding this as a stretch goal)

## Pipeline

The source code will pass through the following steps:

- `parse` uses standard library functions to convert raw source code into an abstract syntax tree
- `compile`, **which you are responsible for implementing**, converts the abstract syntax tree into assembly code
- `assemble` converts assembly code into bytecode
- `runVM` executes the bytecode with the given "input data" and returns the result

## VM Details

There are a few key differences between this VM and the one you implemented for the first session of _Introduction to Computer Systems_.

**Input / Output**

You can assume that as before, the input parameters will always be placed at memory addresses `1` and `2`. In other words, the function `func f(x, y byte) byte` should read `x` and `y` from memory addresses `1` and `2`.

Once the function is done, the return value should be stored in memory address `0`.

**Stack Machine**

The biggest difference is that this one is a "stack machine" rather than a "register machine". In a register machine, arithmetic instructions might look like `add r1 r2`. In a stack machine, arithmetic instructions don't take any arguments; instead, they get their arguments from the stack.

For example, if the stack is `[1 5 2]`, then the `sub` instruction will:

- Pop two values `5` and `2` from the stack
- Perform the subtraction to get `3` (note the ordering)
- Push the result back onto the stack, which will then be `[1 3]`

In addition to arithmetic instructions not taking any arguments, there are also three new instructions involving the stack:

- `push <addr>` will push the contents at memory location `addr` onto the stack
- `pushi <val>` will push the literal value `val` onto the stack
- `pop <addr>` will pop a value from the stack and save it at memory location `addr`

**Comparison Instructions**

There are several new comparison instructions such as `eq`, `lt`. These instructions will pop two values from the stack, perform a comparison, and then push either `0` or `1` back onto the stack, depending on the result.

**Jumps / Labels**

The new `jeqz <addr>` instruction (jump if equal to zero) will:
- Pop a value from the stack (the value will be removed from the stack)
- Check if the value is equal to `0`
- If so, the VM will jump to `addr`; if not, the VM will continue executing the next instruction in order

At the bytecode level, both `Jump` and `Jeqz` absolute (rather than relative) offsets.

At the assembly level, `jump` and `jeqz` take labels (rather than raw addresses), which should make it easier to generate assembly code with jumps. For example:

```
label loop_start
...
jump loop_start
```

## Suggestions

The full version of this problem is VERY difficult. We recommend that you start with a VERY restricted version of the problem, and only gradually add functionality. For example, you might want to start by assuming that the function consists of a single return expression containing only x, y, parentheses, and integer literals, such as:

```go
func f(x, y byte) byte {
    return 2 * (x + 3) * (y + 4)
}
```

Once you have something basic working, you can start introducing more functionality (e.g. multiple statements, if / else, local variables, for loops) to pass more of the test cases.
