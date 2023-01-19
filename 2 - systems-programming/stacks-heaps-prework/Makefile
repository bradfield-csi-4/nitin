viz:
	cc -lncurses -lpanel viz.c -o viz

# The target flag can be omitted for non M1 macs
fact:
	cc --target=x86_64-apple-darwin-macho -fomit-frame-pointer fact.c -o fact

.PHONY : clean
clean :
	rm -f viz fact
