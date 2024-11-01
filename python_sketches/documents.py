from definitions import *
from statements import *
from rules import *
from os import path, curdir
import pathlib

class Document:

    def __init__(self, file_name:str):
        self._rules = []
        self._definitions = []
        self._legal_expressions = []
        self.document_file_name = file_name
        self._doc_reader = Document_sheet_reader(self)
    
    def reset(self):
        self._rules = []
        self._definitions = []
        self._legal_expressions = []
    
    def _add_definition(self, d:Definition):
        self._definitions.append(d)
        self._legal_expressions.append(d.expr)

    def _add_rule(self, rule:Rule):
        self._rules.append(rule)
    
    def _verify_statement(self, statement:Statement):
        self._legal_expressions.append(statement.conclusion)
    
    def start_reader(self):
        self.reset()
        with open(self.document_file_name, "r") as document_file:
            self._doc_reader.set_message_line()
            self._doc_reader.read_sheet(document_file)

class Document_sheet_reader:
    def __init__(self, document:Document):
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

        for l in lines:
            line = self.__parse_line(l)
            key_word = Expression(line).read_to_expression(" ")
            if key_word == "def":
                self.__read_definition(line, line_num)
            elif key_word == "rule":
                self.__read_rule(line, line_num)
            elif key_word == "have":
                self.__read_statement(line, line_num)
            elif key_word == "#":
                pass
            elif key_word == "import":
                self.__read_imported_document(line, line_num)
            elif key_word == None:
                self.message("Compilation error: key word was expected in the line " + str(line_num + 1) )
            else :
                self.message("Compilation error: unknown key word " + key_word + " in the line " + str(line_num+1))
            line_num += 1

    def __parse_line(self, line:str):
        for i in range(len(line)):
            if line[i] == " ":
                line = line[i+1:]
            else:
                return line

    def __read_imported_document(self, line:str, line_num:int):
        i = 0
        doc_name = ""
        while i != " ":
            i+=1
        for j in len(line[i:]):
            doc_name += line[j]
        import_document(self.doc, doc_name)

    def __read_definition(self, line:str, line_num:int):
        def_creator = Definition_creater(self.message, line, line_num + 1)
        def_creator.notify_definition_created = self.doc._add_definition
        def_creator.create()

    def __read_rule(self, line:str, line_num:int):
        rule_creator = Rule_creator(self.message, line, line_num + 1)
        rule_creator.notify_rule_created = self.doc._add_rule
        rule_creator.Create()

    def __read_statement(self, line:str, line_num:int):
        statement_creator = Statement_creator(self.message, line, line_num + 1, self.doc._rules)
        new_statement = statement_creator.create()
        new_statement.notify_statement_verified = self.doc._verify_statement
        new_statement.verify(self.doc._legal_expressions)

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

def import_document(destination:Document, imported_document_name:str) -> bool:
    project_file_name = imported_document_name + ".txt"
    project_path = path.join(imported_document_name, project_file_name)
    project_file = pathlib.Path(project_path)
    if project_file.is_file():
        imported_document = Document(project_path)
        destination._rules.append(imported_document._rules)
        destination._definitions.append(imported_document._definitions)
        destination._legal_expressions.append(imported_document._legal_expressions)
        return True
    else: 
        return False