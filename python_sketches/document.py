from reader import *

class MP_document:

    def __init__(self):
        self.theorems = []
        self.rules = []
        self.definitions = []
        self.legal_expressions = []
        self.doc_reader = Document_sheet_reader(self)
    
    def find_rule(self, rule_name:str):
        for i in self.rules:
            if i.name == rule_name:
                return i
        return None
    def find_theorem(self, theorem_name:str):
        for i in self.theorems:
            if i.name == theorem_name:
                return i
        return None
    def find_definition(self, def_name:str):
        for i in self.definitions:
            if i.name == def_name:
                return i
    
    
    
    def read(self, text:str):
        lines = self.doc_reader.separete_expressions(text)
        for i in range(len(lines)):
            line = lines[i]
            if line[0:3] == "def":
                def_creator = Definition_creater(line, i, self)
                def_creator.create()
            elif line[0:4] == "rule":
                rule_creator = Rule_creator(self)
                rule_creator.Create()

class Document_sheet_reader:
    def __init__(self, document:MP_document):
        self.doc = document
    def separete_expressions(self, line:str):# having a line of chars, we separate them by ";" into several independet lines
        res = []
        expr = ""
        for i in line:
            if i != ";":
                expr += i
            else:
                res.append(expr)
                expr = ""
        return res
    def read_sheet(text:str):
        pass