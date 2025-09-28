# lang_construct

Tool for defining formal systems.

inspired by .lean, coq etc, but easier and works for every formal system (defined by you).
It gives simplest tools for defining rules of inference and assists you checking if every element of the system is correctly infered from these rules.

## Installation

run the following command:

```
go install github.com/FarPlanetSystems/lang_construct
```

## Usage

- write a .txt file with some code:

```
def "A be statement";
rule exMiddle (x): "[x] be statement" -> "[x] or not [x]";
have "A or not A" from exMiddle ("A") "A be statement";
```

- run lang_construct passing the file name:

```
lang_construct "example.txt"
```

- see the interpretor message in the initial code file. If ho errors occured, you will find the `coherence verified!` report:

```
def "A be statement";
rule exMiddle (x): "[x] be statement" -> "[x] or not [x]";
have "A or not A" from exMiddle ("A") "A be statement";

@
Coherence verified!
```
