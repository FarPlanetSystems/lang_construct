from document import MP_document
from rules import *

class Expression:
    def __init__(self, str:str): 
        self.content = str
    
    def compare_content(self, expr):
        if self.content != expr.content:
            return False
        else:
            return True
    def find_expression(self, expr): # Expressions always can consist of other expressions
        size = len(expr.content)
        if len(self.content) < size:
            return False
        i = 0
        while i+size <= len(self.content):
            if(self.content[i:size+i] == expr):
                return True
            i+=1
        return False

class Definition:
    def __init__(self, define:Expression, line:int):
        self.expr = define
        self.line = line

class Definition_creater:
    def __init__(self, define:str, line:int, doc:MP_document):
        self.def_expr = define
        self.def_expr = Expression(line)
        self.def_doc = doc
    def create(self) -> Definition:
        definition = Definition(self.def_expr, self.def_expr)
        self.def_doc.definitions.append(definition)
        return definition
        


