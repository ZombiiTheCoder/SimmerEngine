import os

os.chdir("src")
os.system("go build -o ../bin/SimmerJs.exe")
os.chdir("../bin")
os.system("clear")
os.system("SimmerJs.exe ../tests/fib.js -compact ../tests/fib_compacted.js")
os.chdir("../")
