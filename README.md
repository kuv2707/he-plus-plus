# he++
An interpreted general-purpose programming language built on top of golang.


## Quick Start
* Clone the repository
* Setup the .env fild taking reference from the .env.example file
* Run the main.go file to see the interpreter in action (it executes a sample program in the samples folder)
* There are many sample programs in the samples directory to try out
* The file to be executed can be specified in the .env file
* A command line interface will be implemented soon

## REPL
The interpreter supports REPL mode too.
To test it, set REPL=1 in the .env file.

# Features

## Variables
* declare a variable using the let keyword
* the let keyword is optional

## Template Strings
to be implemented soon

## Conditionals
* Support for if, else if, else statements
* parentheses around the boolean expression are optional

```
if(true){

} elseif(5<7){

}else{
    
}
```

## Loops
* Support for a basic loop structure with break statement
* parentheses around the boolean expression are optional
```
loop(boolean_expression){
    //statements
    if(condition){
        break;
    }
}
```

## Functions
Support for functions with return values and arguments
Arguments are passed by value
All functions in a scope are hoisted to the top of the scope

```
function add(a,b){
    return a+b;
}

let k=add(5,6);
```
Some natively implemented functions are available in every scope, like:
* ```print``` and ```println``` : Prints the argument to stdout.
* ```readNumber``` : Returns a number scanned from stdin.
* ```len``` : Accepts an array or string and returns its length.
* ```makeArray``` : Accepts a size and returns an array of that size, populated with nulls.
* ```random``` : Returns a pseudo-random number between 0 and 1.

## Arrays
An array can have any number of elements of any type, including other arrays. There can be a trailing comma after the last element.
```js
let arr=[1,2,3,4,[5,6,7],"hello world!",];
```


# Peculiarities

### Right to Left precedence
Operators have right to left precedence, so the expression to the right of an operator will be evaluated first.
One implication of this peculiarity is that when concatenating strings, if the rightmost operand is an expression, it will be evaluated arithmetically and then concatenated, while everything else will be concatenated.
```js
let a=2+3+"wow"; //23wow
let b="wow"+2+3; //wow5
```


# Bugs
 
* 8/9\*18 not evaluated correctly (9*18 evaluated first): 
operators of same precedence are evaluated right to left instead of left to right

* unassigned array is not cleaned up when scope is exited
* line numbers in error messages are not correct 
* lexer is crap



# Notes
* The VM runs on 1kb of memory.


# TODO
* Implement template strings
* Implement objects (dot operator)
* Implement tuples (comma operator)
* implement module system and imports
* improve memory management
* implement default args to function
* implement variadic args to function
* pass line no when printing error messages

## Contemplating
* Pratt parsing