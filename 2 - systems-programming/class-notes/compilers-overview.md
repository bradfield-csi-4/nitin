## Objectives

By the end of this session, you should understand:

- The high level steps involved in a language implementation
- Different ways that a program might be represented as it's transformed and executed, benefits of each representation
- Uses cases of compilers / interpreters ideas outside of implementing a language

## Agenda

After a brief discussion about the reading, we will use the majority of the session to discuss the three implementation exercises.

If time allows, we can also try to live-code a "constant folding" optimization.

**Discussion**

- Motivation / use cases for learning about compilers?
	- Is anyone interested in working on a language implementation?
		- Clojure
		- Formal methods / formal logic
			- Recommendation for Coq: Software Foundations
		- DSLs
		- Contributing to Go
			- https://go.dev/doc/contribute
	- Applications outside of programming languages?
		- Editors
			- Structured editing
			- LSP: language server protocol
				- gopls is an example for Go
					- runs in bg
					- listens over unix socket for json
					- runs in go module, has access to go path
					- can answer questions like:
						- "look up module"
						- "look up keyword"
						- "get doc for this"
						- "what's function signature for this"
			- Syntax highlighting
				- Can use regex-based approach
				- Can parse it and then use result to annotate
				- Incremental parsing: https://github.com/tree-sitter/tree-sitter
			- Detecting missing imports
				- goimports
				- Open question that came up: why are they using `context.TODO`?
					- See https://go.dev/blog/context
		- Performance benefits by understanding compilation
			- e.g. understanding what keywords like `volatile` do, understanding what operations prevent the compiler from making certain optimizations
			- Understand what compiler optimizations are possible, what they do
			- Look at output and understand performance better
		- Can understand syntactic sugar, features, etc.

- What are the steps involved in going from a program in a high-level language to a running program?
	- Our goal is to:
		- List relevant steps
			- Scanning / lexing (tokenize input)
			- Parse it (build AST, put structure to tokens)
			- (It depends at that point, languages might do different things)
				- e.g. spit out low level code
				- go to intermediate representation
			- Semantic analysis
				- Open question: is it possible to detect deadlocks at compile time?
			- Intermediate representation analysis
			- Optimization
				- e.g. constant folding, common subexpression elimination
			- Code generation
			- Optional "transpiling" in a couple places
			- Optional "evaluating" directly
		- State the input / output of that step
		- Describe what it does / how it works at a high level

- More open questions
	- Is it possible to parallelize a compiler?
		- See https://gcc.gnu.org/wiki/ParallelGcc
	- Is it possible to do "eval" iteratively?
		- Tricky because a tree is naturally a recursive data structure
		- For an expression, you can parse it iteratively using the ["shunting yard algorithm"](https://en.wikipedia.org/wiki/Shunting_yard_algorithm)

**Review of Implementation Exercises**

## Appendix

**Some notes/references from Eric Ihli**

- LSP: https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/
- How to write code the compiler can actually optimize: https://www.youtube.com/watch?v=GPpD4BBtA1Y

Since Asa Needle mentioned an interest in things like Coq and TLA+:

- https://the-little-prover.github.io/
- https://corecursive.com/023-little-typer-and-pie-language/
- Using TLA+ to find Goroutine bugs: https://www.hillelwayne.com/post/tla-golang/

More resources:

- Tracing JIT: [Pixie, a LISP written in RPython](https://www.youtube.com/watch?v=1AjhFZVfB9c)
- Related to the "constantFold" step: https://nanopass.org/ A compiler framework built around many tiny optimization steps.

**Aside: how would you type-check generics in Java?**

- Don't know how it's done in Java, but one way you MIGHT do it (and possibly the way it's done in C++) is to make a copy of the function / class / etc. for each actual type that's instantiated

```java
class List<T> {
	void add(T item);
}

// You can imagine that if we use List<String> somewhere, then there's a
// preprocessing step that creates a copy of the class that looks something
// like this:
class List_String {
	void add(String item)
}

// Similarly, if we use List<Integer> somewhere, then maybe the preprocessing
// step creates this too:
class List_Integer {
	void add(Integer item)
}

// Now it becomes easy to check that we didn't use a List<Integer> where
// we expected a List<String>:
class MyClass {
	void handle(List_String l);
}

MyClass c = new MyClass();
List_Integer l = new List_Integer();
c.handle(l); // compile error!

// WARNING: There's probably a lot more to this, and what's described
// above probably doesn't work in practice; for example, what if you have
// some subclass of String called SpecialString? Is it valid to pass a
// List<SpecialString> to a method that expects a List<String>?
```
