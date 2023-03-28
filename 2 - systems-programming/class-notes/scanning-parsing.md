## Objectives

By the end of this session, you should understand how to implement a handwritten parser and scanner for a simple language, including details such as:

- How to use a "context-free grammar" to specify the syntax of a language and handle issues such as precedence
- The "recursive descent" algorithm for parsing, as well as a key limitation of this approach

In addition, you should understand what tools are available for automatically generating scanners / parsers, how to use these tools, and tradeoffs between handwritten and automatically generated scanners / parsers.

## Agenda

Logistical note on final session

- Other topic requests?
    - Since there's no security module, buffer overflows / trusting trust seems most interesting
        - On the other hand, there's a lot of tutorials on YouTube; maybe compiler bootstrapping would be more interesting
    - Would prefer to go more in depth
    - Constraint solvers e.g. Z3: how to take an abstract tool and apply it to many different problems
        - Starting point: https://norvig.com/sudoku.html (constraint satisfaction)
    - Implementing compiler optimizations could be cool too (e.g. register coloring)
    - Something concurrency related that's related to future modules
    - Software design related to concurrency
    - More resources on foreign function calls
    - Trusting Trust:
        - https://niconiconi.neocities.org/posts/ken-thompson-really-did-launch-his-trusting-trust-trojan-attack-in-real-life/
- Lightning talks?
    - No
- Consensus:
    - We'll spend ~1/2 the session focusing on compiler bootstrapping / "Trusting Trust attack"
    - Other half, people can do some quick presentations on topics of interest (feel free to choose)

### Discussion / Solution Walkthrough

We will start by discussing the problem / how to approach it, as we step through a reference solution.

What difficulties did you run into / where were there blockers

- Branden: Trying to utilize "generics" / polymorphism in Go, to utilize the visitor pattern
    - What's the problem you're trying to solve?
        - Trying to write a single function that can handle many AST types
    - Review of interfaces vs. structs in Go
- Patch: Some initial ambiguity around concept of a grammar / "how do I turn my tokens into a tree"
    - Started implementing first, then went back to update grammar
    - Figured out precedence issue after implementing initial version
- Eric: Noticing the symmetry between precedence and order of grammar rules, that's when things clicked
    - Started without thinking about precedence
    - Tried both bottom-up and top-down; neither worked until after thinking about precedence
- Asa: Started with fully parenthesized output to focus on how to build syntax tree

```lisp
(OR
    (AND
        (TERM hello),
        (TERM world)),
    (AND
        (TERM alice),
        (NOT
            (TERM BOB)
        ))
```

- John:
    - Let's focus more on parsing
    - Took a while to get started with scanning, had trouble deciding what tokens to start with, eventually decided to stick with bare bones
        - Thinking about everything you could do with the query string
            - Do you tokenize the space?
            - How do you handle two terms that come immediately after each other?
            - Should we require parentheses?
            - Should grouped search terms be a new token?
            - Ended up going with quotes for terms (treated as single term)
            - "hello" "world"

- Should we handle parens in the scanner or the parser?
    - During scanning, don't need to worry about grouping; just output it as a token
    - During parsing, actually figuring out the grouping

- Eric: There's some ambiguity about what you can put in scanner or parser
    - e.g. spaces?
    - Panashe:
        - Scanner's job is to find tokens, smallest units that make sense
        - Parser's job to get the structured representation
            - Dan: "Implementing the grammar of the language"
            - Greg: "treeify"

HELLO WORLD ->
    [TERM(HELLO) WHITESPACE TERM(WORLD)]
    [TERM(HELLO) TERM(WORLD)] <- What I mean by "strip whitespace"

HELLO AND WORLD ->
    [TERM(HELLO) AND TERM(WORLD)]
    [TERM(HELLO) TERM(WORLD)] <- Do this instead

- Could we put everything in parser / avoid scanner?
    - Technically yes, but seems very verbose
    - Nice separation of concerns / decomposition into independent parts

- What helper methods were helpful for structuring code?
    - Peek for both
        - 1 character lookahead / 1 token lookahead
    - Match
        - In textbook, you could hand it multiple tokens

- Bryan: split using a regex
    - `src.split(/([^a-zA-Z0-9])/)`

**Scanning**

- What do we need to think about before starting?
- How do we structure the code?
- What sort of helper functions will be useful?
- What's the overall approach we're using?
- How do we handle:
    - Whitespace
    - Operators
    - Quoted phrases
    - Words

**Parsing**

- What do we need to think about before starting?
- What's a context-free grammar, and why is it useful for solving this problem?
- How do we structure the code?
- What sort of helper functions will be useful?
- What's the overall approach we're using?

Patch's version
- `parse` calls into `__query`
- `expression` import has all the nodes we're going to use
- `Query` is top level entry point

### Context-Free Grammar

- Caveats:
    - Open question: why is it called "context free", what makes a grammar "context-sensitive"?
    - Note: Python3 using something called a "PEG" which has slightly different rules than a context-free grammar

- Open question:
    - Is https://github.com/bradfield-csi-4/patchner/blob/main/systems_programming/scanning_parsing/query_grammar.txt an ambiguous grammer?
        - I think it depends on the nuances of how `*` is defined here

- "Regular expressions are a lower level of power than CFGs?"
    - Regular expressions are a notation for describing regular languages
        - [a-z]+
            - This is a notation for describing all lowercase words
    - You e.g. cannot write down a regular expression for "all balanced parentheses"
        - But you CAN write down a context-free grammar describing it!

```
expr -> '' | '(' expr ')'
```

- What issues can we run into with the "recursive descent" algorithm?

```
expr : expr '+' expr
     | NUMBER
```

```python
# 5 + 5
# input: [NUMBER(5), ADD, NUMBER(5)]

# What goes wrong here?
def __expr():
    if match(NUMBER):
        return Number(current())
	else:
        left = __expr()
        match(PLUS)
        right = __expr()
        return BinaryAdd(left, right)
```

### antlr4 Example
