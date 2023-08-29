#!./ookami

colorify(false, 0, 0, 255)
print("Re-Compiling the shell")

os.execute("go build .")

colorify(false, 0, 255, 0)
print("Done")

resetColor()

print("Executing")

os.execute("./ookami")
