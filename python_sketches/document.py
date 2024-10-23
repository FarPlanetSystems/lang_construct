from definitions import *
from statements import *
from rules import *

class MP_document:

    def __init__(self, file_name:str):
        self.theorems = []
        self.rules = []
        self.definitions = []
        self.legal_expressions = []
        self.document_file_name = file_name
        self.doc_reader = Document_sheet_reader(self)
    
    def reset(self):
        self.theorems = []
        self.rules = []
        self.definitions = []
        self.legal_expressions = []
    
    def add_definition(self, d:Definition):
        self.definitions.append(d)
        self.legal_expressions.append(d.expr)

    def add_rule(self, rule:Rule):
        self.rules.append(rule)
    
    def verify_statement(self, statement:Statement):
        self.legal_expressions.append(statement.conclusion)
    
    def start_reader(self):
        while True:
            i = 0
            while i < 10000000:
                i += 1
            i = 0
            self.reset()
            with open(self.document_file_name, "r") as document_file:
                self.doc_reader.set_message_line()
                self.doc_reader.read_sheet(document_file)

class Document_sheet_reader:
    def __init__(self, document:MP_document):
        self.end_string = "@---"
        self.doc = document
    
    def message(self, msg:str):
        with open(self.doc.document_file_name, "a") as doc_file:
            doc_file.write("\n"+msg)

    def set_message_line(self):
        file_content = ""
        with open(self.doc.document_file_name, "r") as doc_file:
            file_content = doc_file.read()
        content_expression = Expression(file_content)
        code = file_content
        if content_expression.find_expression(self.end_string):
            code = content_expression.read_to_expression(self.end_string)
        new_content = code + "@---"
        with open(self.doc.document_file_name, "w") as doc_file:
            doc_file.write(new_content)

    def read_code(self, code:str):
        lines = self.separete_expressions(code)
        line_num = 0

        for line in lines:
            key_word = Expression(line).read_to_expression(" ")
            if key_word == "def":
                def_creator = Definition_creater(self.message, line, line_num + 1)
                def_creator.notify_definition_created = self.doc.add_definition
                def_creator.create()
            elif key_word == "rule":
                rule_creator = Rule_creator(self.message, line, line_num + 1)
                rule_creator.notify_rule_created = self.doc.add_rule
                rule_creator.Create()
            elif key_word == "have":
                statement_creator = Statement_creator(self.message, line, line_num + 1, self.doc.rules)
                new_statement = statement_creator.create()
                new_statement.notify_statement_verified = self.doc.verify_statement
                new_statement.verify(self.doc.legal_expressions)
            elif key_word == None:
                self.message("Compilation error: key word was expected in the line " + str(line_num + 1) )
            else :
                self.message("Compilation error: unknown key word " + key_word + " in the line " + str(line_num+1))
            line_num += 1
        
    def read_sheet(self, file):
        content = file.read()
        self.content_modified = self.__format_file_content(content)
        code = self.content_modified.read_to_expression(self.end_string)
        self.read_code(code)

    def __format_file_content(self, content:str):
        content = content.replace("\n", "")
        content_expr = Expression(content)
        if not content_expr.find_expression(self.end_string):
            content_expr.content += self.end_string
            self.message("")
        return content_expr
    
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
        

MP_document("D:/stuff/IT/python/New_age/LangConstruct/python_sketches/doc1.txt").start_reader()