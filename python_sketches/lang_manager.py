class Lang_manager:
    def __init__(self):
        self.documents = []
        self.control_file_name = "lang_control.txt"
    def start(self):
        self.control_file = open(self.control_file_name, "w")
        while True:
            content = self.control_file.read()