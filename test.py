
while True:
    f = open("test_doc.txt", "r")
    print(f.read())
    i = 0
    while i < 10000000:
        i += 1
    i = 0


        initial_content = ""
        with open(self.doc.document_file_name, "r") as doc_file:
            initial_content = doc_file.read()
        if Expression(initial_content).find_expression(self.end_string):
            initial_content = Expression(initial_content).read_to_expression(self.end_string)
        initial_content += "\n" + self.end_string
        with open(self.doc.document_file_name, "w") as doc_file:
            doc_file.write(initial_content + msg)