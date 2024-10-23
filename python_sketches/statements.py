from Expressions import Expression
from rules import Rule

class Statement: #have:one + one belong Natural_from sum_1 (one,one) :one belong Natural,one belong Natural
    def __init__(self, rule:Rule, conclusion:Expression, params, premisses):
        self.rule = rule
        self.conclusion = conclusion
        self.params = params
        self.premisses = premisses
        self.verified = False
        self.read_succesfully = True
        self.notify_statement_verified = self.__standard_verified_notification
        self.notify_statement_not_verified = self.__standard_not_verified_notification

    def __standard_verified_notification(self, s):
        print("statement verified in the line ")
    def __standard_not_verified_notification(self, s):
        print("statement was not verified")
    
    def verify(self, legal_expressions):
        if not self.read_succesfully:
            return
        self.verified = self.rule.check(self.params, self.premisses, self.conclusion, legal_expressions)
        if self.verified:
            self.notify_statement_verified(self)
        else:
            return
        



class Statement_creator: #have:one + one belong Natural_from sum_1 (one,one) :one belong Natural,one belong Natural
    def __init__(self, messanger, line:str, line_num:int, doc_rules):
        self.line = line
        self.line_num = line_num
        self.doc_rules = doc_rules
        self.messanger = messanger
        self.is_read_succesfully = True


    
    def create(self) -> Statement:

        if not Expression(self.line).find_expression("have"):
            self.messanger("Compilation error: expressing a statement, have : was expected. line "+ str(self.line_num))
            return
        
        is_structured_succesfully = self.__structure()
        conclusion = None
        rule = None
        params = None
        premisses = None

        if is_structured_succesfully:
            conclusion = self.__read_statement_conclusion(self.conclusion_line)
            rule = self.__read_rule(self.rule_name_line)
            params = self.__read_params(self.params_line)
            premisses = self.__read_premisses(self.premiss_line)

        if conclusion == None:
            conclusion = ""
            self.is_read_succesfully = False
        if premisses == None:
            self.is_read_succesfully = False
            premisses = []
        if rule == None:
            self.is_read_succesfully = False
            rule = Rule("", [], "", Expression(""))
        if params == None:
            self.is_read_succesfully = False
            params = []
        statement = Statement(rule, conclusion, params, premisses)
        if not is_structured_succesfully or not self.is_read_succesfully:
            statement.read_succesfully = False
        return statement
    
    def __structure(self) ->bool: #line = "have {one + one belong Natural} from sum_1 (one,one) [one belong Natural],[one belong Natural]"
        self.conclusion_line = ""
        self.rule_name_line = ""
        self.params_line = ""
        self.premiss_line = ""

        line_expression = Expression(self.line)

        if not line_expression.find_expression("{"): # curly brackets must be there
            self.messanger("Compilation error: declaring a body of the conclusion, closing bracket was expected. Line " + str(self.line_num))
            return False

        i = self.line.index("{") # counter we use to follow the statement line. we start by index i skipping "have :" symbols

        if not line_expression.find_expression("}"): # curly brackets must be closed
            self.messanger("Compilation error: declaring a body of the conclusion, closing bracket was expected. Line " + str(self.line_num))
            return False
        
        i += 1

        while self.line[i] != "}": # we read the concluion line in curly brackets at first
            self.conclusion_line += self.line[i]
            i+=1
        
        i += 1
        line_expression.content =self.line[i:]

        if not line_expression.find_expression("from"):
            self.messanger("Compilation error: declaring the name of a rule, 'from' symbol was expected. Line " + str(self.line_num))
            return False

        while self.line[i:i+4] != "from": # we need to skip  the following "from " symbols
            i+= 1
        
        i += 4

        while self.line[i] == " ": # skip spaces
            i+=1
        
        while self.line[i] != " " and self.line[i] != "(": # we read the name of the used rule
            self.rule_name_line += self.line[i]
            i += 1
        
        while self.line[i] == " ": # skip spaces
            i+=1
        line_expression.content = self.line[i:]

        if not self.line[i] == "(":
            self.messanger("Compilation error: declaring parametres, () was expected. Line " + str(self.line_num))
            return False
        else:
            i+=1
        if not line_expression.find_expression(")"):
            self.messanger("Compilation error: declaring parametres, () was expected. Line "  + str(self.line_num))
            return False
        
        while self.line[i] != ")": # we read parametres in brackets
            self.params_line += self.line[i]
            i+=1
        i+=1

        while self.line[i] == " ": # skip spaces
            i+=1

        self.premiss_line = self.line[i:]

        return True

    def __read_statement_conclusion(self, line:str) -> Expression:  #one+one belong Natural
        return Expression(line)
    
    def __read_rule(self, line:str) -> Rule: #sum_1"
        for i in self.doc_rules:
            if i.name == line:
                return i
        self.messanger("No such rule " + line + " was found")
    
    def __read_params(self, line:str):#one,one
        params = line.split(",")
        for i in params:
            if i == "":
                self.messanger("Compilation error: an non-empty parameter was expected. Line" + str(self.line_num))
                self.is_read_succesfully = False
                return []
        return params
        
    
    def __read_premisses(self, line:str): #[one belong Natural],[one belong Natural]
        premisses = line.split(",")
        for i in range(len(premisses)):
            premiss = premisses[i]
            j = 0
            while premiss[j] == " ":
                j += 1
            if premiss[j] != "[":
                self.messanger("Compilation error: premisses must be declared in square brackets. Line " + str(self.line_num))
                self.is_read_succesfully = False
                return []
            else:
                j += 1
            if not Expression(premiss).find_expression("]"):
                self.messanger("Compilation error: square brackets must be closed. Line " + str(self.line_num))
                self.is_read_succesfully = False
                return []
            closing_bracket_index = premiss.index("]")
            premisses[i] = premiss[j:closing_bracket_index]
            premisses[i] = Expression(premisses[i])
        return premisses
        

    


        
