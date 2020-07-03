 
.PHONY: core
core:
	(cd src/core; go install -ldflags -H=windowsgui)

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
