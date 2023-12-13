# TOYLINGO
An interpreted language built on top of golang.


## Quick Start
* Clone the repository
* Run the main.go file to see the interpreter in action (it executes a sample program in the samples folder)

# Features

## Variables
* declare a variable using the let keyword
* the let keyword is optional

## Template Strings
to be implemented soon

## Conditionals
Support for if, else if, else statements

```
if(true){

} elseif(5<7){

}else{
    
}
```

## Loops
Support for a basic loop structure with break statement
```
loop(boolean_expression){
    //statements
    if(condition){
        break;
    }else{}
}
```

## Functions
Support for functions with return values and arguments
All functions in a scope are hoisted to the top of the scope

```
func(a,b) add{
    return a+b;
}
```




# Peculiarities

### + operator on strings
Any arithmetic expression before or after a string will be converted to a string and concatenated to the string. 
But only the expression succeeding the last occurence of a string in an expression will be evaluated.


# Bugs
 
* 8/9\*18 not evaluated correctly (9*18 evaluated first): right to left evaluated
* if conditional chain is not terminated by an else, the next token is consumed.

* return value and argument lists of functions must be non empty 
* the programs sometimes glitch and control flow becomes entirely different

* if statements must always have a corresponding else 


# Notes
* Unassigned array and object declarations will not even be executed