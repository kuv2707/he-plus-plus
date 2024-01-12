# TOYLINGO
An interpreted general-purpose programming language built on top of golang.


## Quick Start
* Clone the repository
* Setup the .env fild taking reference from the .env.example file
* Run the main.go file to see the interpreter in action (it executes a sample program in the samples folder)
* There are many sample programs in the samples directory to try out
* The file to be executed can be specified in the .env file

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




# Peculiarities

### + operator on strings
Any arithmetic expression before or after a string will be converted to a string and concatenated to the string. 
But only the expression succeeding the last occurence of a string in an expression will be evaluated.


# Bugs
 
* 8/9\*18 not evaluated correctly (9*18 evaluated first): 
operators of same precedence are evaluated right to left instead of left to right

* many bugs in parsing brackets in expressions, like:
- an array declaration in the arguments of a function call
```
print([1,2,3]);
```
* original array returned instead of copy when exp is evaluated, so its memory is freed whenever a scope is exited
* unassigned array is not cleaned up when scope is exited
* empty array declaration crashes interpreter
* line numbers in error messages are not correct sometimes



# Notes
* Expression parsing is terrible!
* The VM runs on 1kb of memory.


# TODO
* coalesce Variable and Pointer structs
* Implement multidimensional arrays
* Implement template strings
* Implement objects (dot operator)
* Implement tuples (comma operator)
* implement module system and imports
* improve memory management