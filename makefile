 
.PHONY: core
core:
	(cd src/core; go install)

.PHONY: game
game:
	(cd src; go build -o ../$@)

all: core game

run:
	./game