from reader import Expression
from document import MP_document
from rules import Rule

class Statement: #have:one + one belong Natural_from sum_1 (one,one) :one belong Natural,one belong Natural
    def __init__(self, rule:Rule, conclusion:Expression, params, premisses, line_num: int):
        self.rule = rule
        self.conclusion = conclusion
        self.params = params
        self.premisses = premisses
        self.line_num = line_num
        self.verified = False
        self.notify_statement_verified = self.__standard_verified_notification
        self.notify_statement_not_verified = self.__standard_not_verified_notification

    def __standard_verified_notification(self):
        print("statement verified in the line " + self.line_num)
    def __standard_not_verified_notification(self):
        print("unverified statement in the line " + self.line_num)
    
    def verify(self):
        self.verified = self.rule.check(self.params, self.premisses, self.conclusion)
        if self.verified:
            self.notify_statement_verified()
        else:
            self.notify_statement_not_verified()
        



class Statement_creator: #have:one + one belong Natural_from sum_1 (one,one) :one belong Natural,one belong Natural
    def __init__(self, line:str, line_num:int, doc: MP_document):
        self.line = line
        self.line_num = line_num
        self.doc = doc
    def create(self) -> Statement:
        if not Expression(self.line).find_expression(Expression("have:")):
            raise Exception("expressing a statement, have: was expected")
        
        self.__structure()
        return Statement(self.rule, self.conclusion, self.params, self.premisses)
    
    def __structure(self): #line = "have:one + one belong Natural_from sum_1 (one,one) :one belong Natural,one belong Natural"
        conclusion_line = ""
        rule_name_line = ""
        arguments_line = ""

        i = 5 # counter we use to follow the statement line. we start by index i skipping "have:" symbols
        while self.line[i:i+5] != "_from":
            conclusion_line += self.line[i] # we read the concluion line at first
        
        i += 6 # after we have read the conclusion line, we need to skip  the following "_from " symbols

        while self.line[i] != " ":
            rule_name_line += self.line[i] # we read the name of the used rule
        
        arguments_line = self.line[i:]  #(one,one) :one belong Natural,one belong Natural

        self.conclusion = self.__read_statement_conclusion(conclusion_line)
        self.rule = self.__read_rule(rule_name_line)
        self.__read_arguments(arguments_line)



    def __read_statement_conclusion(self, line:str) -> Expression:  #one+one belong Natural
        return Expression(line)
    
    def __read_rule(self, line:str) -> Rule: #sum_1"
        return self.doc.find_rule(line)
    
    def __read_arguments(self, line:str): #(one,one) :one belong Natural,one belong Natural
        self.__read_params(line)
        self.__read_premisses(line)
        
    def __read_params(self, line:str):
        params = []
        i = 0
        while line[i] != "(":
            i += 1
        i += 1
        parameter = ""
        while line[i] != ")":
            if line[i] == ",":
                if parameter == "":
                    raise Exception("statement compilation error: parameter was expected in the line "+ self.line_num)
                params.append(parameter)
            else:
                parameter += line[i]
        self.params = params #["one", "one"]
    
    def __read_premisses(self, line:str):
        premisses = []
        i = []
        while line[i] != ":":
            i += 1
        i += 1
        premiss = ""
        while i < len(line):
            if line[i] == ",":
                if premiss == "":
                    raise Exception("statement compilation error: rule was expected in the line "+ self.line_num)
                premisses.append(premiss)
            else:
                premiss += line[i]
        #premisses = ["one belong Natural", "one belong Natural"]
        for i in range(len(premisses)):
            premisses[i] = Expression(premisses[i])
        self.premisses = premisses

    


        
