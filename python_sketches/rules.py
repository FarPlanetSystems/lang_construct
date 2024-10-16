from reader import Expression 
from document import MP_document

class Rule:
    def __init__(self, name:str, variables, premisses, conclusion:Expression): #rule _name any x y : x belong natural , y belong natural -> x+y belong natural
        self.name = name
        self.premisses = premisses
        self.variables = variables
        self.conclusion = conclusion
        
    
    def check(self, var_values, premisses, suppoused_conclusion:Expression):
        if len(var_values) != len(self.variables):
            return False

        context_premisses = self.premisses #reforming the exprected premisses with the given variables values
        for i in range(len(context_premisses)):
            self.__imply_variables(context_premisses[i], var_values)
        context_conclusion = self.__imply_variables(self.conclusion, var_values) # the same to the exprected conclusion


        for i in range(len(premisses)): # here we compare the exprected premisses with the given. 
            if not premisses[i].compare_content(context_premisses[i]):
                return False
        if not suppoused_conclusion.compare_content(context_conclusion):
            return False
        else:
            return True
    

    def __imply_variables(self, expression:Expression, values):
        expr = expression
        for i in range(len(expr.content)):
            for j in range(len(self.variables)):
                if expr.content[i] == self.variables[j]:
                    expr.content[i] = values[j]
        return expression

class Rule_creator:
    def __init__(self, line:str, document:MP_document): #rule sum_1 any x , y : x belong Natural, y belong Natural -> x + y belong Natural
        self.rule_doc = document
        self.line = line

    def Create(self, rule_expr:Expression) -> Rule:

        rule_expr = Expression(self.line) # checking the syntax
        if rule_expr.find_expression(Expression(":")) == False:
            raise Exception("Error, in the rule " + self.line + " symbol : was exprected")
        if rule_expr.find_expression(Expression("->")) == False:
            raise Exception("Error, in the rule " + self.line + " symbol -> was exprected")
        
        self.__strcuture()# dividing the rule line into smaller logical parts

        self.__read_variables()#handling the part with variables

        self.__read_premisses()#handling the part with premisses

        self.__read_conclusion()#handling the part with the conclusion
        
        return Rule(self.ruleName, self.variables, self.premisses, Expression(self.conclusion))
    
    def __strcuture(self): #rule sum_1 any x, y : x belong Natural, y belong Natural -> x + y belong Natural
        name = ""#sum_1
        var_field = ""#any x, y 
        premiss_field = ""# x belong Natural, y belong Natural 
        concl_field = ""# x + y belong Natural

        i = 5 #we pass the "rule" flagg
        while self.line[i + 1 : i + 4] != "any":
            name += self.line[i]
        i+=1 #we pass the space symbol between the name and the section fot variables
        while self.line[i] != ":":
            var_field += self.line[i]
            i += 1
        i+= 1 # we pass the : symbol
        while self.line[i] != "-" and self.line[i+1] != ">":
            premiss_field += self.line[i]
            i+=1
        i += 2 # we pass the -> symbol
        while i < len(self.line):
            concl_field += self.line[i]
        
        self.ruleName = name
        self.variables = var_field
        self.premisses = premiss_field
        self.conclusion = concl_field


    def __read_variables(self): #any x, y
        variables = []
        i = 0
        if self.variables[0:3] != "any": # we could have a rule without variables, that is, a definition
            self.variables = variables
            return 
        else:
            i = 4 # else we pass the "any" symbol
        while i < len(self.variables):
            variable = ""
            if self.variables[i] == ",":
                if variable == "": # we dont allow to add an empty variable - if you have a comma, you have to write a variable before
                    raise Exception("rule compilation error: a variable was expected on symbol "+i)
                variables.append(variable)
            else:
                variable+=self.variables[i]
            i += 1
        if  variable != "": # for we add variables until we achieve the end of the line, the last variable arent going to be taken
            variables.append(variable)
        self.variables = variables

    def __read_premisses(self): # x belong Natural, y belong Natural 
        premisses = []
        i = 0
        while i < len(self.premisses):
            premiss = ""
            if self.premisses[i] == ",":
                if premiss == "":
                    raise Exception("rule compilation error: a premiss was expected on symbol "+i)
                premisses.append(premiss)
            else:
                premiss += self.premisses[i]
        self.premisses = premisses

    
    def __read_conclusion(self):# x + y belong Natural
        concl_expr = Expression(self.conclusion)
        if concl_expr.find_expression(Expression(",")):
            raise Exception("rule compilation error: an unexpected symbol , in the rule's conclusion")