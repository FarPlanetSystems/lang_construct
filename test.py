content = "ahvbriuwfoer"
expr = "ahvbriuwfoer"

size = len(expr)
if len(content) < size:
    print(False)
i = 0
while i+size <= len(content):
    if(content[i:size+i] == expr):
        print(True)
    i+=1
print(False)