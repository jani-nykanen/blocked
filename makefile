 
.PHONY: core
core:
	(cd src/core; go install)

.PHONY: game
game:
	(cd src; go build -o ../$@)

.PHONY: game.exe
game.exe:
	(cd src; go build -o ../$@ -ldflags -H=windowsgui)

all: linux
linux: core game
windows: core game.exe

run:
	./game
