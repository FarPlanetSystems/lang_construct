# LangConstruct
A tool for defining formal languages
## Installation
run the following command
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
have "A or not A" from exMiddle : "A be statement";
```