from Expressions import Expression

class Definition: #def [zero belong Natural]
    def __init__(self, define:Expression):
        self.expr = define

class Definition_creater:#def [zero belong Natural]
    def __init__(self, messanger, line:str, line_num:int):
        self.messanger = messanger
        self.line = line
        self.line_num = line_num
        self.def_line = ""
        self.notify_definition_created = self.__defaultNotifyDefinitionCreated
        self.key_word = "def"

    def __defaultNotifyDefinitionCreated(self, d:Definition):
        pass
    
    def create(self) -> Definition:
        is_read_succesfully = self.read()
        if is_read_succesfully:
            result =  Definition(Expression(self.def_line))
            self.notify_definition_created(result)
            return result
        else:
            return Definition(Expression(""))

        
    def read(self):
        i = 0
        line_expression = Expression(self.line)
        while self.line[i] == " ":
            i += 1
        i += len(self.key_word)
        while self.line[i] == " ":
            i += 1
        if self.line[i] != "[":
            self.messanger("Compilation error: definitions must be placed in square brackets. Line" + str(self.line_num))
            return False
        else:
            i += 1
        line_expression.content = self.line[i:]
        if not line_expression.find_expression("]"):
            self.messanger("Compilation error: the square brackets must be closed. Line" + str(self.line_num))
            return False
        closing_brackets_index = line_expression.content.index("]")
        self.def_line = line_expression.content[0:closing_brackets_index]
        return True        


