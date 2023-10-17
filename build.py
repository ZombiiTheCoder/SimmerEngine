import os
src = []
for root, paths, files in os.walk("src"):
    for file in files:
        print(os.path.abspath(os.path.join(root, file)))
        if (file.__contains__(".go")):
            os.path.abspath(os.path.join(root, file))

os.chdir("src")
os.system("go build -o ../bin/SimmerJs.exe "+" ".join(src))
os.chdir("../")
os.system("./bin/SimerJs tests/numbers.js")