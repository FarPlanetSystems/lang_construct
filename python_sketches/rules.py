from Expressions import Expression 

class Rule:
    def __init__(self, name:str, variables, premisses, conclusion:Expression): #rule _name any x y : x belong natural , y belong natural -> x+y belong natural
        self.name = name
        self.premisses = premisses
        self.variables = variables
        self.conclusion = conclusion
        self.read_succesfully = True
        self.sendMessage = self.__default_send_message

    def __default_send_message(line:str):
        print(line)
    
    def check(self, var_values, premisses, suppoused_conclusion:Expression, legal_expressions):
        if not self.read_succesfully:
            return False

        if len(var_values) != len(self.variables):
            return False

        for i in premisses:
            premiss_found = False
            for j in legal_expressions:
                if i.content == j.content:
                    premiss_found = True
            if premiss_found == False:
                self.sendMessage("No verified premiss " + i.content + " found while verifing rule " + self.name)


        context_premisses = []
        for i in self.premisses:
            context_premisses.append(i.clone())
        #reforming the exprected premisses with the given variables values
        for i in range(len(context_premisses)):
            context_premisses[i] = self.__imply_variables(context_premisses[i], var_values)
        context_conclusion = self.__imply_variables(self.conclusion.clone(), var_values) # the same to the exprected conclusion

        for i in range(len(premisses)): # here we compare the exprected premisses with the given. 
            if not premisses[i].compare_content(context_premisses[i]):
                self.sendMessage("given premiss "+ premisses[i].content +" does not match the requirements of the rule " + self.premisses[i].content)
        if not suppoused_conclusion.compare_content(context_conclusion):
            self.sendMessage("concluded statement "+ suppoused_conclusion.content +" does not match the expected conclusion " + self.conclusion.content)
        else:
            return True
    

    def __imply_variables(self, expr:Expression, values):
        for i in range(len(expr.content)):
            for j in range(len(self.variables)):
                if expr.content[i] == self.variables[j]:
                    expr.content = expr.content.replace(expr.content[i], values[j])
        return expr

class Rule_creator:
    def __init__(self, messanger, line:str, line_num:int): #rule sum_1 any x , y : x belong Natural, y belong Natural -> x + y belong Natural
        self.line = line
        self.line_num = line_num
        self.notify_rule_created = messanger
        self.messanger = messanger
        self.key_word = "rule"
        self.is_read_succesfully = True

    def Create(self) -> Rule:
        is_srtuctured_succesfully = self.__strcuture()# dividing the rule line into smaller logical parts

        params = None
        premisses = None
        conclusion = None
        name = None

        if is_srtuctured_succesfully:
            params = self.__read_params(self.params_line)#handling the part with variables
            premisses = self.__read_premisses(self.premiss_line)#handling the part with premisses
            conclusion = self.__read_conclusion(self.conclusion_line)#handling the part with the conclusion
            name = self.__read_name(self.name_line)

        if name == None:
            name = ""
            self.is_read_succesfully = False
        if params == None:
            params = []
            self.is_read_succesfully = False
        if conclusion == None:
            conclusion = Expression("")
            self.is_read_succesfully = False
        if premisses == None:
            self.premisses = []
            self.is_read_succesfully = False
        
        result = Rule(name, params, premisses, conclusion)
        result.sendMessage = self.messanger
        self.notify_rule_created(result)
        if not is_srtuctured_succesfully or not self.is_read_succesfully:
            result.read_succesfully = False
        return result
    
    def __strcuture(self) -> bool: #rule sum_1 (x, y) : [x belong Natural], [y belong Natural] > {x + y belong Natural}
        self.name_line = ""
        self.params_line = ""
        self.premiss_line = ""
        self.conclusion_line = ""
        i = 0
        line_expression = Expression(self.line)
        # here we are looking for the key_word
        while self.line[i:i+4] != self.key_word: 
            i+=1
        i += len(self.key_word)
        # we are skipping the empty spaces between the key word and the rule name
        while self.line[i] == " ":
            i+=1
        # we are reading the name of our rule
        while self.line[i] != " " and self.line[i] != "(":
            self.name_line += self.line[i] 
            i+=1
        #we are skipping the empty spaces between the name and parametres
        while self.line[i] == " ":
            i += 1
        if self.line[i] != "(":
            self.messanger("Compilation error: defining a rule, brackets () for parametres were expected. Line " + str(self.line_num))
            return False
        else:
            i += 1
        line_expression.content = self.line[i:]
        if not line_expression.find_expression(")"):
            self.messanger("Compilation error: defining a rule, brackets () for parametres should be closed. Line " + str(self.line_num))
            return False
        # we are reading the required parametres
        closing_brackets_index = self.line.index(")")
        self.params_line = self.line[i:closing_brackets_index]
        i = closing_brackets_index
        i += 1
        #we are skipping the empty spaces
        while self.line[i] == " ":
            i += 1
        # checking the presence of ':' symbol
        if self.line[i] != ":":
            self.messanger("Compilation error: defining a rule, transition symbol : was expected. Line " + str(self.line_num))
            return False
        else:
            i += 1
        #we are skipping the empty spaces
        while self.line[i] == " ":
            i += 1
        line_expression.content = self.line[i:]
        # checking the presence of '>' symbol
        if not line_expression.find_expression(">"):
            self.messanger("Compilation error: defining a rule, refering symbol > was expected. Line " + str(self.line_num))
            return False
        # we are reading the premisses
        refering_symbol_index = line_expression.content.index(">")
        self.premiss_line = line_expression.content[0:refering_symbol_index]
        i = i + refering_symbol_index
        i += 1
        #we are skipping the empty spaces
        while self.line[i] == " ":
            i += 1
        # checking the presence of the curly brackets
        if self.line[i] != "{":
            self.messanger("Compilation error: introducing the conclusion of a rule, curly brackets were expected. Line " + str(self.line_num))
            return False
        else:
            i += 1
        line_expression.content = self.line[i:]
        if not line_expression.find_expression("}"):
            self.messanger("Compilation error: introducing the conclusion of a rule, curly brackets must be closed. Line " + str(self.line_num))
            return False
        # we are reading the conclusion
        closing_brackets_index = line_expression.content.index("}")
        self.conclusion_line = line_expression.content[0:closing_brackets_index]
        i = i + closing_brackets_index
        return True

    def __read_params(self, line:str): #x, y
        params = line.split(",")
        for i in params:
            if i == "":
                self.messanger("Compilation error: an non-empty parameter was expected. Line " + str(self.line_num))
                self.is_read_succesfully = False
                return []
        return params
        
    def __read_premisses(self, line:str): # [x belong Natural], [y belong Natural]
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

    def __read_conclusion(self, line:str):# x + y belong Natural
        concl_expr = Expression(line)
        return concl_expr
    def __read_name(self, line:str):
        rule_name = line
        return rule_name