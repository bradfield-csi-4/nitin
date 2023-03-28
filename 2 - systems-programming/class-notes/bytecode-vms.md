## Objectives

At this point, you've seen scanning and parsing in the context of a real language implementation (via Go's `ast` and `parser` packages). You've also worked directly with x86-64 assembly (in _Introduction to Computer Systems_).

However, we've only been working with toy VMs, with very simplified instruction sets. The goal of this session is to explore the concepts we've been discussing in the context of a production language.

By the end of this session, you should understand:

- Structure of Python 2.7 "compiled" bytecode (`.pyc` files)
- Basic architecture of a Python bytecode interpreter
- How various Python features are implemented

## Agenda

This will be an open-ended, exploratory session focused on going over the prework. Depending on time/interest, we can then implement various stretch goals as a group, or explore more details of the CPython interpreter.

Possible stretch goals:
- generators
- closures

**Exercise**

- Quick recap of the goal
	- Get familiar with a real-world virtual machine
	- To build a simple Python interpreter
	- Input: "machine code" for the virtual machine (aka bytecode—it's just raw bytes)
	- Output: instructions that the CPU actually runs
		- Question: are we converting Python bytecode into x86 machine code at some point in the process?
			- Don't think so! We are NOT compiling down to machine code at any actual point; our actual machine takes the output of virtual machine
			- `python` itself is the only thing that has x86 bytecode
			- The bytecode is an input to the `python` program (a C program)

- `pyc` files in the wild?
	- When you run Python, it'll compile your code down into these `pyc` objects; you'll see them in some `__pycache__`
	- Every time you run your program, you can save the work of compiling

- structure of the file
	- magic number (identifies version of Python)
	- modification date (can tell us if we can avoid recompiling source)
	- rest of bytes are code object that has been "marshalled" using the [`marshal` library](https://docs.python.org/3/library/marshal.html)
	- fields of the code object:
		- See https://docs.python.org/3/reference/datamodel.html#objects-values-and-types

- how did it go?
	- went smoothly until the iterator one `8.pyc`
		- recursive functions wasn't too bad
		- resolving environments was somewhat tricky
	- up to 7 with just `interpret` function, but then realized wanted actual VM class with helper `Frame` classes, etc.
	- hacked the crap out of it, second Eric's point ("oh no, reached point where warning about abstraction gets finicky")
	- did whole thing without reading article: trickiest was "block" (implemented it with extra function call)

- what parts felt new?
	- higher order constructs like functions, blocks, iterators
	- nested scopes (starting at locals, going to globals, then looking at builtins)
		- nested locals for functions inside of functions
	- ceval code checked at various levels
		- was weird to read C and then think "what's Python like"
	- was tempted to just read interpreter.py to get started
	- tried Python3, header changed, timestamp format changed, but otherwise still pretty much same

- what was the process for working on this?