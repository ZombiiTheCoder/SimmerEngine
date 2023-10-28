import os

os.chdir("src")
os.system("go build -o ../bin/SimmerJs.exe")
os.chdir("../bin")
# os.system("clear")
os.system("SimmerJs.exe ../tests/bytecode.js -compact ../tests/bytecode_compacted.js")
os.chdir("../")
