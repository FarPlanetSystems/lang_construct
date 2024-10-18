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

class Definition: #def zero belong Natural
    def __init__(self, define:Expression, line:int):
        self.expr = define
        self.line = line

class Definition_creater:#def zero belong Natural
    def __init__(self, define:str, line:int, doc:MP_document):
        
        self.def_line = line
        self.def_doc = doc
        self.def_expr = Expression(define[5:]) #zero belong Natural
    def create(self) -> Definition:
        definition = Definition(self.def_expr, self.def_line)
        self.def_doc.definitions.append(definition)
        self.def_doc.legal_expressions.append(self.def_expr)
        return definition
        


