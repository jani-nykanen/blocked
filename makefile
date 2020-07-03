 
.PHONY: core
core:
	(cd src/core; go install)

.PHONY: game
game:
	(cd src; go build -o ../$@)

.PHONE: game.exe
game.exe:
	(cd src; go build -o ../$@)

all: linux
linux: core game
windows: core game.exe

run:
	./game
