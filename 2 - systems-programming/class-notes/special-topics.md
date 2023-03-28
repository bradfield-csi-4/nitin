## Objectives

This is an informal session where we will:

- Discuss compiler bootstrapping and Ken Thompson's "Trusting Trust" attack
- Hear a few short talks from students about topics of interest
- Address lingering questions from the module
- Do a quick retrospective of the module, with the goal of identifying possible ways to:
	- Improve your experience in later CSI-4 modules
	- Improve future cohorts' experience in this module

## Discussion

- How quickly would a topic on FFI be? How prepared am I to give a high level overview?
	- Not prepared:
		- High level idea is "straightforward" ("just use the calling convention to set up the stack / registers properly, then jump to the start of the function")
		- Details e.g. in Go are really messy, because Go runtime is pretty complicated (e.g. goroutine scheduler / garbage collection are all relevant)
	- Recommendation: look for some docs on how to get it working
		- https://pkg.go.dev/cmd/cgo or https://go.dev/blog/cgo
		- "Java Native Interface"
			- See e.g. https://www.baeldung.com/jni
	- Open question: can we add it to one of the interpreters we've been building?

- How much do we care about "how long it takes a compiler to compile itself"?
	- More generally, how do we balance getting a more optimized binary vs. having faster compilation?
	- e.g. V language brags about being able to compile its own compiler in ~0.5 seconds

- Why will the language a compiler is written in matter?
	- If we have a source file implemented in Go (i.e. `compiler.go`), we need an executable that can compile Go source code, in order to build it!

- Let's say it's 2008, and there aren't any Go compilers
- How do we solve this "chicken and egg problem"
	- Let's say we write `go_compiler.go`
		- Can we turn this into an executable?
			- No!
	- Step 1: Implement a Go compiler in C
		- `go_compiler.c`
		- Can we turn this into an executable?
			- Yes! We have C compilers
			- `gcc go_compiler.c -o go_compiler_v1.exe`
			- Now we have a binary!
				- `go_compiler_v1.exe hello_world.go`
					- Produces `hello_world.exe`
	- Step 2: 
		- Run `go_compiler.v1.exe go_compiler.go` to produce `go_compiler_v2.exe`
	- Step 3:
		- Throw away `go_compiler.c`, pretend that never happened
	- Step 4: Profit
		- What happens if we run `go_compiler_v2.exe go_compiler.go`?
			- The output should be identical, byte for byte, with `go_compiler_v2.exe`

- Let's say `go_compiler.c` is not "full" implementation of Go, but it only has a limited set of Go features
	- Let's say it doesn't handle `range` statement
	- `gcc go_compiler.c -o go_compiler_v1.exe`

OK, let's say `hello.go` contains this function:

```go
func f(arr []int) int {
	sum := 0
	for i := range arr {
		sum += a[i]
	}
	return sum
}
```

- What would happen if we do `go_compiler_v1.exe hello.go`?
	- Syntax error, "unrecognized token `range`"

OK, now let's think about the source code in `go_compiler.go`

- This can't contain any range statements (otherwise `go_compiler_v1.exe` can't handle it)
	- But can it SUPPORT range statements?
	- Yes, it might contain code like the following, which SUPPORTs range statements but doesn't actually USE range statements:

```go
// Maybe this is in `go_compiler_v3.go`
if nextToken == "range" {
	// Handle logic of range statement
}
```

- Self-hosting:
	- https://en.wikipedia.org/wiki/Bootstrapping_(compilers)

- What are the steps of a compiler?
	- Tokenization (aka scanning)
	- Parsing
	- Semantic Analysis
	- Code Generation
	- Optimization

- Compiler bootstrapping
	- What's the motivation for implementing the compiler for a language in that language?
		- Test it ("dogfooding")
		- Assume your language is good (it's pleasant)
			- Go is nicer to use than C, for example
		- It's easier for people to get started contributing
		- See [Go 1.3 Compiler Overhaul](https://docs.google.com/document/u/0/d/1P3BLR31VA8cvLJLfMibSuTdwTuF7WWLux71CYD0eeD8/mobilebasic)
	- How is it possible for the Go compiler to be implemented in Go?
	- Did anyone try implementing a quine (not as essential, but amusing)?
	- What's the "Trusting Trust" attack?
		- See https://www.cs.cmu.edu/~rdriley/487/papers/Thompson_1984_ReflectionsonTrustingTrust.pdf

```go
// evil_compiler_v0.go
if isGoCompiler(inputSourceCode) {
	// Manipulate binary so that it:
	// - makes bank programs evil
	// - makes future compilers evil
}

// evil_compiler_v0.exe go_compiler.go  # produces go_compiler_v1.exe
// go_compiler_v1.exe bank_program.go   # produces an evil bank program
// go_compiler_v1.exe go_compiler_v2.go # produces an evil compiler that will corrupt future bank programs AND compilers
```

- Idea: Can make an evil binary more resilient against reverse engineering by randomly flipping bits so that the program still works but a disassembler / etc. crashes when you open it

- Retrospective
	- Overall thoughts on the module?
		- Relevance?
		- Interest level?
		- Difficulty?
			- Eric: Thought most difficult was compilers / interpreters one, but also liked that one
			- Panashe: Enjoyed that one too
			- Branden: Liked more difficult
			- Bryan: Compilers topics seemed bimodal among members of class (unlike e.g. assembly programming)
			- John: Seemed difficult to target sweet spot for everybody
				- Good to order things in prework easy -> difficult
				- **Helpful to draw line "if you get past this point you won't be confused"**
			- Dan: Had more trouble fully completing prework for Thursday classes vs. having the weekend
				- Branden: More incremental assignment, put indicator of "good enough" point
				- Patch: Scanning / parsing most concrete example (parsing landed on Thursday)
	- What were the most and least useful sessions?
		- John: Depends a lot on personal preference
			- Found concurrency / GC interesting
		- Panashe: stack / heap and GC useful
			- compilers and interpreters fav. prework
		- Branden: liked compilers
		- Greg: +1 concurrency
		- Eric: GC get people most promotions
			- Personal interest liked interpreters
	- Preworks to iterate on?
		- Places to make it more fine-grained:
			- Bytecode generation exercise: smaller steps between test cases
		- Making the scanner / parser handle subset of Go?
			- Branden: prefer
			- John: prefer something that builds on itself each class
			- Eric prefers simple query language
		- Python2.7 pain to install using conda
		- Stack / heap / memory allocation prework could use iteration
			- Exercise on buffer overflow there could be helpful
		- GC could come after interpretation, could then build a simple one there
		- Asa: Want more support with testing / benchmarks on the stack / heap / memory allocation one
			- Liked compilers / interpreters prework
		- Thoughts on concurrency preworks?
			- John: Liked Go one (debugging)
				- Would've preferred to write more practical program (than the toy ID generation problem)
			- Branden: liked concurrency though seemed different from everything else (everything else felt more related)
				- But could mitigate it by making part of assignment related to other stuff?
			- Asa: VSCode gave some hints about Go concurrency questions
				- Patch: Could look up default LSP impl that's servicing hint suggestions for Go ("what's giving me these spoilers?")
		- Patch: Harder to determine if concurrency stuff will stick (we only have one week with it, and we don't interact with it as frequently)
	- What takeaways do you wish you had about concurrency?
		- Bryan: "Interleavings" was probably best part of module
			- Say I'm writing DB client / server
				- "How many threads to delegate to each before we get thread starvation / resource usage too high, thread pool vs. thread per client"
		- Branden: Most people deal with distributed systems; need to make sophisticated concurrent system is not as high
			- Exercise of knowing when you want to start breaking things into concurrent programs (rather than handling things in one process / thread), e.g. threads for audio / visual / main game loop: "why you want it in separate threads"
	- Is there anything (topics, exercises) you wish we added / removed?
	- Thoughts on specific topics:
		- Memory management / garbage collection?
		- Compilers?
			- Should we split the compilers-intro prework into two parts (code generation by itself)?
			- Should we expand the scanning / parsing content across two units?
				- John: Would only want to do this if we also add more on intermediate representation, LLVM, etc.
			- Branden: If we spend more time on compilers, could add more on optimizations (e.g. lattice structures, type analysis)
	- Opinions on ordering of the units?
	- Opinions on the title / possible renaming of the entire module?
		- Systems Programming and Compilers
			- +1 from Eric
		- Languages, Compilers, and Interpreters
		- Other?
			- Source to Machine Code
			- Systems Programming and Language Internals
	- Other feedback?
		- Ilyas: Liked first half more and felt it was more applicable
		- John A: similar experience, and felt more behind on preworks for compilers when things got busy
		- John O: concurrency was very important, don't want to lose that
	- Should we take out "special topics"?
		- Greg: Thought bootstrapping discussion was very useful
			- Dan: +1
		- Bharat: Liked peer aspect of last session
		- Eric: Wished had more time to prepare, liked idea but wished had been thinking about it from beginning
			- **Mention at beginning of module**
	- Greg: Some of subject matter, wish I had a cheat sheet (e.g. efficiency improvements from Intro to Systems) or a wiki to trace through what we discussed and review it (especially if it comes up on the job)
		- John A: +1 to that
		- Dan: can imagine "mountain diagram" from Crafting Interpreters being on there
		- Greg: something along the lines of answers to exercises in cheat sheet form, or crowdsource Anki cards?
		- Patch: active recall without a prompt is usually a pretty effective means of checking you have those encoded
			- e.g. draw how these things are related
		- Branden: the best way to bring back memory / recall is to create a "memory entrance" to it, imo
	- Branden: Good takeaway: taking advanced topic and "group classing" it seemed pretty positive (implementing bytecode interpreter)
		- Take it further in live sessions, implement together
	- John: Mentioned in first session, favorite thing is clarifying, agree with that
		- Might be interesting to have a clarification or two queued up in the hopper to get things rolling
	- Greg: I think in some of the classes, it may have been good to establish some terms, and agree on definitions, although I think the garbage collection session was the one where we had the most diverging definitions of terms (like block, etc)
		- (reduce ambiguity being the goal, as well as narrowing focus)