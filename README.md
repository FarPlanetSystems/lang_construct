# LangConstruct

A tool for defining formal systems.

inspired by .lean, coq etc, but easier and works for every formal system (defined by you).
It gives simplest tools for defining rules of inference and assists you checking if every element of the system is correctly infered from these rules.

## Installation

run the following command:

```
go install github.com/FarPlanetSystems/lang_construct
```

## Usage

- run lang_construct.
- after this the folder "lang_construct_projects" is supposed to appear in the current directory. Create there a plain text file, for example test.txt.
- write some code:

```
def "A be statement";
rule exMiddle (x): "[x] be statement" -> "[x] or not [x]";
have "A or not A" from exMiddle (A) "A be statement";
```

- enter the name of the file in the command line
- then check the initiate file again: there you will find a message from lang_construct after @ symbol. If everything is ok you will get the "coherence verified!" report

```
def "A be statement";
rule exMiddle (x): "[x] be statement" -> "[x] or not [x]";
have "A or not A" from exMiddle (A) "A be statement";

@
Coherence verified!
```
