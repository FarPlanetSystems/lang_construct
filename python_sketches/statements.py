from Expressions import Expression
from rules import Rule

class Statement: #have:one + one belong Natural_from sum_1 (one,one) :one belong Natural,one belong Natural
    def __init__(self, rule:Rule, conclusion:Expression, params, premisses):
        self.rule = rule
        self.conclusion = conclusion
        self.params = params
        self.premisses = premisses
        self.verified = False
        self.notify_statement_verified = self.__standard_verified_notification
        self.notify_statement_not_verified = self.__standard_not_verified_notification

    def __standard_verified_notification(self, s):
        print("statement verified in the line ")
    def __standard_not_verified_notification(self, s):
        print("statement was not verified")
    
    def verify(self, legal_expressions):
        self.verified = self.rule.check(self.params, self.premisses, self.conclusion, legal_expressions)
        if self.verified:
            self.notify_statement_verified(self)
        else:
            self.notify_statement_not_verified(self)
        



class Statement_creator: #have:one + one belong Natural_from sum_1 (one,one) :one belong Natural,one belong Natural
    def __init__(self, messanger, line:str, line_num:int, doc_rules):
        self.line = line
        self.line_num = line_num
        self.doc_rules = doc_rules
        self.messanger = messanger
    
    def create(self) -> Statement:
        if not Expression(self.line).find_expression("have"):
            self.messanger("Compilation error: expressing a statement, have: was expected. line "+ str(self.line_num))
            return
        
        self.__structure()
        return Statement(self.rule, self.conclusion, self.params, self.premisses)
    
    def __structure(self): #line = "have:one + one belong Natural_from sum_1 (one,one) :one belong Natural,one belong Natural"
        conclusion_line = ""
        rule_name_line = ""
        arguments_line = ""

        i = 6 # counter we use to follow the statement line. we start by index i skipping "have :" symbols
        while self.line[i:i+5] != "_from":
            conclusion_line += self.line[i] # we read the concluion line at first
            i+=1
        
        i += 6 # after we have read the conclusion line, we need to skip  the following "_from " symbols

        while self.line[i] != " ":
            rule_name_line += self.line[i] # we read the name of the used rule
            i+=1
        
        arguments_line = self.line[i:]  #(one,one) :one belong Natural,one belong Natural

        self.conclusion = self.__read_statement_conclusion(conclusion_line)

        self.rule = self.__read_rule(rule_name_line)

        self.__read_arguments(arguments_line)



    def __read_statement_conclusion(self, line:str) -> Expression:  #one+one belong Natural
        return Expression(line)
    
    def __read_rule(self, line:str) -> Rule: #sum_1"
        for i in self.doc_rules:
            if i.name == line:
                return i
    
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
                    self.messanger("Compilation error: stating an expression a parameter was expected. Line "+ self.line_num)
                    return
                params.append(parameter)
            else:
                parameter += line[i]
            i+=1
        params+=parameter
        self.params = params #["one", "one"]
    
    def __read_premisses(self, line:str):
        premisses = []
        i = 0
        while line[i] != ":":
            i += 1
        i += 1
        premiss = ""
        while i < len(line):
            if line[i] == ",":
                if premiss == "":
                    self.messanger("Compilation error: stating an expression a reference to a rule was expected. Line "+ self.line_num)
                    return
                premisses.append(premiss)
            else:
                premiss += line[i]
            i+=1
        premisses.append(premiss)
        #premisses = ["one belong Natural", "one belong Natural"]
        for j in range(len(premisses)):
            premisses[j] = Expression(premisses[j])
        self.premisses = premisses

    


        
