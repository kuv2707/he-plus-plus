# A toy compiler
I use this project to learn about compilers and interpreters. It is a work in progress.
The syntax is freaky, and the keywords are in Spanish.

# Notable things done so far:
### Lexical analysis
Supports single line comments, string literals, numbers in bases 2, 8, 10, and 16, identifiers, keywords, operators, and punctuation.
I've used Trie to match keywords and operators.
I've written some tests for the lexer (doesn't cover many cases yet).

### Parsing to AST
Supports variable declarations, function declarations, if-else statements, return statements, expression statements, function calls, member access, array access, binary expressions, unary expressions, literals (string, number, boolean, null), identifiers.

Uses *Pratt parsing* for expressions.

*Lexing and parsing are done concurrently using Go channels and goroutines*.

### Basic static analysis (type checking, variable declaration checks, etc.)
This is a WIP.
I've planned to support type casting.

The rough end goal is to generate assembly for the language, and make it able to work with C libraries.