from Expressions import Expression

class Definition: #def zero belong Natural
    def __init__(self, define:Expression):
        self.expr = define

class Definition_creater:#def zero belong Natural
    def __init__(self, messanger, define:str):
        self.messanger = messanger
        self.def_expr = Expression(define[4:]) #zero belong Natural
        self.notify_definition_created = self.__defaultNotifyDefinitionCreated

    def __defaultNotifyDefinitionCreated(self, d:Definition):
        pass
    
    def create(self) -> Definition:
        definition = Definition(self.def_expr)
        self.notify_definition_created(definition)
        return definition
        


