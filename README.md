# TOYLINGO
An interpreted language built on top of golang.

# Features

## Variables
declare a variable using the let keyword
the let keyword is optional

## Template Strings
Full support for JS style template strings

## Conditionals
Support for if, else if, else statements



# Peculiarities

### + operator on strings
Any arithmetic expression before or after a string will be converted to a string and concatenated to the string. 
But only the expression succeeding the last occurence of a string in an expression will be evaluated.