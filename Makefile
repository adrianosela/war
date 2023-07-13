NAME:=$(shell basename `git rev-parse --show-toplevel`)

all: setbin

setbin: build
	cp $(NAME) /usr/local/bin

build: 
	go build -o $(NAME)
