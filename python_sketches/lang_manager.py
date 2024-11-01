import os
import pathlib
from documents import Document

class Lang_manager:
    def __init__(self):
        self.document = None
        self.current_project_name = ""
        self.projects_folder_path = "pytnon_sketches\\projects"
    def run(self):
        print_greetings()
        while True:
            res = input()
            if res == "-help":
                print_info()
            elif res == "-open":
                self.open_project()
            elif res == "-create":
                self.create_project()
            elif res == "-status":
                self.send_status()
            elif res == "-run":
                self.run_project()
            else:
                print("unknown command: " + res)

    def send_status(self):
        if self.document == None:
            print("No project selected")
        else:
            print("your current project is: " + self.current_project_name)
            
    def open_project(self):
        self.current_project_name = ask_project_name()
        project_file_name = self.current_project_name + ".txt"
        project_path = os.path.join(os.path.curdir, "python_sketches", "projects", project_file_name)
        project_path = os.path.abspath(project_path)
        if os.path.exists(project_path):
            print("opening project...")
            self.document = Document(project_path)
            print("opened!")
            print("type -run to run the chosen project")
        else:
            print("no such project " + self.current_project_name + " was found in /projects.")
            self.current_project_name = ""
    
    def create_project(self):
        self.current_project_name = ask_project_name()
        project_file_name = self.current_project_name + ".txt"
        project_path = os.path.join(os.path.curdir, "python_sketches", "projects", project_file_name)
        project_path = os.path.abspath(project_path)
        project_file = pathlib.Path(project_path)
        project_file.touch()
        self.document = Document(project_path)
        print("Project "+ self.current_project_name + "was succesfully created!")

    def run_project(self):
        if self.document == None:
            print("No project selected")
        else:
            self.document.start_reader()

def print_info():
    print("type -create to create new project")
    print("type -open to open an existing project")
    print("type -status to see information about your current project")
    print("type -run to run your current project")
    
def print_greetings():
    print("welcome to Lang Construct! type -help to see all available commands.")
def ask_project_name():
    print("enter the name of you project:")
    name = input()
    return name

def __Main__():
    Lang_manager().run()
__Main__()