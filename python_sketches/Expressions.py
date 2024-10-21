class Expression:
    def __init__(self, str:str): 
        self.content = str
    
    def compare_content(self, expr):
        if self.content != expr.content:
            return False
        else:
            return True
    def find_expression(self, expr:str): # Expressions always can consist of other expressions
        size = len(expr)
        if len(self.content) < size:
            return False
        i = 0
        while i+size <= len(self.content):
            if(self.content[i:size+i] == expr):
                return True
            i+=1
        return False
    def read_to_expression(self, expr:str) -> str: # return string of symbols defore the given combination of symbols
        i = 0
        size = len(expr)
        result = ""
        while i < len(self.content):
            if self.content[i : i + size] == expr :
                return result
            else:
                result += self.content[i]
            i+=1