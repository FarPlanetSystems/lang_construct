from Expressions import Expression 

class Rule:
    def __init__(self, name:str, variables, premisses, conclusion:Expression): #rule _name any x y : x belong natural , y belong natural -> x+y belong natural
        self.name = name
        self.premisses = premisses
        self.variables = variables
        self.conclusion = conclusion
        self.sendMessage = self.__default_send_message

    def __default_send_message(line:str):
        print(line)
    
    def check(self, var_values, premisses, suppoused_conclusion:Expression, legal_expressions):
        if len(var_values) != len(self.variables):
            return False

        for i in premisses:
            premiss_found = False
            for j in legal_expressions:
                if i.content == j.content:
                    premiss_found = True
            if premiss_found == False:
                self.sendMessage("No verified premiss " + i.content + " found while verifing rule " + self.name)


        context_premisses = self.premisses #reforming the exprected premisses with the given variables values
        for i in range(len(context_premisses)):
            self.__imply_variables(context_premisses[i], var_values)
        context_conclusion = self.__imply_variables(self.conclusion, var_values) # the same to the exprected conclusion

        for i in range(len(premisses)): # here we compare the exprected premisses with the given. 
            if not premisses[i].compare_content(context_premisses[i]):
                self.sendMessage("given premiss "+ premisses[i].content +" does not match the requirements of the rule " + self.premisses[i].content)
        if not suppoused_conclusion.compare_content(context_conclusion):
            self.sendMessage("concluded statement "+ suppoused_conclusion.content +" does not match the expected conclusion " + self.conclusion.content)
        else:
            return True
    

    def __imply_variables(self, expression:Expression, values):
        expr = expression
        for i in range(len(expr.content)):
            for j in range(len(self.variables)):
                if expr.content[i] == self.variables[j]:
                    expr.content = expr.content.replace(expr.content[i], values[j])
        return expression

class Rule_creator:
    def __init__(self, messanger, line:str): #rule sum_1 any x , y : x belong Natural, y belong Natural -> x + y belong Natural
        self.line = line
        self.notify_rule_created = messanger
        self.messanger = messanger

    def Create(self) -> Rule:

        rule_expr = Expression(self.line) # checking the syntax
        if rule_expr.find_expression(":") == False:
            self.messanger("Compilation error, in the rule " + self.line + " symbol : was exprected")
            return

        if rule_expr.find_expression("->") == False:
            self.messanger("Compilation error, in the rule " + self.line + " symbol -> was exprected")
            return
        
        self.__strcuture()# dividing the rule line into smaller logical parts

        self.__read_variables()#handling the part with variables

        self.__read_premisses()#handling the part with premisses

        self.__read_conclusion()#handling the part with the conclusion

        result = Rule(self.ruleName, self.variables, self.premisses, Expression(self.conclusion))
        result.sendMessage = self.messanger
        self.notify_rule_created(result)
        return result
    
    def __strcuture(self): #rule sum_1 any x, y : x belong Natural, y belong Natural -> x + y belong Natural
        name = ""#sum_1
        var_field = ""#any x, y 
        premiss_field = ""# x belong Natural, y belong Natural 
        concl_field = ""# x + y belong Natural

        i = 5 #we pass the "rule" flagg
        while self.line[i + 1 : i + 4] != "any":
            name += self.line[i]
            i+=1
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
            i+=1
        
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
        variable = ""
        while i < len(self.variables):
            if self.variables[i] == ",":
                if variable == "": # we dont allow to add an empty variable - if you have a comma, you have to write a variable before
                    self.messanger("Compilation error: a variable was expected in rule "+self.ruleName)
                    return
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
        premiss = ""
        while i < len(self.premisses):
            if self.premisses[i] == ",":
                if premiss == "":
                    self.messanger("Compilation error: a premiss was expected in rule "+ self.ruleName)
                    return
                premisses.append(premiss)
            else:
                premiss += self.premisses[i]
            i+=1
        premisses.append(Expression(premiss))
        self.premisses = premisses

    
    def __read_conclusion(self):# x + y belong Natural
        concl_expr = Expression(self.conclusion)
        if concl_expr.find_expression(","):
            self.messanger("Compilation error: an unexpected symbol , in " + self.ruleName + " conclusion")
            return